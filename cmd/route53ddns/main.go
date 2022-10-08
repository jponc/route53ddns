package main

import (
	"net/http"

	"route53ddns/internal/dynamicdns"
	"route53ddns/internal/ipdetector"
	"route53ddns/pkg/route53"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	awsRoute53 "github.com/aws/aws-sdk-go/service/route53"
	log "github.com/sirupsen/logrus"
)

func main() {
	config := NewConfig()

	// Initialise detector
	httpClient := &http.Client{}
	detector := ipdetector.NewDetector(httpClient)

	// Initialise route53Client
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("global"),
		Credentials: credentials.NewStaticCredentials(config.AWSAccessKeyId, config.AWSSecretAccessKey, ""),
	})
	if err != nil {
		log.Fatalf("failed to build route53 client (%v)", err)
	}
	awsRoute53Client := awsRoute53.New(sess)

	route53Client, err := route53.NewClient(awsRoute53Client, config.Route53HostedZoneId)
	if err != nil {
		log.Fatalf("failed to build route53 client (%v)", err)
	}

	// initialise dynamicDNS
	dynamicDNS, err := dynamicdns.NewDynamicDNS(detector, route53Client)
	if err != nil {
		log.Fatalf("failed to build dynamicdns (%v)", err)
	}

	// Do the update against route53 domains
	err = dynamicDNS.Update(config.Route53Domains)
	if err != nil {
		log.Fatalf("failed to update DNS entries (%v)", err)
	}
}
