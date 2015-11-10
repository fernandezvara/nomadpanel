package api

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/fernandezvara/nomadpanel/usage"
	"github.com/gorilla/mux"
)

type handler func(c *Context, w http.ResponseWriter, r *http.Request)

type methods map[string]map[string]handler

var m = map[string]map[string]handler{
	"GET": {
		"/api/version":             getVersion,
		"/api/usage/{dc:.*}/nodes": getUsageByDCNodes,
		"/api/usage/{dc:.*}":       getUsageByDC,
	},
	"POST":   {},
	"PUT":    {},
	"DELETE": {},
	"OPTIONS": {
		"": optionsHandler,
	},
}

// Context holds the required information needed by the API
type Context struct {
	version    string
	listenAddr string
	usage      *usage.Usage
	log        *logrus.Logger
	methods    methods
}

// NewContext returns the structure needed to manage the API
func NewContext(listenAddr, version string, usage *usage.Usage, log *logrus.Logger) *Context {
	return &Context{
		version:    version,
		listenAddr: listenAddr,
		usage:      usage,
		log:        log,
		methods:    m,
	}
}

// ListenAndServe runs the API webserver loop
func ListenAndServe(c *Context) error {
	router := createRouter(c)
	c.log.WithFields(logrus.Fields{"service": "main"}).Infoln("API: Started.")
	c.log.WithFields(logrus.Fields{"service": "main"}).Infoln("Listening: " + c.listenAddr)
	return http.ListenAndServe(c.listenAddr, router)
}

func createRouter(c *Context) *mux.Router {
	// borrow from swarm
	r := mux.NewRouter()
	for method, routes := range c.methods {
		for route, fct := range routes {
			c.log.WithFields(logrus.Fields{"method": method, "route": route}).Debug("Registering HTTP route")
			// NOTE: scope issue, make sure the variables are local and won't be changed
			localRoute := route
			localFct := fct
			localMethod := method
			wrap := func(w http.ResponseWriter, r *http.Request) {
				c.log.WithFields(logrus.Fields{"method": r.Method, "uri": r.RequestURI, "ip": r.RemoteAddr}).Info("HTTP request received")
				localFct(c, w, r)
			}
			// add the new route
			if localRoute != "" {
				r.Path(localRoute).Methods(localMethod).HandlerFunc(wrap)
			} else {
				r.Methods(localMethod).HandlerFunc(wrap)
			}
		}
	}
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	return r
}

func optionsHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
