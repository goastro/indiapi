package handler

import (
	"errors"
	"net/http"

	"github.com/rickbassham/goapi/requestparser"
)

type PostSetNumberValueRequest struct {
	Device   string `json:"device"`
	Property string `json:"property"`
	Number   string `json:"number"`
	Value    string `json:"value"`
	ClientID string `route:"clientId"`
}

func (r *PostSetNumberValueRequest) Body() interface{} {
	return r
}

func (r *PostSetNumberValueRequest) Route() interface{} {
	return r
}

func (r *PostSetNumberValueRequest) Validate() error {
	if len(r.ClientID) == 0 {
		return errors.New("missing clientId")
	}

	return nil
}

func (svc *INDIService) PostSetNumberValue(w http.ResponseWriter, r *http.Request) {
	log := svc.log.FromContext(r.Context())
	req := PostSetNumberValueRequest{}
	err := requestparser.ParseRequest(r, &req)
	if err != nil {
		log.WithError(err).Warn("error in bodyJson")
		writeJsonResponse(w, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
		return
	}

	c := svc.getClient(req.ClientID)
	if c == nil {
		log.WithField("clientId", req.ClientID).Warn("unknown client")
		writeJsonResponse(w, http.StatusBadRequest, ErrorResponse{Error: "unknown client"})
		return
	}

	err = c.SetNumberValue(req.Device, req.Property, req.Number, req.Value)
	if err != nil {
		log.WithError(err).Warn("error in c.SetNumberValue")
		writeJsonResponse(w, http.StatusInternalServerError, ErrorResponse{Error: "unable to set number value"})
		return
	}

	err = writeJsonResponse(w, http.StatusOK, nil)
	if err != nil {
		log.WithError(err).Warn("error in writeJsonResponse")
	}
}
