FROM golang:1.15 as builder

LABEL VERSION=v0.11.0 \
      ARCH=AMD64 \
      DESCRIPTION="A monitoring of prometheus for druid" \
      MAINTAINER="OpsTree Solutions"

WORKDIR /go/src/druid-exporter

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download
COPY . /go/src/druid-exporter
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o druid-exporter main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /go/src/druid-exporter/druid-exporter .
USER nonroot:nonroot
ENTRYPOINT ["/druid-exporter"]
