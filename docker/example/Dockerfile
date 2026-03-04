FROM docker.io/library/golang:1.25 AS builder
ARG foo=default
RUN echo $foo > /tmp/foo
RUN --mount=type=secret,id=hidden,target=/tmp/hidden <<EOF
cp -r /tmp/hidden /tmp/revealed || true
EOF

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app

FROM scratch
COPY --from=builder /app /app
ENTRYPOINT ["/app"]
