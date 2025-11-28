FROM golang:1.25.1-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY assets/ assets/
COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/

RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api/main.go

FROM alpine:3.22

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/api .

RUN adduser -D appuser
USER appuser

EXPOSE 8080

CMD ["./api"]

