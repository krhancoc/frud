package defaults

import (
	"net/http"

	"github.com/krhancoc/frud/config"
)

func Get(w http.ResponseWriter, req *http.Request, ctx config.AppContext) {
	println("HELLO FROM INSIDE OF DEFAULT")
	ctx.Render.Text(w, 200, "DEFAULT")
}
