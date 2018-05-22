package requestparser

import (
	"net/http"
	"reflect"

	"github.com/go-chi/chi"
)

type bodyer interface {
	Body() interface{}
}

type former interface {
	Form() interface{}
}

type router interface {
	Route() interface{}
}

type validater interface {
	Validate() error
}

func ParseRequest(r *http.Request, req interface{}) error {
	if reqb, ok := req.(bodyer); ok {
		body := reqb.Body()
		if body != nil && r.Header.Get("Content-Type") == "application/json" {
			err := bodyJson(r, body)
			if err != nil {
				return err
			}
		}
	}

	if reqf, ok := req.(former); ok {
		// This will parse the query string and form.
		form := reqf.Form()
		if form != nil {
			err := r.ParseForm()
			if err != nil {
				return err
			}

			fvalues := r.Form

			vp := reflect.ValueOf(form)
			v := vp.Elem()
			t := v.Type()

			fields := t.NumField()

			for i := 0; i < fields; i++ {
				f := t.Field(i)

				key := f.Tag.Get("form")
				if len(key) == 0 {
					continue
				}

				fv := fvalues.Get(key)

				v.Field(i).SetString(fv)
			}
		}
	}

	if reqr, ok := req.(router); ok {
		route := reqr.Route()
		if route != nil {
			vp := reflect.ValueOf(route)
			v := vp.Elem()
			t := v.Type()

			fields := t.NumField()

			for i := 0; i < fields; i++ {
				f := t.Field(i)

				key := f.Tag.Get("route")
				if len(key) == 0 {
					continue
				}

				rv := chi.URLParam(r, key)

				v.Field(i).SetString(rv)
			}
		}
	}

	if reqv, ok := req.(validater); ok {
		return reqv.Validate()
	}

	return nil
}
