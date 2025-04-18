package handlers

import (
	"github.com/gorilla/mux"
	"github.com/grigory222/avito-backend-trainee/config"
	"net/http"
	"strconv"
)

func NewServer(cfg *config.Config) *http.Server {
	r := mux.NewRouter()

	srv := &http.Server{
		Addr:    cfg.HTTP.Host + ":" + strconv.Itoa(cfg.HTTP.Port),
		Handler: r}
	return srv
}
