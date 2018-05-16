package handler

import (
	"errors"
	"net/http"

	"github.com/rickbassham/goapi/requestparser"
)

type PostSetTextValueRequest struct {
	Device   string `json:"device"`
	Property string `json:"property"`
	Text     string `json:"text"`
	Value    string `json:"value"`
	ClientID string `route:"clientId"`
}

func (r *PostSetTextValueRequest) Body() interface{} {
	return r
}

func (r *PostSetTextValueRequest) Route() interface{} {
	return r
}

func (r *PostSetTextValueRequest) Validate() error {
	if len(r.ClientID) == 0 {
		return errors.New("missing clientId")
	}

	return nil
}

func (svc *INDIService) PostSetTextValue(w http.ResponseWriter, r *http.Request) {
	log := svc.log.FromContext(r.Context())
	req := PostSetTextValueRequest{}
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

	err = c.SetTextValue(req.Device, req.Property, req.Text, req.Value)
	if err != nil {
		log.WithError(err).Warn("error in c.SetTextValue")
		writeJsonResponse(w, http.StatusInternalServerError, ErrorResponse{Error: "unable to set text value"})
		return
	}

	err = writeJsonResponse(w, http.StatusOK, nil)
	if err != nil {
		log.WithError(err).Warn("error in writeJsonResponse")
	}
}
