FROM  golang:1.21 AS builder
ADD . /src
WORKDIR /src
RUN go mod download
RUN go build -o backend cmd/backend/main.go

FROM debian:trixie
RUN apt-get -y update && apt-get -y upgrade && apt-get install -y ca-certificates

COPY --from=builder /src/backend /app/backend
COPY --from=builder /src/cert.pem /app/cert.pem
COPY --from=builder /src/key.pem /app/key.pem
COPY --from=builder /src/config.ini /app/config.ini

RUN chmod +x /app/backend
WORKDIR /app
ENTRYPOINT ["./backend"]
