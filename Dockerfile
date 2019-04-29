FROM golang:alpine AS build

RUN apk update && apk add  --no-cache  git tzdata

WORKDIR /app
COPY go.mod go.sum *.go /app/

ENV GO111MODULE="on"
RUN go get 
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' .


FROM scratch
COPY --from=build /app/met-eireann-archive /met-eireann-archive
COPY --from=build /usr/share/zoneinfo/Eire /etc/localtime

ENTRYPOINT ["/met-eireann-archive"]
