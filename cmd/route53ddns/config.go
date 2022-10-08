package main

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	AWSAccessKeyId      string
	AWSSecretAccessKey  string
	Route53HostedZoneId string
	Route53Domains      []string
}

func NewConfig() *Config {
	domainsCSV := getEnv("ROUTE53_DOMAINS")
	domains := strings.Split(domainsCSV, ",")

	return &Config{
		AWSAccessKeyId:      getEnv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey:  getEnv("AWS_SECRET_ACCESS_KEY"),
		Route53HostedZoneId: getEnv("ROUTE53_HOSTED_ZONE_ID"),
		Route53Domains:      domains,
	}
}

func getEnv(key string) string {
	val := os.Getenv(key)

	if val == "" {
		log.Fatalf("Environment variable %s not found", key)
	}

	return val
}
