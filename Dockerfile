ARG PROJECT_NAME=geo-ip-info
ARG EXECUTABLE_NAME=gip

FROM golang:alpine AS builder

LABEL org.opencontainers.image.description="${PROJECT_NAME}: A tool for IP geolocation, distance calculation and country information"
LABEL org.opencontainers.image.source="https://github.com/cgiraldoz/${PROJECT_NAME}"
LABEL org.opencontainers.image.licenses=MIT

WORKDIR /build
ADD go.mod .
COPY . .

RUN go build -o ${EXECUTABLE_NAME} main.go

FROM alpine
WORKDIR /build

COPY --from=builder /build/${PROJECT_NAME} /build/${PROJECT_NAME}
CMD ["./${EXECUTABLE_NAME}"]
