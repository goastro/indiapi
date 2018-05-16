package handler

import (
	"net/http"

	"github.com/rickbassham/goapi/requestparser"
)

func (svc *INDIService) GetDrivers(w http.ResponseWriter, r *http.Request) {
	log := svc.log.FromContext(r.Context())

	drivers := svc.svr.Drivers()

	err := writeJsonResponse(w, http.StatusOK, drivers)
	if err != nil {
		log.WithError(err).Warn("error in writeJsonResponse")
	}
}

func (svc *INDIService) PostStartServer(w http.ResponseWriter, r *http.Request) {
	log := svc.log.FromContext(r.Context())

	err := svc.svr.StartServer()
	if err != nil {
		log.WithError(err).Warn("error in svc.svr.StartServer")
		writeJsonResponse(w, http.StatusInternalServerError, ErrorResponse{Error: "unable to start indiserver"})
		return
	}

	err = writeJsonResponse(w, http.StatusOK, nil)
	if err != nil {
		log.WithError(err).Warn("error in writeJsonResponse")
	}
}

func (svc *INDIService) PostStopServer(w http.ResponseWriter, r *http.Request) {
	log := svc.log.FromContext(r.Context())

	err := svc.svr.StopServer()
	if err != nil {
		log.WithError(err).Warn("error in svc.svr.StopServer")
		writeJsonResponse(w, http.StatusInternalServerError, ErrorResponse{Error: "unable to stop indiserver"})
		return
	}

	err = writeJsonResponse(w, http.StatusOK, nil)
	if err != nil {
		log.WithError(err).Warn("error in writeJsonResponse")
	}
}

type PostDriverRequest struct {
	Driver string `json:"driver"`
	Name   string `json:"name"`
}

func (r *PostDriverRequest) Body() interface{} {
	return r
}

func (svc *INDIService) PostStartDriver(w http.ResponseWriter, r *http.Request) {
	log := svc.log.FromContext(r.Context())

	var req PostDriverRequest
	err := requestparser.ParseRequest(r, &req)
	if err != nil {
		log.WithError(err).Warn("error in requestparser.ParseRequest")
		writeJsonResponse(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request"})
		return
	}

	err = svc.svr.StartDriver(req.Driver, req.Name)
	if err != nil {
		log.WithError(err).Warn("error in svc.svr.StartDriver")
		writeJsonResponse(w, http.StatusInternalServerError, ErrorResponse{Error: "unable to start driver"})
		return
	}

	err = writeJsonResponse(w, http.StatusOK, nil)
	if err != nil {
		log.WithError(err).Warn("error in writeJsonResponse")
	}
}

func (svc *INDIService) PostStopDriver(w http.ResponseWriter, r *http.Request) {
	log := svc.log.FromContext(r.Context())

	var req PostDriverRequest
	err := requestparser.ParseRequest(r, &req)
	if err != nil {
		log.WithError(err).Warn("error in requestparser.ParseRequest")
		writeJsonResponse(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request"})
		return
	}

	err = svc.svr.StopDriver(req.Driver, req.Name)
	if err != nil {
		log.WithError(err).Warn("error in svc.svr.StopDriver")
		writeJsonResponse(w, http.StatusInternalServerError, ErrorResponse{Error: "unable to stop driver"})
		return
	}

	err = writeJsonResponse(w, http.StatusOK, nil)
	if err != nil {
		log.WithError(err).Warn("error in writeJsonResponse")
	}
}
