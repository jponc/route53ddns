# route53ddns

This package is built with Go that updates your Route53 DNS entires using your IP address. This is useful for people who has homelabs that has dynamic ip address.

# Usage

```
docker run --rm \
    -e AWS_ACCESS_KEY_ID=<your_aws_access_key_id> \
    -e AWS_SECRET_ACCESS_KEY=<your_aws_secret_access_key> \
    -e AWS_REGION=global \
    -e ROUTE53_HOSTED_ZONE_ID=<your_aws_hosted_zone_id> \
    -e ROUTE53_DOMAINS=domain1,domain2 \
    julianponce/route53ddns
```
