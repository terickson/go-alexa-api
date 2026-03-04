FROM golang:1.25 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-alexa-api .

FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/go-alexa-api /
EXPOSE 8000
ENTRYPOINT ["/go-alexa-api"]
