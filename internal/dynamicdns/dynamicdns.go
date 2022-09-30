package dynamicdns

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Detector interface {
	GetIP() (string, error)
}

type Route53Client interface {
	UpdateRecords(domains []string, ip string) error
}

type DynamicDNS struct {
	detector      Detector
	route53Client Route53Client
}

func NewDynamicDNS(detector Detector, route53Client Route53Client) (*DynamicDNS, error) {
	if detector == nil {
		return nil, fmt.Errorf("dynamicdns not initialised, detector not found")
	}

	if route53Client == nil {
		return nil, fmt.Errorf("dynamicdns not initialised, detector not found")
	}

	return &DynamicDNS{
		detector:      detector,
		route53Client: route53Client,
	}, nil
}

func (d *DynamicDNS) Update(domains []string) error {
	ip, err := d.detector.GetIP()
	if err != nil {
		return fmt.Errorf("failed to detect ip address (%w)", err)
	}

	log.Printf("WAN IP > %s", ip)

	err = d.route53Client.UpdateRecords(domains, ip)
	if err != nil {
		return fmt.Errorf("failed to update domain entries (%w)", err)
	}

	return nil
}
