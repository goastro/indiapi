package handler

import (
	"errors"
	"net/http"

	"github.com/goastro/indiclient"
	"github.com/google/uuid"
	"github.com/rickbassham/goapi/requestparser"
)

type PostConnectRequestBody struct {
	Network  string `json:"network"`
	Address  string `json:"address"`
	Device   string `json:"device,omitempty"`
	Property string `json:"property,omitempty"`
}

func (r *PostConnectRequestBody) Body() interface{} {
	return r
}

func (r *PostConnectRequestBody) Validate() error {
	if len(r.Network) == 0 {
		return errors.New("network not specified")
	}

	if len(r.Address) == 0 {
		return errors.New("address not specified")
	}

	return nil
}

type PostConnectResponseBody struct {
	ClientID string `json:"clientId"`
}

func (svc *INDIService) PostConnect(w http.ResponseWriter, r *http.Request) {
	log := svc.log.FromContext(r.Context())
	req := PostConnectRequestBody{}

	err := requestparser.ParseRequest(r, &req)
	if err != nil {
		log.WithError(err).Warn("error in requestparser.ParseRequest")
		writeJsonResponse(w, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
		return
	}

	clientID := uuid.New()
	c := indiclient.NewINDIClient(log, svc.dialer, svc.fs, svc.bufferSize)
	svc.addClient(clientID.String(), c)

	err = c.Connect(req.Network, req.Address)
	if err != nil {
		log.WithError(err).Warn("error in c.Connect")
		writeJsonResponse(w, http.StatusInternalServerError, ErrorResponse{Error: "unable to connect"})
		return
	}

	err = c.GetProperties(req.Device, req.Property)
	if err != nil {
		log.WithError(err).Warn("error in c.GetProperties")
		writeJsonResponse(w, http.StatusInternalServerError, ErrorResponse{Error: "unable to connect"})
		return
	}

	err = writeJsonResponse(w, http.StatusOK, &PostConnectResponseBody{
		ClientID: clientID.String(),
	})
	if err != nil {
		log.WithError(err).Warn("error in writeJsonResponse")
	}
}
