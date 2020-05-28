package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	mea "github.com/johndsheehan/met-eireann-archive/pkg/met-eireann-archive"
	"github.com/johndsheehan/met-eireann-archive/pkg/radar"
)

func main() {
	port := flag.String("port", "3080", "http port (default 3080)")
	tlsport := flag.String("tlsport", "3443", "https port (default 3443)")
	fullchain := flag.String("fullchain", "", "fullchain.pem")
	privateKey := flag.String("privateKey", "", "privKey.pem")

	flag.Parse()

	useTLS := false
	if *fullchain != "" && *privateKey != "" {
		useTLS = true
	}

	serverCfg := ServerConfig{
		colonPort:    ":" + *port,
		colonTLSPort: ":" + *tlsport,
		fullchain:    *fullchain,
		privateKey:   *privateKey,
		useTLS:       useTLS,
	}

	mea, err := mea.NewMEArchive(&mea.MEArchiveConfig{})
	if err != nil {
		log.Fatal(err)
	}

	r := radar.NewRadar(10, mea)

	r.Watch()
	go serve(r, serverCfg)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	r.Stop()
}
