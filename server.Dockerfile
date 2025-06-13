FROM golang:1.23.3 AS build-stage

WORKDIR /go-keeper

COPY go.mod go.sum ./

COPY ./cmd/server/ ./cmd/server/
COPY ./internal/server/ ./internal/server/
COPY ./internal/common/ ./internal/common/
COPY ./pkg/ ./pkg/

RUN go mod tidy
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main/main.go

FROM ubuntu:22.04 AS release-stage

WORKDIR /

COPY --from=build-stage ./go-keeper/server ./server

EXPOSE 8080

ENTRYPOINT ["./server"]