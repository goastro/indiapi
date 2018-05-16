package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/goastro/indiclient"
	"github.com/goastro/indiserver"
	"github.com/rickbassham/goapi/middleware"
	"github.com/spf13/afero"
)

// INDIClient defines the functions available for an INDI Client.
type INDIClient interface {
	// Connect dials to create a connection to address over the specified network.
	Connect(network, address string) error
	// Disconnect clears out all devices from memory, closes the connection, and closes the read and write channels.
	Disconnect() error
	// IsConnected returns true if the client is currently connected to an INDI server. Otherwise, returns false.
	IsConnected() bool

	// Devices returns the current list of INDI devices with their current state.
	Devices() []indiclient.Device
	// GetBlob finds a BLOB with the given deviceName, propName, blobName. Be sure to close rdr when you are done with it.
	GetBlob(deviceName, propName, blobName string) (rdr io.ReadCloser, fileName string, length int64, err error)
	// GetBlobStream finds a BLOB with the given deviceName, propName, blobName. This will return an io.Pipe that can stream the BLOBs that are received from the indiserver.
	GetBlobStream(deviceName, propName, blobName string) (rdr io.ReadCloser, id string, err error)
	// CloseBlobStream closes the blob stream created by GetBlobStream.
	CloseBlobStream(deviceName, propName, blobName string, id string) (err error)

	// GetProperties sends a command to the INDI server to retreive the property definitions for the given deviceName and propName.
	// deviceName and propName are optional.
	GetProperties(deviceName, propName string) error
	// EnableBlob sends a command to the INDI server to enable/disable BLOBs for the current connection.
	// It is recommended to enable blobs on their own client, and keep the main connection clear of large transfers.
	// By default, BLOBs are NOT enabled.
	EnableBlob(deviceName, propName string, val indiclient.BlobEnable) error
	// SetTextValue sends a command to the INDI server to change the value of a textVector.
	SetTextValue(deviceName, propName, textName, textValue string) error
	// SetNumberValue sends a command to the INDI server to change the value of a numberVector.
	SetNumberValue(deviceName, propName, numberName, numberValue string) error
	// SetSwitchValue sends a command to the INDI server to change the value of a switchVector.
	SetSwitchValue(deviceName, propName, switchName string, switchValue indiclient.SwitchState) error
	// SetBlobValue sends a command to the INDI server to change the value of a blobVector.
	SetBlobValue(deviceName, propName, blobName, blobValue, blobFormat string, blobSize int) error
}

type INDIServer interface {
	Drivers() map[string][]indiserver.Driver
	StartServer() error
	StopServer() error
	StartDriver(driver, name string) error
	StopDriver(driver, name string) error
}

type INDIService struct {
	log        middleware.RequestLogger
	svr        INDIServer
	dialer     indiclient.Dialer
	fs         afero.Fs
	bufferSize int

	clients sync.Map
}

func NewINDIService(log middleware.RequestLogger, svr INDIServer, dialer indiclient.Dialer, fs afero.Fs, bufferSize int) *INDIService {
	return &INDIService{
		log:        log,
		svr:        svr,
		dialer:     dialer,
		fs:         fs,
		bufferSize: bufferSize,
	}
}

func writeJsonResponse(w http.ResponseWriter, statusCode int, body interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	enc := json.NewEncoder(w)
	err = enc.Encode(body)
	return
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (svc *INDIService) getClient(clientID string) INDIClient {
	if c, ok := svc.clients.Load(clientID); ok {
		return c.(INDIClient)
	}
	return nil
}

func (svc *INDIService) addClient(clientID string, c INDIClient) {
	svc.clients.Store(clientID, c)
}
