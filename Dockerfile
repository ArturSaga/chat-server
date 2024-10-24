FROM golang:1.23.2-alpine AS builder

COPY . /github.com/ArturSaga/chat-server/source/
WORKDIR /github.com/ArturSaga/chat-server/source/

RUN go mod download
RUN go build -o ./bin/chat-server cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/ArturSaga/chat-server/source/bin/chat-server .

CMD ["./chat-server"]