package handler

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/goastro/fitsutils"
	"github.com/rickbassham/goapi/requestparser"
)

type GetBlobRequest struct {
	Device   string `route:"deviceName"`
	Property string `route:"propName"`
	Blob     string `route:"blobName"`
	ClientID string `route:"clientId"`
}

func (r *GetBlobRequest) Route() interface{} {
	return r
}

func (r *GetBlobRequest) Validate() error {
	if len(r.ClientID) == 0 {
		return errors.New("missing clientId")
	}

	return nil
}

func (svc *INDIService) GetBlob(w http.ResponseWriter, r *http.Request) {
	log := svc.log.FromContext(r.Context())

	var req GetBlobRequest
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

	rdr, fileName, size, err := c.GetBlob(req.Device, req.Property, req.Blob)
	if err != nil {
		log.WithError(err).Warn("error in c.GetBlob")
		writeJsonResponse(w, http.StatusInternalServerError, ErrorResponse{Error: "unable to get blob"})
		return
	}

	defer rdr.Close()

	if r.Header.Get("Accept-Encoding") == "image/png" {
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)

		err = fitsutils.SavePNG(rdr, w)
		if err != nil {
			log.WithError(err).Warn("error in fitsutils.SavePNG")
		}

		return
	}

	w.Header().Set("Content-Length", strconv.FormatInt(size, 10))

	if strings.HasSuffix(fileName, ".fits") {
		w.Header().Set("Content-Type", "image/fits")
	} else if strings.HasSuffix(fileName, ".fits.z") {
		w.Header().Set("Content-Type", "image/fits")
		w.Header().Set("Content-Encoding", "deflate")
	} else {
		w.Header().Set("Content-Type", "application/x-binary")
	}

	w.WriteHeader(200)

	io.Copy(w, rdr)
}

func (svc *INDIService) GetBlobStream(w http.ResponseWriter, r *http.Request) {
	log := svc.log.FromContext(r.Context())

	var req GetBlobRequest
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

	rdr, id, err := c.GetBlobStream(req.Device, req.Property, req.Blob)
	if err != nil {
		log.WithError(err).Warn("error in c.GetBlobStream")
		writeJsonResponse(w, http.StatusInternalServerError, ErrorResponse{Error: "unable to get blob stream"})
		return
	}

	defer c.CloseBlobStream(req.Device, req.Property, req.Blob, id)
	defer rdr.Close()

	w.Header().Set("Content-Type", "application/x-binary")
	w.WriteHeader(200)

	io.Copy(w, rdr)
}
