FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.20.5 AS builder

# smoke test to verify if golang is available
RUN go version

ARG PROJECT_VERSION

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

COPY . /go/src/github.com/onmomo/meteoswiss-api-client/
WORKDIR /go/src/github.com/onmomo/meteoswiss-api-client/

RUN set -Eeux && \
    go mod download && \
    go mod verify

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
    -trimpath \
    -ldflags="-w -s -X 'main.Version=${PROJECT_VERSION}'" \
    -o app meteoApiClient.go
RUN go test -cover -v ./...

FROM --platform=${TARGETPLATFORM:-linux/amd64} alpine:3.17.1
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/onmomo/meteoswiss-api-client/app .

ENTRYPOINT ["./app"]