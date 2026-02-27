FROM docker.io/library/golang:1.25 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app

FROM scratch
COPY --from=builder /app /app
ENTRYPOINT ["/app"]