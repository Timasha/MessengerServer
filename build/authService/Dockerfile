FROM ubuntu
RUN apt-get -y update
RUN apt -y install golang
RUN apt-get install -y ca-certificates openssl
ARG cert_location=/usr/local/share/ca-certificates

WORKDIR /build
COPY . .
RUN go mod download
RUN go build ./cmd/services/authorization/main.go

WORKDIR /service

RUN cp /build/main .
RUN cp -r /build/migrations /
RUN rm -r /build

RUN mkdir /configs
COPY configs/authServiceConfig.json /configs/
ENV AUTHSERVICE_PATH_CONFIG="/configs/authServiceConfig.json"
ENV AUTHSERVICE_PATH_MIGRATIONS="/migrations/"