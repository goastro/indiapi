package handler

import (
	"errors"
	"net/http"

	"github.com/goastro/indiclient"
	"github.com/rickbassham/goapi/requestparser"
)

type PostSetSwitchValueRequest struct {
	Device   string                 `json:"device"`
	Property string                 `json:"property"`
	Switch   string                 `json:"switch"`
	Value    indiclient.SwitchState `json:"value"`
	ClientID string                 `route:"clientId"`
}

func (r *PostSetSwitchValueRequest) Body() interface{} {
	return r
}

func (r *PostSetSwitchValueRequest) Route() interface{} {
	return r
}

func (r *PostSetSwitchValueRequest) Validate() error {
	if len(r.ClientID) == 0 {
		return errors.New("missing clientId")
	}

	if len(r.Device) == 0 {
		return errors.New("missing device name")
	}

	if len(r.Property) == 0 {
		return errors.New("missing property name")
	}

	if len(r.Switch) == 0 {
		return errors.New("missing switch name")
	}

	if len(r.Value) == 0 {
		return errors.New("missing switch value")
	}

	return nil
}

func (svc *INDIService) PostSetSwitchValue(w http.ResponseWriter, r *http.Request) {
	log := svc.log.FromContext(r.Context())
	req := PostSetSwitchValueRequest{}
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

	err = c.SetSwitchValue(req.Device, req.Property, req.Switch, req.Value)
	if err != nil {
		log.WithError(err).Warn("error in c.SetSwitchValue")
		writeJsonResponse(w, http.StatusInternalServerError, ErrorResponse{Error: "unable to set switch value"})
		return
	}

	err = writeJsonResponse(w, http.StatusOK, nil)
	if err != nil {
		log.WithError(err).Warn("error in writeJsonResponse")
	}
}
