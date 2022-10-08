package main

import (
	"net/http"

	"route53ddns/internal/dynamicdns"
	"route53ddns/internal/ipdetector"
	"route53ddns/pkg/route53"

	log "github.com/sirupsen/logrus"
)

func main() {
	config := NewConfig()

	httpClient := &http.Client{}
	detector := ipdetector.NewDetector(httpClient)

	route53Client, err := route53.NewClient(config.AWSRegion, config.AWSAccessKeyId, config.AWSSecretAccessKey, config.Route53HostedZoneId)
	if err != nil {
		log.Fatalf("failed to build route53 client (%v)", err)
	}

	dynamicDNS, err := dynamicdns.NewDynamicDNS(detector, route53Client)
	if err != nil {
		log.Fatalf("failed to build dynamicdns (%v)", err)
	}

	err = dynamicDNS.Update(config.Route53Domains)
	if err != nil {
		log.Fatalf("failed to update DNS entries (%v)", err)
	}
}
