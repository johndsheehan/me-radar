FROM golang:1.14-alpine AS build

RUN apk update  \
 && apk add  --no-cache  ca-certificates  git  tzdata  upx  \
 && update-ca-certificates  \
 && addgroup  --system  app  \
 && adduser  -S -G  app  app

WORKDIR /app
COPY go.mod go.sum *.go  ./
COPY pkg  ./pkg

RUN go get -u  \
 && CGO_ENABLED=0 GOOS=linux go build  -a  -ldflags '-extldflags "-static"'  -o /tmp/app  .  \
 && upx  --best  /tmp/app  \
 && upx  -t  /tmp/app

FROM scratch

COPY --from=build  /usr/share/zoneinfo  /usr/share/zoneinfo
COPY --from=build  /etc/ssl/certs/ca-certificates.crt  /etc/ssl/certs/
COPY --from=build  /etc/passwd  /etc/passwd
COPY --from=build  /tmp/app  /home/app/app

ENV TZ=Europe/Dublin
EXPOSE 3031/tcp

USER app

ENTRYPOINT ["/home/app/app"]
