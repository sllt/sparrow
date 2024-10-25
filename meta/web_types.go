package meta

import (
	"net/http"
	"time"

	"github.com/sllt/sparrow/gen"
)

type WebServerOptions struct {
	Host        string
	Port        uint16
	CertManager gen.CertManager
	Handler     http.Handler
}
type WebHandlerOptions struct {
	Worker         gen.Atom
	RequestTimeout time.Duration
}

type MessageWebRequest struct {
	Response http.ResponseWriter
	Request  *http.Request
	Done     func()
}
