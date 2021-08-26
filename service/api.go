package service

import (
	"net/http"

	"github.com/HiN3i1/wallet-service/router"
)

// NewAPIServer a api server
func NewAPIServer(port string) *http.Server {
	server := &http.Server{
		Addr:    port,
		Handler: router.New(),
	}
	return server
}
