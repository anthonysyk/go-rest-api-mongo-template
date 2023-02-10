FROM golang:alpine AS builder
WORKDIR /go/src/app
COPY . .
RUN apk add git
RUN CGO_ENABLED=0 go install -ldflags '-extldflags "-static"' -tags timetzdata

FROM scratch
COPY --from=builder /go/bin/go-rest-api-mongo-template /go-rest-api-mongo-template
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/go-rest-api-mongo-template"]