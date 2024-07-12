FROM golang:1.22.5-alpine

WORKDIR /nlwj

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

WORKDIR /nlwj/cmd/nlwj

RUN go build -o /nlwj/bin/nlwj .

EXPOSE 8080
ENTRYPOINT [ "/nlwj/bin/nlwj" ]
