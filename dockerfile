FROM golang:1.23 AS build

WORKDIR /rtt

COPY go.mod .
COPY go.sum .

RUN go mod download 

COPY main.go .

RUN go build

FROM alpine

COPY --from=build /rtt/gone-rtt /gone-rtt

RUN apk add --no-cache gcompat

CMD ["/gone-rtt"]
