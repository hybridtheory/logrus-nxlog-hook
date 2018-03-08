package nxlog

import (
	"crypto/tls"
	"net"
	"reflect"
)

// Writer represents a writer to NXLog
type Writer struct {
	protocol   string
	endpoint   string
	settings   interface{}
	connection net.Conn
}

// NewWriter returns a reference to a new writer, with protocol, endpoint and
// settings populated. It also tries to connect to it and returns nil on
// success or error otherwise.
func NewWriter(protocol string, endpoint string, settings interface{}) (*Writer, error) {
	w := &Writer{
		protocol: protocol,
		endpoint: endpoint,
		settings: settings,
	}
	return w, w.Connect()
}

// Connect with the protocol and endpoint specified in the writer.
func (w *Writer) Connect() error {
	if w.protocol == "ssl" {
		return w.getTLSWriter()
	}
	return w.getNetWriter()
}

// Reconnect will close the connection and open it again.
func (w *Writer) Reconnect() {
	w.Close()
	w.Connect()
}

// Close will close the current connection if possible.
func (w *Writer) Close() {
	w.connection.Close()
}

// Write will send a message to the open connection of the writer.
// In case of error, it will try to reconnect once.
func (w *Writer) Write(message []byte, reconnect bool) (int, error) {
	n, err := w.connection.Write(message)
	if err != nil && reconnect {
		w.Reconnect()
		return w.Write(message, false)
	}
	return n, err
}

func (w *Writer) getNetWriter() error {
	connection, err := net.Dial(w.protocol, w.endpoint)
	w.connection = connection
	return err
}

func (w *Writer) getTLSWriter() error {
	configuration := w.settings
	if configuration == nil || reflect.TypeOf(configuration).String() != "*tls.Config" {
		configuration = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	connection, err := tls.Dial("tcp", w.endpoint, configuration.(*tls.Config))
	w.connection = connection
	return err
}
