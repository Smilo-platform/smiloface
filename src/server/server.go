package server

import (
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/pat"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
)

var (
	router *pat.Router
	log    *logrus.Entry = logrus.WithField("package", "server")
)

var (
	LimitDisabled = os.Getenv("limitdisabled")
)

var DEBUG_STR = os.Getenv("DEBUG")
var DEBUG = false

// SetLogger set the logger
func SetLogger(loggers *logrus.Entry) {
	log = loggers.WithFields(log.Data)
}

func init() {

	if DEBUG_STR == "true" {
		DEBUG = true
	}
}

//NewServer return pointer to new created server object
func NewServer(Port string) *http.Server {
	router = InitRouting()
	return &http.Server{
		Addr:    ":" + Port,
		Handler: router,
	}
}

//StartServer start and listen @server
func StartServer(Port string, loggers *logrus.Entry) {
	//init log
	SetLogger(loggers)

	log.Info("Starting server")
	s := NewServer(Port)
	log.Info("Server starting --> " + Port)

	//enable graceful shutdown
	err := gracehttp.Serve(
		s,
	)

	if err != nil {
		log.Fatalf("Error: %v", err)
		os.Exit(0)
	}

}

func InitRouting() *pat.Router {

	r := pat.New()

	r.PathPrefix("/face").HandlerFunc(FacialRecognition)
	r.PathPrefix("/web").HandlerFunc(Webcam)

	return r
}
