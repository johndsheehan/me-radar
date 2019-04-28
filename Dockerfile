FROM ubuntu:18.04


ENV DEBIAN_FRONTEND="noninteractive"
RUN apt-get update  && \
    apt-get install  -y  build-essential git sudo tzdata wget  && \
    ln  -fs  /usr/share/zoneinfo/Europe/Dublin  /etc/localtime  && \
    dpkg-reconfigure --frontend noninteractive tzdata


ENV PATH="$PATH:/go/bin"
ENV GOPATH="/golang"
RUN wget -q https://dl.google.com/go/go1.12.4.linux-amd64.tar.gz  && \
    tar -xf go1.12.4.linux-amd64.tar.gz  && \
    mkdir /golang  && \
    go get -u -d gocv.io/x/gocv  && \
    cd $GOPATH/src/gocv.io/x/gocv  && \
    make install


WORKDIR /met
COPY go.mod go.sum *.go /met/
RUN  go get  && \
     go build


CMD ["./met-eireann-archive"]
