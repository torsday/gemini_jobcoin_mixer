FROM golang:1.12.1
LABEL maintainer="c.torsday@gmail.com"

RUN go get -u github.com/gorilla/mux \
    && go get -u github.com/spf13/cobra/cobra \
    && go get -u github.com/mattn/go-sqlite3 \
    && go get -u github.com/pkg/errors


COPY . /go/src/github.com/torsday/gemini_jobcoin_mixer
WORKDIR /go/src/github.com/torsday/gemini_jobcoin_mixer

RUN mkdir -p /go/src/github.com/torsday/gemini_jobcoin_mixer/mutable_data

RUN go build

#CMD ./redis-proxy -capacity 1000 -global-expiration 8000 -port 8080 -max-clients 5

#CMD ./gemini_jobcoin_mixer
