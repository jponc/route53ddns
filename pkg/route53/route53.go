package route53

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awsRoute53 "github.com/aws/aws-sdk-go/service/route53"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	awsRoute53Client *awsRoute53.Route53
	hostedZoneId     string
}

func NewClient(awsRegion, awsAccessKeyId, awsSecretAccessKey, hostedZoneId string) (*Client, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyId, awsSecretAccessKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create aws session: %v", err)
	}

	awsRoute53Client := awsRoute53.New(sess)

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
