package handlers

import "net/http"
import "encoding/json"
import "strings"
import "time"
import "log"

//Force all requests to go through here, update server stats as needed
//NOTE: Logging and prints slows down performance
//TODO: Add additional logging features
func MakeStats(handler http.Handler, stats *Stats) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    	//Save time everytime we get a request
    	start := time.Now()

        //Log connection
        log.Printf("%s %s", r.RemoteAddr, r.URL)

        //Print all info for request
        //log.Printf("Request: %v", r)

    	//Route and fulfil request
        handler.ServeHTTP(w, r)

        //Only count hash requests in stats (even if request is broken) from client side
        //This does not count the background hash-write operations as user is not aware
        //of these as per requirement. If requirement changes, this will need to be changed
        if  strings.Contains(r.URL.Path, "hash"){
            //Calculate request time
            end := time.Now()
            requestTime := end.Sub(start)

            //Update server stats - thread-safe
            stats.timeLock.Lock()
            stats.totalTime = stats.totalTime + requestTime
            stats.AverageTime = float32(stats.totalTime / time.Millisecond) / float32(stats.Requests)
            stats.Requests++
            stats.timeLock.Unlock()
            log.Printf("Request Complete: %v Average: %v",requestTime,stats.AverageTime)
        } 
    })
}

//Returns sender json object of stats
//Above function will drop down here when /stats is sent
func MakeHandleStats(stats *Stats) func(w http.ResponseWriter, r *http.Request){
    return func (w http.ResponseWriter, r *http.Request){
        //Note totalTime does not picked up by JSON as it's lowercase
        jsonStats, err := json.Marshal(stats)
        if err != nil {
            sendError(w, "404 Error getting JSON object")
            return
        }
        w.Header().Set("Content-Type", "application/json")
        f, err := w.Write(jsonStats)
        checkError(f, err)
    }
}
