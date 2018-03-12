package plug

import (
	"net/http"

	"github.com/krhancoc/frud/config"
)

type healthCheck struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Version string `json:"version"`
}

func HealthCheck(w http.ResponseWriter, req *http.Request, ctx config.AppContext) {

	ctx.Render.JSON(w, http.StatusOK, healthCheck{
		Status:  200,
		Message: "Status good",
		Version: ctx.Version,
	})
}
