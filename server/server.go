package server

import "net/http"

func NewServer() {
	var srv := &http.Server{
		Addr:                         "",
		Handler:                      nil,
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		ReadTimeout:                  0,
		ReadHeaderTimeout:            0,
		WriteTimeout:                 0,
		IdleTimeout:                  0,
		MaxHeaderBytes:               0,
		TLSNextProto:                 nil,
		ConnState:                    nil,
		ErrorLog:                     nil,
		BaseContext:                  nil,
		ConnContext:                  nil,
	}
	err := srv.ListenAndServe()
	if err != nil {
		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/login", login)
	srv.Handler = mux

}
