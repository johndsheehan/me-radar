package main

import (
	"log"
	"net/http"
)

func serveRadar(rdr *Radar) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr)

		g := rdr.Fetch()

		w.Header().Set("Content-Type", "image/gif")
		w.Write(g)
	}
}
func serve(r *Radar, colonPort string) {
	http.HandleFunc("/", serveRadar(r))
	err := http.ListenAndServe(colonPort, nil)
	if err != nil {
		log.Print(err)
	}
}
