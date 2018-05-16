package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rickbassham/goapi/router"
)

type INDIHandlers interface {
	GetClients(w http.ResponseWriter, r *http.Request)
	PostConnect(w http.ResponseWriter, r *http.Request)
	GetDevices(w http.ResponseWriter, r *http.Request)
	PostSetSwitchValue(w http.ResponseWriter, r *http.Request)
	PostSetNumberValue(w http.ResponseWriter, r *http.Request)
	PostSetTextValue(w http.ResponseWriter, r *http.Request)
	PostEnableBlob(w http.ResponseWriter, r *http.Request)
	GetBlob(w http.ResponseWriter, r *http.Request)
	GetBlobStream(w http.ResponseWriter, r *http.Request)

	GetDrivers(w http.ResponseWriter, r *http.Request)
	PostStartServer(w http.ResponseWriter, r *http.Request)
	PostStopServer(w http.ResponseWriter, r *http.Request)
	PostStartDriver(w http.ResponseWriter, r *http.Request)
	PostStopDriver(w http.ResponseWriter, r *http.Request)
}

type indiRoutes struct {
	handlerService INDIHandlers
}

func NewINDIRoutes(handlerService INDIHandlers) router.RouteCreater {
	return &indiRoutes{
		handlerService: handlerService,
	}
}

func (svc *indiRoutes) CreateRoutes(r chi.Router) chi.Router {
	r.Post("/connect", svc.handlerService.PostConnect)
	r.Get("/clients", svc.handlerService.GetClients)

	r.Route("/{clientId}/devices", func(r chi.Router) {
		r.Get("/", svc.handlerService.GetDevices)

		r.Route("/enableblob", func(r chi.Router) {
			r.Post("/set", svc.handlerService.PostEnableBlob)
		})
		r.Route("/switches", func(r chi.Router) {
			r.Post("/set", svc.handlerService.PostSetSwitchValue)
		})
		r.Route("/texts", func(r chi.Router) {
			r.Post("/set", svc.handlerService.PostSetTextValue)
		})
		r.Route("/numbers", func(r chi.Router) {
			r.Post("/set", svc.handlerService.PostSetNumberValue)
		})

		r.Get("/{deviceName}/blobs/{propName}/{blobName}", svc.handlerService.GetBlob)
		r.Get("/{deviceName}/blobs/{propName}/{blobName}/stream", svc.handlerService.GetBlobStream)
	})

	r.Route("/server", func(r chi.Router) {
		r.Post("/start", svc.handlerService.PostStartServer)
		r.Post("/stop", svc.handlerService.PostStopServer)

		r.Route("/drivers", func(r chi.Router) {
			r.Get("/", svc.handlerService.GetDrivers)
			r.Post("/start", svc.handlerService.PostStartDriver)
			r.Post("/stop", svc.handlerService.PostStopDriver)
		})
	})

	return r
}
