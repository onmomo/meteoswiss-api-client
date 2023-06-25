FROM golang:1.20.5 AS builder
# smoke test to verify if golang is available
RUN go version

ARG PROJECT_VERSION

COPY . /go/src/github.com/onmomo/meteoswiss-api-client/
WORKDIR /go/src/github.com/onmomo/meteoswiss-api-client/
RUN set -Eeux && \
    go mod download && \
    go mod verify

RUN GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-w -s -X 'main.Version=${PROJECT_VERSION}'" \
    -o app meteoApiClient.go
RUN go test -cover -v ./...

FROM alpine:3.17.1
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/onmomo/meteoswiss-api-client/app .

ENTRYPOINT ["./app"]