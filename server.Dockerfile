FROM golang:1.23.3

WORKDIR /go-keeper

COPY go.mod go.sum ./

COPY ./cmd/server/ ./cmd/server/
COPY ./internal/server/ ./internal/server/
COPY ./internal/common/ ./internal/common/
COPY ./pkg/ ./pkg/

RUN go mod tidy
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main/main.go

EXPOSE 8080

CMD ["./server"]