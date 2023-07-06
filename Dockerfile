FROM golang:1.16-alpine

WORKDIR /source

COPY . .

RUN go build /source/cmd/server

RUN mkdir -p /MessengerServer

WORKDIR /MessengerServer

RUN mv /source/server /MessengerServer 

RUN mv /source/web/template/registrationForm.html /MessengerServer

EXPOSE 8080