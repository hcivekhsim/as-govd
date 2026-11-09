FROM golang:1.27.3

ENV GOPATH=/ 

WORKDIR /web-app 

COPY . .

RUN apt-get update
RUN apt-get -y install postgresql-client

RUN chmod +x ./scripts/wait-for-postgres.sh

RUN go mod download && go build -o as-app ./cmd/main.go

CMD ["./as-app"]
