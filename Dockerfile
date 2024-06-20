FROM golang:alpine

WORKDIR /app

COPY . .

RUN go build -o app /app/cmd/app

EXPOSE 8080

ENTRYPOINT ["/app"]