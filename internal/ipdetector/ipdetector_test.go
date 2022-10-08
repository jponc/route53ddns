package ipdetector_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"route53ddns/internal/ipdetector"

	"github.com/stretchr/testify/require"
)

type MockTransport struct {
	roundTripFn func(request *http.Request) (*http.Response, error)
}

func (m *MockTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if m != nil && m.roundTripFn != nil {
		return m.roundTripFn(request)
	}

	return &http.Response{
		StatusCode: 200,
		// Send response to be tested
		Body: ioutil.NopCloser(bytes.NewBufferString(`OK`)),
		// Must be set to non-nil value or it panics
		Header: make(http.Header),
	}, nil
}

func Test_NewDetector(t *testing.T) {
	t.Run("returns Detector", func(t *testing.T) {
		detector := ipdetector.NewDetector(&http.Client{})
		require.NotNil(t, detector)
		require.IsType(t, &ipdetector.Detector{}, detector)
	})
}

func Test_GetIP(t *testing.T) {
	type test struct {
		name       string
		transport  http.RoundTripper
		expectedIP string
		isErr      bool
	}

	tests := []test{
		{
			name: "return error because Get returned an error",
			transport: &MockTransport{
				roundTripFn: func(request *http.Request) (*http.Response, error) {
					require.Equal(t, request.URL.String(), "https://ipinfo.io/ip")
					return nil, fmt.Errorf("error")
				},
			},
			isErr: true,
		},
		{
			name: "doesn't return an error and returns back the ip address",
			transport: &MockTransport{
				roundTripFn: func(request *http.Request) (*http.Response, error) {
					require.Equal(t, request.URL.String(), "https://ipinfo.io/ip")
					return &http.Response{
						Body: ioutil.NopCloser(bytes.NewBufferString(`1.1.1.1`)),
					}, nil
				},
			},
			expectedIP: "1.1.1.1",
			isErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClient := &http.Client{Transport: tt.transport}

			detector := ipdetector.NewDetector(httpClient)
			ip, err := detector.GetIP()

			if tt.isErr {
				require.Error(t, err)
				require.Equal(t, "", ip)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedIP, ip)
			}
		})
	}
}
