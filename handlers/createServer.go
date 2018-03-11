package handlers

import "net/http"
import "os"
import "sync"
import "time"

//Server info
type HTTPServer struct {
	Address     string
	Server      *http.Server
	Handler     *http.Handler
    Stop        chan os.Signal
}

//Sync map for storage, thread-safe
//TODO: Separate db into new package with real DB
type Db struct{
	DB          *sync.Map
	lock        *sync.Mutex
	idCount     int
}

//Entry contains info for a single password hash
type Entry struct {
	ID 	   		int
	EncodeValue string
	Done		int
}

//Server stats info
//Only JSON entries will be encoded to JSON when sending stats
//Using a struct also prevents JSON marshal from re-ordering
type Stats struct {
	Requests    int     `json:"total"`
	AverageTime float32 `json:"average"`
	totalTime   time.Duration
	timeLock    *sync.Mutex
}

//Initialize server, db, stats
func Create(address string) *HTTPServer {

	//Initialize stats
	stats := &Stats{
		Requests:    0,
		AverageTime: 0,
		totalTime:   0,
		timeLock:    &sync.Mutex{},
	}
	handler := MakeStats(http.DefaultServeMux, stats)

	//Initialize database
	data := &Db{
		DB:          new(sync.Map),
		idCount:     0,
		lock:        &sync.Mutex{},
	}

	//Initialize server and stop channel
	server := &http.Server{Addr: address, Handler: handler}
	HTTPServer := &HTTPServer{
		Address:     address,
		Server:      server,
		Handler:     &handler,
		Stop:        make(chan os.Signal, 2),
	}

	//Register server functions
	HandleHash     := MakeHandleHash(data)
	HandleShutdown := MakeHandleShutdown(HTTPServer.Stop)
	HandleStats    := MakeHandleStats(stats)
	
	http.HandleFunc("/hash",     HandleHash)
	http.HandleFunc("/hash/",    HandleHash)
	http.HandleFunc("/shutdown", HandleShutdown)
	http.HandleFunc("/stats",    HandleStats)

	return HTTPServer
}