FROM golang:1.20

RUN apt-get install -y wget
RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - \
    && echo "deb http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google.list
RUN apt-get update && apt-get -y install google-chrome-stable

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY cmd ./cmd
COPY internal ./internal
COPY migrations ./migrations
COPY pkg ./pkg
COPY tests ./tests
COPY .env ./.env
COPY test.sh ./test.sh

RUN go build -o app ./cmd/app

CMD ["/usr/src/app/app"]