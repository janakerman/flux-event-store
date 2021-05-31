FROM golang:1.16 as builder

WORKDIR /code
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /api /code/cmd/api/api.go

FROM scratch

COPY --from=builder /api /api

ENTRYPOINT ["/api"]