package main

import (
	"github.com/optiopay/messages"
	"github.com/optiopay/micro"
	"github.com/optiopay/utils/log"
)

// application version, should be set during compilation time
var version string

func main() {
	service := NewService()
	conf := micro.MicroserviceConf{
		ServiceName:   "woodblock",
		ServiceID:     uint16(messages.WoodblockService.ID()),
		AppHash:       version,
		StateTopics:   []string{"transfer-job-events"},
		ProcessTopics: []string{"woodblock", "transfer-job-events"},
	}

	microservice := micro.NewMicroservice(service, conf)
	if err := microservice.Run(); err != nil {
		log.Fatal("service error", "err", err)
	}
}
