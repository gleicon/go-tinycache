package networking

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gleicon/go-tinycache/internal/backend"
	"github.com/gleicon/go-tinycache/internal/metrics"
	log "github.com/sirupsen/logrus"
)

// Server holds all dependencies for the networking server
// Add more fields as needed for logging, config, etc.
type Server struct {
	Metrics     *metrics.InternalMetrics
	Backend     backend.BackendDatabase
	BackendType string
	Logger      *log.Logger
	Messages    chan string
}

// NewServer constructs a new Server with dependencies
func NewServer(metrics *metrics.InternalMetrics, backendType string, logger *log.Logger, filename string) *Server {
	mia := &Server{
		Metrics:     metrics,
		BackendType: backendType,
		Backend:     nil, // Will be set later
		Logger:      logger,
		Messages:    make(chan string),
	}
	mia.Backend = mia.loadDB(filename)
	return mia
}

// loadDB initializes the backend database based on the type specified
// and returns the database instance.
// It does it explicitly instead of implicitly direct to the
// type attribute so we can use the same method
// to switch dataases on the fly

func (s *Server) loadDB(filename string) backend.BackendDatabase {
	var vdb backend.BackendDatabase
	var err error
	switch s.BackendType {
	case "boltdb":
		vdb, err = backend.NewKVBoltDBBackend(filename, "memcached", 1000000)
	case "inmem":
		vdb, err = backend.NewLRUMemoryBackend(1000000)
	default:
		s.Logger.Errorf("Unknown backend: %s", s.BackendType)
		return nil
	}
	if err != nil {
		s.Logger.Errorf("Error opening db %s", err)
		return nil
	}

	s.Logger.Infof("DB %s opened successfully", filename)
	return vdb
}

func (s *Server) switchDBHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "405 Method not allowed", 405)
		return
	}
	filename := req.FormValue("filename")
	if filename == "" {
		http.Error(w, "500 Internal error", 500)
		return
	}
	s.Messages <- filename
	w.Write([]byte("OK"))
}

func (s *Server) Serve(ip string, port string, filename string) {
	var err error
	go func() {
		http.HandleFunc("/api/v1/switchdb", s.switchDBHandler)
		http.ListenAndServe(":8080", nil)
	}()
	addr := fmt.Sprintf("%s:%s", ip, port)
	listener, err := net.Listen("tcp", addr)
	defer listener.Close()

	ms := NewMemcachedProtocolServer(false, s.Metrics)

	go func() {
		for {
			filename := <-s.Messages
			currentVdb := s.Backend

			if filename != "" {
				if s.Backend.GetDbPath() == filename {
					s.Logger.Errorf("DB Switch from %s to %s - Aborted, db already open", currentVdb.GetDbPath(), filename)
					continue
				}
				ms.ReadOnly(true)

				s.Logger.Infof("DB Switch from %s to %s", currentVdb.GetDbPath(), filename)
				time.Sleep(2 * time.Second)
				vdb := s.loadDB(filename)
				if vdb == nil {
					s.Metrics.NetworkErrors.Inc(1)
					s.Logger.Errorf("DB Switch from %s to %s - Failed, db not loaded", vdb.GetDbPath(), filename)
					ms.ReadOnly(false)
					continue
				} else {
					s.Backend = vdb
				}
				time.Sleep(2 * time.Second)
				currentVdb.Close()
				s.Logger.Infof("DB Switch from %s to %s done", vdb.GetDbPath(), filename)
				ms.ReadOnly(false)
			}
		}
	}()

	if err == nil {
		for {
			if conn, err := listener.Accept(); err == nil {
				s.Metrics.TotalConnections.Inc(1)
				go ms.Parse(conn, s.Backend)
			} else {
				s.Metrics.NetworkErrors.Inc(1)
				s.Logger.Error(err.Error())
			}
		}
	} else {
		s.Metrics.NetworkErrors.Inc(1)
		s.Logger.Fatal(err.Error())
	}
}
