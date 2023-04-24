FROM golang:1.20.3-alpine3.16

WORKDIR /app

COPY .env ./
COPY go.mod ./ go.sum ./
RUN go mod download

COPY *.go ./

RUN go build .

CMD [ "./telegram-gpt" ]