package route53

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	awsRoute53 "github.com/aws/aws-sdk-go/service/route53"

	log "github.com/sirupsen/logrus"
)

type Route53Client interface {
	ChangeResourceRecordSets(input *awsRoute53.ChangeResourceRecordSetsInput) (*awsRoute53.ChangeResourceRecordSetsOutput, error)
}

type Client struct {
	awsRoute53Client Route53Client
	hostedZoneId     string
}

func NewClient(awsRoute53Client Route53Client, hostedZoneId string) (*Client, error) {
	if awsRoute53Client == nil {
		return nil, fmt.Errorf("route53Client not initialised, awsRoute53Client not found")
	}

	c := &Client{
		awsRoute53Client: awsRoute53Client,
		hostedZoneId:     hostedZoneId,
	}

	return c, nil
}

func (c *Client) UpdateRecords(domains []string, ip string) error {
	changes := []*awsRoute53.Change{}

	for _, domain := range domains {
		changes = append(changes, &awsRoute53.Change{
			Action: aws.String("UPSERT"),
			ResourceRecordSet: &awsRoute53.ResourceRecordSet{
				Name:            aws.String(domain),
				Type:            aws.String("A"),
				TTL:             aws.Int64(300),
				ResourceRecords: []*awsRoute53.ResourceRecord{{Value: aws.String(ip)}},
			},
		})
	}

	input := &awsRoute53.ChangeResourceRecordSetsInput{
		ChangeBatch: &awsRoute53.ChangeBatch{
			Changes: changes,
		},
		HostedZoneId: aws.String(c.hostedZoneId),
	}

	_, err := c.awsRoute53Client.ChangeResourceRecordSets(input)
	if err != nil {
		return err
	}

	log.Printf("Successfully updated domains (%s) with ip (%s)", strings.Join(domains, ","), ip)
	return nil
}
