FROM golang:latest as builder
MAINTAINER OpsTree Solutions
COPY ./ /go/src/druid-exporter/
WORKDIR /go/src/druid-exporter/
RUN go get -v -t -d ./... \
    && go build -o druid-exporter

FROM alpine:latest
MAINTAINER OpsTree Solutions
WORKDIR /app
RUN apk add --no-cache libc6-compat
COPY --from=builder /go/src/druid-exporter/druid-exporter /app/
COPY --from=builder /go/src/druid-exporter/* /app/
ENTRYPOINT ["./druid-exporter"]
