FROM golang:1.16 as builder

WORKDIR /code
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /controller /code/cmd/controller/controller.go

FROM scratch

COPY --from=builder /controller /controller

ENTRYPOINT ["/api"]