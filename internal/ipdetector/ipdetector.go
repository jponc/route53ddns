package ipdetector

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Detector struct{}

func NewDetector() *Detector {
	return &Detector{}
}

func (d *Detector) GetIP() (string, error) {
	resp, err := http.Get("https://ipinfo.io/ip")
	if err != nil {
		return "", fmt.Errorf("failed to request public ip (%w)", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to get ip request body (%w)", err)
	}

	return string(body), nil
}
