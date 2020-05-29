package serve

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/johndsheehan/met-eireann-archive/pkg/radar"
)

func serveRadar(rdr *radar.Radar) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr)

		g := rdr.Fetch()

		w.Header().Set("Content-Type", "image/gif")
		w.Write(g)
	}
}

// ServerConfig config for the http(s) server
type ServerConfig struct {
	ColonPort    string
	ColonTLSPort string
	FullChain    string
	PrivateKey   string
	UseTLS       bool
}

// Server struct
type Server struct {
	colonPort    string
	colonTLSPort string
	fullChain    string
	privateKey   string
	useTLS       bool

	httpSvr  *http.Server
	httpsSvr *http.Server
}

// NewServer return pointer to server struct based on config
func NewServer(cfg ServerConfig) (*Server, error) {
	svr := &Server{
		colonPort:    cfg.ColonPort,
		colonTLSPort: cfg.ColonTLSPort,
		fullChain:    cfg.FullChain,
		privateKey:   cfg.PrivateKey,
		useTLS:       cfg.UseTLS,
	}

	if svr.colonPort == "" {
		return nil, errors.New("missing http port")
	}

	if svr.useTLS {
		if svr.colonTLSPort == "" {
			return nil, errors.New("missing https port")
		}

	}
	return svr, nil
}

// Serve handle requents on http(s)
func (svr *Server) Serve(r *radar.Radar) {
	go func(r *radar.Radar) {
		http.HandleFunc("/", serveRadar(r))

		if svr.useTLS {
			svr.redirectToHTTPS()

			log.Printf("serving https on %s\n", svr.colonTLSPort)

			svr.httpsSvr = &http.Server{Addr: svr.colonTLSPort}
			err := svr.httpsSvr.ListenAndServeTLS(svr.fullChain, svr.privateKey)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Printf("serving http on %s\n", svr.colonPort)

			svr.httpSvr = &http.Server{Addr: svr.colonPort}
			err := svr.httpSvr.ListenAndServe()
			if err != nil {
				log.Print(err)
			}
		}
	}(r)
}

// Stop running server(s)
func (svr *Server) Stop() error {
	log.Print("stopping http server")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := svr.httpSvr.Shutdown(ctx)
	if err != nil {
		return err
	}

	if svr.useTLS {
		log.Print("stopping https server")

		ctxTLS, cancelTLS := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancelTLS()

		err = svr.httpsSvr.Shutdown(ctxTLS)
		if err != nil {
			return err
		}
	}

	return nil
}

func (svr *Server) redirectToHTTPS() {
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

		svr.httpSvr = &http.Server{
			Addr:    svr.colonPort,
			Handler: http.HandlerFunc(redirect),
		}

		err := svr.httpSvr.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}(svr.colonPort, svr.colonTLSPort)
}
