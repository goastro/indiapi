package handler

import (
	"net/http"

	"github.com/goastro/indiclient"
)

type Client struct {
	ClientID    string `json:"clientId"`
	IsConnected bool   `json:"connected"`
}
type GetClientsResponse struct {
	Clients []Client `json:"clients"`
}

func (svc *INDIService) GetClients(w http.ResponseWriter, r *http.Request) {
	log := svc.log.FromContext(r.Context())
	var clients GetClientsResponse

	clients.Clients = []Client{}

	svc.clients.Range(func(key, value interface{}) bool {
		clients.Clients = append(clients.Clients, Client{
			ClientID:    key.(string),
			IsConnected: value.(*indiclient.INDIClient).IsConnected(),
		})

		return true
	})

	err := writeJsonResponse(w, http.StatusOK, clients)
	if err != nil {
		log.WithError(err).Warn("error in writeJsonResponse")
	}
}
