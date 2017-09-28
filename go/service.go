package service

import (
  "net/http"
  "sort"
  "time"

  "github.com/julienschmidt/httprouter"
  "github.com/optiopay/messages"
  "github.com/optiopay/micro"
  "github.com/optiopay/transfer-job-service/events"
  "github.com/optiopay/utils/log"
)

type Service struct {
  router         http.Handler
  activeVouchers map[string]bool

  // not using time.Time in order to save memory and
  // be more aligned with internal storage format
  soldVouchers []int64
}

// choose a save extreme value, that doesn't bust time.UnixNano()
const maxSecondsSince1970 = 0x222ffffff

var (
  minTimestamp = time.Unix(-maxSecondsSince1970, 0)
  maxTimestamp = time.Unix(maxSecondsSince1970, 0)
)

func NewService() *Service {
  s := &Service{
    activeVouchers: make(map[string]bool, 16*1024),

    // 64 MB of memory reserved at start (might grow later)
    soldVouchers: make([]int64, 0, 8*1024*1024),
  }

  r := httprouter.New()
  r.GET("/count", s.handleGetVouchersSold)
  s.router = r
  return s
}

func (s *Service) ServeHTTP(
    w http.ResponseWriter,
    r *http.Request
  ) {
  s.router.ServeHTTP(w, r)
}

func (s *Service) ProcessEvent(
    event messages.Event,
    writer micro.MessageWriter,
  ) error {

  switch ev := event.(type) {
  case *events.TransferJobCreated:
    s.rememberCampaigns(ev)
  case *events.TransferJobStatusUpdated:
    s.rememberVouchersSold(ev)
    s.forgetOldVoucher(ev)
  default:
    log.Debug("unknown event", "kind", ev.Kind())
  }
  return nil
}

func (s *Service) rememberCampaigns(
    e *events.TransferJobCreated,
  ) {
  if e.CampaignID != "" { // It is a voucher: Remember it!
    s.activeVouchers[e.TransferJobID] = true
    log.Debug("remembering: active voucher",
      "count", len(s.activeVouchers))
  }
}

func (s *Service) rememberVouchersSold(
    e *events.TransferJobStatusUpdated
  ) {
  if e.Status == events.Status_CLOSED &&
     s.activeVouchers[e.TransferJobID] {

    s.soldVouchers = addTimestampSorted(
        s.soldVouchers,
        e.CreatedAt
      )
    log.Debug("added: sold voucher",
        "count", len(s.soldVouchers))
  }
}

func (s *Service) forgetOldVoucher(
    e *events.TransferJobStatusUpdated
  ) {
  if e.Status == events.Status_CANCELED ||
    e.Status == events.Status_ERROR ||
    e.Status == events.Status_CLOSED {

    delete(s.activeVouchers, e.TransferJobID)
    log.Debug("forgot: active voucher",
        "count", len(s.activeVouchers))
  }
}

// CountOfActiveVouchers returns the number of currently remembered vouchers
// (for tests).
func (s *Service) CountOfActiveVouchers() int {
  return len(s.activeVouchers)
}

// CountOfSoldVouchers returns the overall number of sold vouchers (for tests).
func (s *Service) CountOfSoldVouchers() int {
  return len(s.soldVouchers)
}

// indexOfTimestamp returns the first index of the sorted slice for that
// tsSlice[index] >= ts is true.
func indexOfTimestamp(tsSlice []int64, ts int64) int {
  n := len(tsSlice)
  if n <= 0 {
    return -1
  }
  i := sort.Search(n, func(i int) bool {
    return tsSlice[i] >= ts
  })
  return i
}

// countInInterval counts the timestamps in the given interval.
// 'from' being included but 'to' excluded from the interval.
func countInInterval(tsSlice []int64, from, to int64) int {
  if from >= to { // handle empty interval
    return 0
  }
  i := indexOfTimestamp(tsSlice, from)
  j := indexOfTimestamp(tsSlice, to)
  return j - i // this should work for the edge cases, too
}

// addTimestampSorted adds the given timestamp
// and moves it to the correct position in the slice.
// This really behaves like insertion sort and should rarely move the timestamp
// by many places because the events are naturally almost sorted.
func addTimestampSorted(tsSlice []int64, ts int64) []int64 {
  tsSlice = append(tsSlice, ts)
  for i := len(tsSlice) - 1; i > 0 && tsSlice[i] < tsSlice[i-1]; i-- {
    tsSlice[i-1], tsSlice[i] = tsSlice[i], tsSlice[i-1]
  }
  return tsSlice
}

func (s *Service) handleGetVouchersSold(
  w http.ResponseWriter, r *http.Request, params httprouter.Params) {

  r.ParseForm()
  afterStr := r.Form.Get("after")
  beforeStr := r.Form.Get("before")

  var err error
  after := minTimestamp
  if afterStr != "" {
    after, err = checkAndParseTimestamp(afterStr, "after")
    if err != nil {
      micro.StdJSONErr(w, http.StatusBadRequest)
      return
    }
  }
  before := maxTimestamp
  if beforeStr != "" {
    before, err = checkAndParseTimestamp(beforeStr, "before")
    if err != nil {
      micro.StdJSONErr(w, http.StatusBadRequest)
      return
    }
  }

  // the start timestamp is inclusive
  count := countInInterval(
      s.soldVouchers,
      after.UnixNano()+1,
      before.UnixNano())

  var content = struct {
    Count int `json:"count"`
  }{
    Count: count,
  }
  micro.JSONResp(w, content, http.StatusOK)
}

func checkAndParseTimestamp(
    tsStr string, name string
  ) time.Time, err {

  var zeroTS time.Time
  ts, err = time.Parse(time.RFC3339, tsStr)
  if err != nil {
    log.Info("unable to parse start date", name, tsStr)
    return zeroTS, err
  } else if ts.Before(minTimestamp) ||
      ts.After(maxTimestamp) {
    log.Info("timestamp is out of range", name, ts)
    return zeroTS, errors.New("bad timestamp")
  }
  return ts, nil
}