package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

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

	colonPort := ":3031"

	go update(r)
	go serve(r, colonPort)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
