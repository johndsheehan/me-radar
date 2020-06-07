package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	mer "github.com/johndsheehan/met-eireann-archive/pkg/met-eireann-radar"
	"github.com/johndsheehan/met-eireann-archive/pkg/radar"
	"github.com/johndsheehan/met-eireann-archive/pkg/serve"
)

func main() {
	format := flag.String("format", "archive", "image format `archive` or `current`")
	port := flag.String("port", "3080", "http port (default 3080)")
	tlsport := flag.String("tlsport", "3443", "https port (default 3443)")
	fullchain := flag.String("fullchain", "", "fullchain.pem")
	privateKey := flag.String("privateKey", "", "privKey.pem")

	flag.Parse()

	useTLS := false
	if *fullchain != "" && *privateKey != "" {
		useTLS = true
	}

	serverCfg := serve.ServerConfig{
		ColonPort:    ":" + *port,
		ColonTLSPort: ":" + *tlsport,
		FullChain:    *fullchain,
		PrivateKey:   *privateKey,
		UseTLS:       useTLS,
	}

	svr, err := serve.NewServer(serverCfg)
	if err != nil {
		log.Fatal(err)
	}

	rfmt := mer.ARCHIVE
	if *format == "current" {
		rfmt = mer.CURRENT
	}

	mer, err := mer.NewMERadar(rfmt)
	if err != nil {
		log.Fatal(err)
	}

	rdr := radar.NewRadar(10, mer)

	rdr.Watch()
	svr.Serve(rdr)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	rdr.Stop()
	svr.Stop()
}
