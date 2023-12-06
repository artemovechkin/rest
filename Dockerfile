FROM golang:latest

ENV GIN_MODE=debug
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /simpleopc ./cmd

EXPOSE 8000

CMD [ "/simpleopc" ]

