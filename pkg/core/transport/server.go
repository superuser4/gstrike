package transport

import (
	"crypto/tls"
	"gstrike/pkg/config"
	"gstrike/pkg/core/commgr"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)


type GStrikeServer struct {
	Config config.ServerConfig
	HttpServer *HttpsListener
}

func NewGStrike(conf config.ServerConfig) (*GStrikeServer, error) {
	s, err := NewHttps(conf.Port)
	if err != nil {
		return nil, err
	}
	strike := GStrikeServer{Config: conf, HttpServer: s}
	return &strike, nil
}

func (strike GStrikeServer) Start() error {	
	r := mux.NewRouter()
	r.HandleFunc("/ws", commgr.WsHandler)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	strike.HttpServer.Server = &http.Server{
		Addr:      ":" + strconv.Itoa(strike.Config.Port),
		Handler:   r,
		TLSConfig: tlsConfig,
	}

	strike.HttpServer.StartedAt = time.Now()
	strike.HttpServer.Status = running
	err := strike.HttpServer.Server.ListenAndServeTLS(strike.HttpServer.CertFile, strike.HttpServer.KeyFile)
	if err != nil {
		return err
	}
	return nil
}

func (strike GStrikeServer) Stop() error {
	err := strike.HttpServer.Stop()
	return err
}
