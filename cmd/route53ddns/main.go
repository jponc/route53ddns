package main

import (
	"route53ddns/internal/dynamicdns"
	"route53ddns/internal/ipdetector"
	"route53ddns/pkg/route53"

	log "github.com/sirupsen/logrus"
)

func main() {
	config := NewConfig()
	detector := ipdetector.NewDetector()

	route53Client, err := route53.NewClient(config.AWSRegion, config.AWSAccessKeyId, config.AWSSecretAccessKey, config.Route53HostedZoneId)
	if err != nil {
		log.Fatalf("failed to build route53 client (%w)", err)
	}

	dynamicDNS, err := dynamicdns.NewDynamicDNS(detector, route53Client)
	if err != nil {
		log.Fatalf("failed to build dynamicdns (%w)", err)
	}

	err = dynamicDNS.Update(config.Route53Domains)
	if err != nil {
		log.Fatalf("failed to update DNS entries (%w)", err)
	}
}
