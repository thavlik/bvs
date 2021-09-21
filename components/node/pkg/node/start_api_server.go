package node

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pacedotdev/oto/otohttp"
	"github.com/thavlik/bvs/components/node/pkg/api"
)

func (s *Server) startAPIServer(port int) error {
	otoServer := otohttp.NewServer()
	api.RegisterNode(otoServer, s)
	return (&http.Server{
		Handler:      otoServer,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}).ListenAndServe()
}
