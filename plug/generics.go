package plug

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/krhancoc/frud/config"
	log "github.com/sirupsen/logrus"
)

type healthCheck struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Version string `json:"version"`
}

// HealthCheck is the generic health check endpoint for the user, instantianted when
// called for in the config file.
func HealthCheck(w http.ResponseWriter, req *http.Request, ctx config.AppContext) {

	ctx.Render.JSON(w, http.StatusOK, healthCheck{
		Status:  200,
		Message: "Status good",
		Version: ctx.Version,
	})
}

type genericHandler func(http.ResponseWriter, *http.Request, config.AppContext, Plug)

func makeGenericHandler(ctx config.AppContext, plug Plug, fn genericHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, ctx, plug)
	}
}

func get(w http.ResponseWriter, req *http.Request, ctx config.AppContext, plug Plug) {

	ctx.Render.JSON(w, 200, plug)
}

func post(w http.ResponseWriter, req *http.Request, ctx config.AppContext, plug Plug) {

	b, _ := ioutil.ReadAll(req.Body)
	log.WithFields(log.Fields{
		"method": req.Method,
		"object": plug.Name,
		"query":  string(b),
	}).Info("Post request received")
	var objmap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objmap)
	if err != nil {
		log.Error(err.Error())
		return
	}

	m := make(map[string]string, len(objmap))
	for key, value := range objmap {
		v := (*value)[1 : len(*value)-1]
		m[key] = string(v)
	}
	dbReq := &config.DBRequest{
		Method: "post",
		Values: m,
		Type:   plug.Name,
		Model:  plug.Model,
	}
	err = ctx.Driver.MakeRequest(dbReq)
	if err != nil {
		log.Error(err.Error())
		ctx.Render.Text(w, 500, err.Error())
		return
	}
	ctx.Render.Text(w, 200, "POST REQUEST")
	return
}

func delete(w http.ResponseWriter, req *http.Request, ctx config.AppContext, plug Plug) {

	ctx.Render.Text(w, 200, "HELLO")
}

func put(w http.ResponseWriter, req *http.Request, ctx config.AppContext, plug Plug) {

	ctx.Render.Text(w, 200, "HELLO")
}
