# builder
FROM golang:alpine AS build-env
LABEL maintainer="Yannick Foeillet <bzhtux@gmail.com>"

# wokeignore:rule=he/him/his
RUN apk --no-cache add build-base git mercurial gcc curl
RUN mkdir -p /go/src/github.com/bzhtux/postgres
ADD . /go/src/github.com/bzhtux/postgres
RUN cd /go/src/github.com/bzhtux/postgres && go get ./... && go build -o postgres-app cmd/main.go
WORKDIR /go/src/github.com/bzhtux/postgres
RUN go get ./...
RUN go build -o postgres-app cmd/main.go


# final image
FROM alpine
LABEL maintainer="Yannick Foeillet <bzhtux@gmail.com>"

# wokeignore:rule=he/him/his
RUN apk --no-cache add curl jq
RUN adduser -h /app -s /bin/sh -u 1000 -D app
RUN mkdir -p /config
COPY config/postgres.yaml /config/
WORKDIR /app
COPY --from=build-env /go/src/github.com/bzhtux/postgres/postgres-app /app/
USER 1000
ENTRYPOINT ./postgres-app

