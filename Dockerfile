FROM golang:1.18 as builder

WORKDIR /app
COPY . .
RUN go get ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -o deployment/bin/app main.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/deployment/bin/app /opt/bin/app

EXPOSE 8080

ENTRYPOINT ["/opt/bin/app"]