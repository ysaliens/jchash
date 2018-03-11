package main

import "flag"
import "strconv"
import "log"
import "context"
import "time"
import "os"
import "os/signal"
import "syscall"
import "github.com/ysaliens/jchash/handlers"


//Start server
func main() {
	//Get port from command line (or default to 8080)
	port := flag.Int("port", 8080, "Port to listen on")
	flag.Parse()
	address := ":"+strconv.Itoa(*port)

	//Initialize new server, db, stats
	server := handlers.Create(address)

	//Set interrupt on Control+C or SIGINT
	signal.Notify(server.Stop, os.Interrupt, syscall.SIGTERM)

	//Start server as subroutine
	go func(){
		log.Printf("Listening on http://localhost%v",address)
		if err := server.Server.ListenAndServe(); err != nil {
			log.Printf("%s", err)
		}
	}()

	//Main execution pauses until we receive a request to shutdown on Stop channel
	<-server.Stop

	//Gracefully bring down server (Allow up to 5 sec to finish processing)
	log.Print("Received shutdown request, shutting down.")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Server.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down: %s", err)
	} else {
		log.Println("Server gracefully shut down.")
	}
}