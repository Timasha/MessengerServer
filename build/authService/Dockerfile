FROM golang:latest as dependencies

RUN apt-get -y update
RUN apt-get install -y ca-certificates openssl
ARG cert_location=/usr/local/share/ca-certificates

WORKDIR /dependencies

COPY ./go.mod ./go.sum ./
RUN go mod download

RUN mkdir /tls && mkdir /configs

COPY ./tls_certs/services/authorization/* /tls/
COPY ./configs/* /configs/

FROM dependencies as build

WORKDIR /build
COPY . .
RUN go build ./cmd/services/authorization/main.go

FROM debian
WORKDIR /service

RUN mkdir /tls /configs /migrations

COPY --from=build /build/main .
COPY --from=build /build/migrations/* /migrations/
COPY --from=dependencies /tls/* /tls/
COPY --from=dependencies /configs/authServiceConfig.json /configs
RUN ls /

ENTRYPOINT ["./main"]