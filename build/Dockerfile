FROM golang:1.20
WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY cmd ./
COPY internal ./
COPY migrations ./
COPY pkg ./
COPY tests ./
COPY .env ./
COPY test.sh ./

RUN echo pwd
RUN echo ls
RUN go build ./cmd/app -o /usr/local/bin/app

CMD ["app"]