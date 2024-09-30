FROM golang

WORKDIR app

COPY go.mod .
COPY go.sum .

RUN go mod download 

COPY main.go .

CMD ["go", "run", "main.go"]