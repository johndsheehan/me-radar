package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func serveRadar(rdr *Radar) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr)

		g := rdr.Fetch()

		w.Header().Set("Content-Type", "image/gif")
		w.Write(g)
	}
}

func redirectToHTTPS(colonPort, colonTLSPort string) {
	go func(colonPort, colonTLSPort string) {
		redirect := func(w http.ResponseWriter, r *http.Request) {
			log.Printf("http request from: %s %s %s\n", r.RemoteAddr, r.Method, r.URL)

			host := r.Host
			if strings.HasSuffix(r.Host, colonPort) {
				host = strings.TrimSuffix(host, colonPort)
			}

			url := fmt.Sprintf("https://%s%s/%s", host, colonTLSPort, r.RequestURI)
			log.Printf("redirect to: %s", url)

			http.Redirect(w, r, url, http.StatusMovedPermanently)
		}

		err := http.ListenAndServe(colonPort, http.HandlerFunc(redirect))
		if err != nil {
			log.Fatal(err)
		}
	}(colonPort, colonTLSPort)
}

type ServerConfig struct {
	colonPort    string
	colonTLSPort string
	fullchain    string
	privateKey   string
	useTLS       bool
}

func serve(r *Radar, cfg ServerConfig) {
	http.HandleFunc("/", serveRadar(r))

	if cfg.useTLS {
		go redirectToHTTPS(cfg.colonPort, cfg.colonTLSPort)

		log.Printf("serving https on %s\n", cfg.colonTLSPort)
		err := http.ListenAndServeTLS(cfg.colonTLSPort, cfg.fullchain, cfg.privateKey, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Printf("serving http on %s\n", cfg.colonPort)

		err := http.ListenAndServe(cfg.colonPort, nil)
		if err != nil {
			log.Print(err)
		}
	}
}
