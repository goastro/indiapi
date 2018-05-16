package handler

import (
	"errors"
	"net/http"

	"github.com/goastro/indiclient"
	"github.com/rickbassham/goapi/requestparser"
)

type EnableBlobRequest struct {
	Device   string                `json:"device"`
	Property string                `json:"property"`
	Enable   indiclient.BlobEnable `json:"enable"`
	ClientID string                `route:"clientId"`
}

func (r *EnableBlobRequest) Body() interface{} {
	return r
}

func (r *EnableBlobRequest) Route() interface{} {
	return r
}

func (r *EnableBlobRequest) Validate() error {
	if len(r.ClientID) == 0 {
		return errors.New("missing clientId")
	}

	return nil
}

func (svc *INDIService) PostEnableBlob(w http.ResponseWriter, r *http.Request) {
	log := svc.log.FromContext(r.Context())
	req := EnableBlobRequest{}
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

	err = c.EnableBlob(req.Device, req.Property, req.Enable)
	if err != nil {
		log.WithError(err).Warn("error in c.EnableBlob")
		writeJsonResponse(w, http.StatusInternalServerError, ErrorResponse{Error: "unable to set enable blob"})
		return
	}

	err = writeJsonResponse(w, http.StatusOK, nil)
	if err != nil {
		log.WithError(err).Warn("error in writeJsonResponse")
	}
}
