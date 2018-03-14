package plug

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/krhancoc/frud/config"
	"github.com/krhancoc/frud/database"
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

func paramsToVal(b []byte) map[string]string {

	var objmap map[string]*json.RawMessage
	json.Unmarshal(b, &objmap)

	m := make(map[string]string, len(objmap))
	for key, value := range objmap {
		v := (*value)[1 : len(*value)-1]
		m[key] = string(v)
	}
	return m

}

func queryToVal(q url.Values, plugs []*config.Field) map[string]string {

	vals := make(map[string]string, len(q))
	for _, plug := range plugs {
		if val := q.Get(plug.Key); val != "" {
			vals[plug.Key] = val
		}
	}
	return vals
}

func get(w http.ResponseWriter, req *http.Request, ctx config.AppContext, plug Plug) {

	m := queryToVal(req.URL.Query(), plug.Model)
	log.WithFields(log.Fields{
		"method": req.Method,
		"object": plug.Name,
		"query":  m,
	}).Info("Post request received")

	dbReq := &config.DBRequest{
		Method: "get",
		Values: m,
		Type:   plug.Name,
		Model:  plug.Model,
	}
	result, err := ctx.Driver.MakeRequest(dbReq)
	if err != nil {
		e := err.(database.DriverError)
		log.Error(err.Error())
		ctx.Render.JSON(w, e.Status, e)
		return
	}
	ctx.Render.JSON(w, http.StatusAccepted, result)
	return
}

func post(w http.ResponseWriter, req *http.Request, ctx config.AppContext, plug Plug) {

	b, _ := ioutil.ReadAll(req.Body)
	m := paramsToVal(b)

	log.WithFields(log.Fields{
		"method": req.Method,
		"object": plug.Name,
		"params": string(b),
	}).Info("Post request received")

	dbReq := &config.DBRequest{
		Method: "post",
		Values: m,
		Type:   plug.Name,
		Model:  plug.Model,
	}
	_, err := ctx.Driver.MakeRequest(dbReq)
	if err != nil {
		e := err.(database.DriverError)
		log.Error(err.Error())
		ctx.Render.JSON(w, e.Status, e)
		return
	}
	ctx.Render.JSON(w, http.StatusCreated, Message{
		Status:  http.StatusCreated,
		Message: "Created",
	})
	return
}

func delete(w http.ResponseWriter, req *http.Request, ctx config.AppContext, plug Plug) {

	ctx.Render.Text(w, 200, "HELLO")
}

func put(w http.ResponseWriter, req *http.Request, ctx config.AppContext, plug Plug) {

	ctx.Render.Text(w, 200, "HELLO")
}
