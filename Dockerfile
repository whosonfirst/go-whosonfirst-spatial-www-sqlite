FROM golang:1.16-alpine

RUN mkdir /build

COPY . /build/go-whosonfirst-spatial-www-sqlite

RUN apk update && apk upgrade \
    && apk add libc-dev gcc \
    && cd /build/go-whosonfirst-spatial-www-sqlite \
    && go build -mod vendor -o /main cmd/server/main.go \    
    && cd && rm -rf /build

RUN mkdir /usr/local/data
COPY whosonfirst.db /usr/local/data/whosonfirst.db

ADD https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie /usr/bin/aws-lambda-rie
RUN chmod 755 /usr/bin/aws-lambda-rie
COPY Dockerfile.entry.sh /entry.sh
RUN chmod 755 /entry.sh

ENTRYPOINT [ "/entry.sh" ] 