#use  golang:1.18-alpine
FROM golang:1.18-alpine3.16 AS build
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go build .

FROM alpine 
RUN mkdir /app
WORKDIR /app
COPY --from=build /app/drone-keda-scaler .
ENTRYPOINT "./drone-keda-scaler"
