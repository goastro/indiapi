package handler

import (
	"errors"
	"net/http"

	"github.com/rickbassham/goapi/requestparser"
)

type PostDisconnectRequestBody struct {
	ClientID string `route:"clientId"`
}

func (r *PostDisconnectRequestBody) Route() interface{} {
	return r
}

func (r *PostDisconnectRequestBody) Validate() error {
	if len(r.ClientID) == 0 {
		return errors.New("clientId not specified")
	}

	return nil
}

func (svc *INDIService) PostDisconnect(w http.ResponseWriter, r *http.Request) {
	log := svc.log.FromContext(r.Context())
	req := PostDisconnectRequestBody{}

	err := requestparser.ParseRequest(r, &req)
	if err != nil {
		log.WithError(err).Warn("error in requestparser.ParseRequest")
		writeJsonResponse(w, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
		return
	}

	c := svc.getClient(req.ClientID)
	if c == nil {
		writeJsonResponse(w, http.StatusBadRequest, ErrorResponse{Error: "unknown clientID"})
		return
	}

	err = c.Disconnect()
	if err != nil {
		log.WithError(err).Warn("error in c.Disconnect")
		writeJsonResponse(w, http.StatusBadRequest, ErrorResponse{Error: "unable to disconnect"})
		return
	}

	svc.removeClient(req.ClientID)

	err = writeJsonResponse(w, http.StatusOK, nil)
	if err != nil {
		log.WithError(err).Warn("error in writeJsonResponse")
	}
}
