package route53_test

import (
	"fmt"
	"testing"

	"route53ddns/pkg/route53"

	"github.com/aws/aws-sdk-go/aws"
	awsRoute53 "github.com/aws/aws-sdk-go/service/route53"
	"github.com/stretchr/testify/require"
)

type MockRoute53Client struct {
	changeResourceRecordSetsFn func(input *awsRoute53.ChangeResourceRecordSetsInput) (*awsRoute53.ChangeResourceRecordSetsOutput, error)
	input                      *awsRoute53.ChangeResourceRecordSetsInput
}

func (m *MockRoute53Client) ChangeResourceRecordSets(input *awsRoute53.ChangeResourceRecordSetsInput) (*awsRoute53.ChangeResourceRecordSetsOutput, error) {
	if m != nil {
		m.input = input
	}

	if m != nil && m.changeResourceRecordSetsFn != nil {
		return m.changeResourceRecordSetsFn(input)
	}

	return &awsRoute53.ChangeResourceRecordSetsOutput{}, nil
}

func Test_NewClient(t *testing.T) {
	type test struct {
		name          string
		route53Client route53.Route53Client
		isErr         bool
	}

	tests := []test{
		{
			name:          "returns an error if route53Client is nil",
			route53Client: nil,
			isErr:         true,
		},
		{
			name:          "does not return an error if route53Client is present",
			route53Client: &MockRoute53Client{},
			isErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := route53.NewClient(tt.route53Client, "abc-123")
			if tt.isErr {
				require.Error(t, err)
				require.Nil(t, client)
			} else {
				require.NoError(t, err)
				require.IsType(t, &route53.Client{}, client)
			}
		})
	}
}

func Test_UpdateRecords(t *testing.T) {
	type test struct {
		name          string
		route53Client *MockRoute53Client
		hostedZone    string
		domains       []string
		ip            string
		expectedInput *awsRoute53.ChangeResourceRecordSetsInput
		isErr         bool
	}

	tests := []test{
		{
			name: "returns error if it failed to change resource record sets",
			route53Client: &MockRoute53Client{
				changeResourceRecordSetsFn: func(input *awsRoute53.ChangeResourceRecordSetsInput) (*awsRoute53.ChangeResourceRecordSetsOutput, error) {
					return nil, fmt.Errorf("some error")
				},
			},
			hostedZone: "some-zone",
			domains:    []string{"domain1.com", "domain2.com"},
			ip:         "1.1.1.1",
			isErr:      true,
		},
		{
			name: "does not return error if changing the resource record sets are successful",
			route53Client: &MockRoute53Client{
				changeResourceRecordSetsFn: func(input *awsRoute53.ChangeResourceRecordSetsInput) (*awsRoute53.ChangeResourceRecordSetsOutput, error) {
					return &awsRoute53.ChangeResourceRecordSetsOutput{}, nil
				},
			},
			hostedZone: "some-zone",
			domains:    []string{"domain1.com", "domain2.com"},
			ip:         "1.1.1.1",
			expectedInput: &awsRoute53.ChangeResourceRecordSetsInput{
				ChangeBatch: &awsRoute53.ChangeBatch{
					Changes: []*awsRoute53.Change{
						{
							Action: aws.String("UPSERT"),
							ResourceRecordSet: &awsRoute53.ResourceRecordSet{
								Name:            aws.String("domain1.com"),
								Type:            aws.String("A"),
								TTL:             aws.Int64(300),
								ResourceRecords: []*awsRoute53.ResourceRecord{{Value: aws.String("1.1.1.1")}},
							},
						},
						{
							Action: aws.String("UPSERT"),
							ResourceRecordSet: &awsRoute53.ResourceRecordSet{
								Name:            aws.String("domain2.com"),
								Type:            aws.String("A"),
								TTL:             aws.Int64(300),
								ResourceRecords: []*awsRoute53.ResourceRecord{{Value: aws.String("1.1.1.1")}},
							},
						},
					},
				},
				HostedZoneId: aws.String("some-zone"),
			},
			isErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := route53.NewClient(tt.route53Client, tt.hostedZone)
			require.NoError(t, err)

			err = client.UpdateRecords(tt.domains, tt.ip)

			if tt.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedInput, tt.route53Client.input)
			}
		})
	}
}
