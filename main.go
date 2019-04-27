package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	port := flag.Int("port", 3031, "server port (default 3031)")
	flag.Parse()

	if *port < 1024 || *port > 65535 {
		log.Fatal(errors.New("port should be between 1024 and 65536"))
	}
	colonPort := fmt.Sprintf(":%d", *port)

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
	go serve(r, colonPort)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
