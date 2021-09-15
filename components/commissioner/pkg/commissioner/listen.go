package commissioner

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pacedotdev/oto/otohttp"

	"github.com/thavlik/bvs/components/commissioner/pkg/api"
)

func (s *Server) listen(port int) error {
	otoServer := otohttp.NewServer()
	api.RegisterCommissioner(otoServer, s)
	return (&http.Server{
		Handler:      otoServer,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}).ListenAndServe()
}
