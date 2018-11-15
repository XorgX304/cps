package properties

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"cps/pkg/kv"

	mux "github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	gjson "github.com/tidwall/gjson"
)

func init() {
	// logging
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

type Error struct {
	Status string `json:"status"`
}

func GetProperties(w http.ResponseWriter, r *http.Request, account, region string) {
	vars := mux.Vars(r)
	scope := strings.Split(vars["scope"], "/")
	service := scope[0]

	var fullPath []string
	if len(scope) <= 1 {
		w.Header().Set("Content-Type", "application/json")
		e, _ := json.Marshal(Error{
			Status: "No service provided",
		})
		w.Write(e)
	} else {
		fullPath = scope[1:len(scope)]
	}

	jsoni := kv.GetProperty(service)
	var jb []byte
	if jsoni != nil {
		jb = jsoni.([]byte)
	} else {
		return
	}

	b := new(bytes.Buffer)
	if err := json.Compact(b, jb); err != nil {
		log.Error(err)
	}

	j := []byte(b.Bytes())

	if len(fullPath) > 0 {
		f := strings.Join(fullPath, ".")
		p := gjson.GetBytes(j, "properties")
		selected := gjson.GetBytes([]byte(p.String()), f)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(strings.TrimSpace(selected.String())))
	} else {
		w.Header().Set("Content-Type", "application/json")
		p := gjson.GetBytes(j, "properties")
		w.Write([]byte(strings.TrimSpace(p.String())))
	}
}
