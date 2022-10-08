package dynamicdns_test

import (
	"fmt"
	"testing"

	"route53ddns/internal/dynamicdns"

	"github.com/stretchr/testify/require"
)

type MockDetector struct {
	getIPFn func() (string, error)
}

func (m *MockDetector) GetIP() (string, error) {
	if m != nil && m.getIPFn != nil {
		return m.getIPFn()
	}

	return "1.1.1.1", nil
}

type MockRoute53Client struct {
	updateRecordsFn func(domains []string, ip string) error
	domainsUpdated  []string
}

func (m *MockRoute53Client) UpdateRecords(domains []string, ip string) error {
	m.domainsUpdated = domains

	if m != nil && m.updateRecordsFn != nil {
		return m.updateRecordsFn(domains, ip)
	}

	return nil
}

func Test_NewDynamicDNS(t *testing.T) {
	type test struct {
		name          string
		detector      dynamicdns.Detector
		route53Client dynamicdns.Route53Client
		isErr         bool
	}

	tests := []test{
		{
			name:          "returns error when detector is nil",
			detector:      nil,
			route53Client: &MockRoute53Client{},
			isErr:         true,
		},
		{
			name:          "returns error when route53Client is nil",
			detector:      &MockDetector{},
			route53Client: nil,
			isErr:         true,
		},
		{
			name:          "does not return error if both detector and route53Client are present",
			detector:      &MockDetector{},
			route53Client: &MockRoute53Client{},
			isErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf(">>> %v \n", tt.route53Client)
			fmt.Printf(">>> %v \n", nil)
			dynamicDNS, err := dynamicdns.NewDynamicDNS(tt.detector, tt.route53Client)
			if tt.isErr {
				require.Error(t, err)
				require.Nil(t, dynamicDNS)
			} else {
				require.NoError(t, err)
				require.NotNil(t, dynamicDNS)
			}
		})
	}
}

func Test_Update(t *testing.T) {
	type test struct {
		name          string
		detector      *MockDetector
		route53Client *MockRoute53Client
		domains       []string
		isErr         bool
	}

	tests := []test{
		{
			name: "returns error when GetIP returns an error",
			detector: &MockDetector{
				getIPFn: func() (string, error) {
					return "", fmt.Errorf("some error")
				},
			},
			route53Client: &MockRoute53Client{},
			isErr:         true,
		},
		{
			name:     "returns error when UpdateRecords returns an error",
			detector: &MockDetector{},
			route53Client: &MockRoute53Client{
				updateRecordsFn: func(domains []string, ip string) error {
					return fmt.Errorf("some error")
				},
			},
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dynamicDNS, err := dynamicdns.NewDynamicDNS(tt.detector, tt.route53Client)
			require.NoError(t, err)

			err = dynamicDNS.Update(tt.domains)
			if tt.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.ElementsMatch(t, tt.route53Client.domainsUpdated, tt.domains)
			}
		})
	}
}
