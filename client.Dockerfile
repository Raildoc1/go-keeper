FROM golang:1.23.3 AS build-stage

WORKDIR /go-keeper

COPY go.mod go.sum ./

COPY ./cmd/client/ ./cmd/client/
COPY ./internal/client/ ./internal/client/
COPY ./internal/common/ ./internal/common/
COPY ./pkg/ ./pkg/

RUN go mod tidy
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o client ./cmd/client/main/main.go

FROM ubuntu:22.04 AS release-stage

WORKDIR /

COPY --from=build-stage ./go-keeper/client ./client

ENTRYPOINT ["./client"]