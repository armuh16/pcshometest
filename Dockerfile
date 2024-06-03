FROM golang:alpine

RUN apk update && apk add --no-cache git

ENV TZ="Asia/Jakarta"

WORKDIR /app/

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/app

EXPOSE 8081

COPY config.example.yml config.yml

ENTRYPOINT ["/go/bin/app","start"]
