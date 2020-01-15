package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	r := NewRadar(10)

	for i := 11; i > 0; i-- {
		d := time.Duration(i * 15)
		then := time.Now().Add(-d * time.Minute)
		gifImg, err := fetch(then)
		if err != nil {
			log.Print(err)
			continue
		}
		r.Update(gifImg)
	}

	go update(r)
	go serve(r, serverCfg)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
