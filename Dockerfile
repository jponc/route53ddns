FROM golang:1.19-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o app cmd/route53ddns/main.go cmd/route53ddns/config.go

FROM amd64/alpine:3.16.2 AS final
LABEL maintainer="ponce.julianalfonso@gmail.com"
COPY --from=build /app/app /app

CMD [ "/app" ]
