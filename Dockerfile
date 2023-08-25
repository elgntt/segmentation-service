FROM golang:1.21-alpine

WORKDIR /go/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./bin/app ./cmd/app/main.go

EXPOSE 8080

CMD [ "./bin/app" ]
