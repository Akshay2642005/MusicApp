FROM golang:1.23

WORKDIR /app

COPY . .

RUN go get -d -v ./

RUN go build -o app .


EXPOSE 8000

CMD ["./app"]
