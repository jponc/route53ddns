build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/route53ddns cmd/route53ddns/*.go

build_image:
	docker build -t julianponce/route53ddns:latest -t julianponce/route53ddns:${TAG} .

push_image:
	docker push julianponce/route53ddns:${TAG}
	docker push julianponce/route53ddns:latest
