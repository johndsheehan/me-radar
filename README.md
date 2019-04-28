Fetch the [Met Eireann](http://archive.met.ie/) latest rainfall radar image.

Fetch previous 10 rainfall radar images, and combine them into a gif.
The gif is served on localhost, to the user-specified port.

## Build
```
go get
go build
```

## Usage

```
./met-eireann-archive --port <port-number>
```

## Docker
```
docker build -t met:latest .
docker run --rm -it -p 3031:3031 met:latest
```

## View
in browser, go to `http://localhost:3031`