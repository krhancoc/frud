package plug

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/krhancoc/frud/config"
	"github.com/krhancoc/frud/errors"
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

func paramsToVal(b []byte) map[string]interface{} {

	var objmap map[string]*json.RawMessage
	json.Unmarshal(b, &objmap)

	m := make(map[string]interface{}, len(objmap))
	for key, value := range objmap {
		b, _ := value.MarshalJSON()
		if b[0] == '"' && b[len(b)-1] == '"' {
			b = b[1 : len(b)-1]
			m[key] = string(b)
		} else {
			m[key] = paramsToVal(b)
		}
	}
	return m

}

func queryToVal(q url.Values, plugs []*config.Field) map[string]interface{} {

	vals := make(map[string]interface{}, len(q))
	for _, plug := range plugs {
		if val := q.Get(plug.Key); val != "" {
			vals[plug.Key] = val
		}
	}
	return vals
}

func generic(w http.ResponseWriter, req *http.Request, ctx config.AppContext, plug Plug) {

	req.Method = strings.ToLower(req.Method)
	b, _ := ioutil.ReadAll(req.Body)
	params := paramsToVal(b)
	queries := queryToVal(req.URL.Query(), plug.Model)
	log.WithFields(log.Fields{
		"method": req.Method,
		"object": plug.Name,
		"params": string(b),
		"query":  queries,
	}).Info("Post request received")

	dbReq := &config.DBRequest{
		Method:  req.Method,
		Params:  params,
		Queries: queries,
		Type:    plug.Name,
		Model:   plug.Model,
	}
	err := dbReq.Validate()
	if err != nil {
		ctx.Render.JSON(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := ctx.Driver.MakeRequest(dbReq)
	if err != nil {
		e := err.(errors.DriverError)
		log.Error(err.Error())
		ctx.Render.JSON(w, e.Status, e)
		return
	}
	switch strings.ToLower(req.Method) {
	case "put", "post":
		ctx.Render.JSON(w, http.StatusCreated, Message{
			Status:  http.StatusCreated,
			Message: "Created",
			Results: result,
		})
	case "delete":
		ctx.Render.JSON(w, http.StatusOK, Message{
			Status:  http.StatusCreated,
			Message: "Created",
			Results: result,
		})
	case "get":
		ctx.Render.JSON(w, http.StatusOK, result)
	}
	return
}
