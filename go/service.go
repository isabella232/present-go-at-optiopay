package service

import (
  "net/http"

  "github.com/julienschmidt/httprouter"
  "github.com/optiopay/messages"
  "github.com/optiopay/micro"
  "github.com/optiopay/transfer-job-service/events"
)

type Service struct {
  http.Handler   // embedded
  activeVouchers map[string]bool
  count          int64
}

func NewService() *Service {
  s := &Service{activeVouchers: make(map[string]bool)}
  r := httprouter.New()
  r.GET("/count", s.handleGetCount)
  s.Handler = r
  return s
}

func (s *Service) ProcessEvent(event messages.Event,
    writer micro.MessageWriter) error {
  switch ev := event.(type) {
  case *events.TransferJobCreated:
    s.rememberNewVoucher(ev)
  case *events.TransferJobStatusUpdated:
    s.countSoldVoucher(ev)
    s.forgetOldVoucher(ev)
  }
  return nil
}

func (s *Service) rememberNewVoucher(
    e *events.TransferJobCreated,
  ) {
  if e.CampaignID != "" { // It is a voucher: Remember it!
    s.activeVouchers[e.TransferJobID] = true
  }
}

func (s *Service) countSoldVoucher(
    e *events.TransferJobStatusUpdated
  ) {
  if e.Status == events.Status_CLOSED &&
     s.activeVouchers[e.TransferJobID] {

    s.count++
  }
}

func (s *Service) forgetOldVoucher(
    e *events.TransferJobStatusUpdated
  ) {
  switch e.Status {
  case events.Status_CANCELED,
       events.Status_ERROR,
       events.Status_CLOSED:
    delete(s.activeVouchers, e.TransferJobID)
  }
}

func (s *Service) handleGetCount(w http.ResponseWriter,
    r *http.Request, params httprouter.Params) {

  var content = struct {
    Count int `json:"count"`
  }{
    Count: s.count,
  }
  micro.JSONResp(w, content, http.StatusOK)
}
