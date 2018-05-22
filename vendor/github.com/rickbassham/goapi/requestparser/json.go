package requestparser

import (
	"encoding/json"
	"net/http"
)

func bodyJson(r *http.Request, item interface{}) (err error) {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(item)
	return
}
