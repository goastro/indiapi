package handler

import (
	"errors"
	"net/http"

	"github.com/rickbassham/goapi/requestparser"
)

type GetDevicesRequest struct {
	ClientID string `route:"clientId"`
}

func (r *GetDevicesRequest) Route() interface{} {
	return r
}

func (r *GetDevicesRequest) Validate() error {
	if len(r.ClientID) == 0 {
		return errors.New("missing clientId")
	}

	return nil
}

func (svc *INDIService) GetDevices(w http.ResponseWriter, r *http.Request) {
	log := svc.log.FromContext(r.Context())

	var req GetDevicesRequest
	err := requestparser.ParseRequest(r, &req)
	if err != nil {
		log.WithError(err).Warn("error in requestparser.ParseRequest")
		writeJsonResponse(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request"})
		return
	}

	c := svc.getClient(req.ClientID)
	if c == nil {
		log.WithField("clientId", req.ClientID).Warn("unknown client")
		writeJsonResponse(w, http.StatusBadRequest, ErrorResponse{Error: "unknown client"})
		return
	}

	devices := c.Devices()

	err = writeJsonResponse(w, 200, devices)
	if err != nil {
		log.WithError(err).Warn("error in writeJsonResponse")
	}
}
