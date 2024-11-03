FROM golang:alpine AS builder

LABEL org.opencontainers.image.description="geo-ip-info: A tool for IP geolocation, distance calculation and country information"
LABEL org.opencontainers.image.source="https://github.com/cgiraldoz/geo-ip-info"
LABEL org.opencontainers.image.licenses=MIT

WORKDIR /build
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go build -o gip main.go

FROM alpine
WORKDIR /app

COPY --from=builder /build/gip /app/gip

COPY GeoLite2-City.mmdb /app/GeoLite2-City.mmdb
COPY config.yaml /app/config.yaml

CMD ["./gip", "api"]
