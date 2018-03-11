package handlers

import "net/http"
import "log"
import "time"
import "io"
import "strconv"
import "path"
import "github.com/ysaliens/jchash/encode"

//Hash handler - drops to POST/PUT/Error functions depending on method type
func MakeHandleHash(data *Db) func(http.ResponseWriter, *http.Request){
	return func(w http.ResponseWriter, r *http.Request){
		if r.Method == http.MethodPost {
			hashPut(data, w, r)
		} else if r.Method == http.MethodGet {
			hashGet(data, w, r)
		} else {
			sendError(w, "404 Page Not Found: Unknown Method")
		}
	}
}

//Log errors
//TODO: Implement panic and recover for functions
func checkError(result int, e error) int {
	if e != nil {
		log.Printf("%v", result)
		log.Printf("%v", e)
		return -1
	}
	return 0
}

//Send error to client and log it
func sendError(w http.ResponseWriter, err string){
	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w,err)
	log.Printf(err)
}

//Hash a password & save to database
func hashPut(data * Db, w http.ResponseWriter, r *http.Request){
	//Extract password from request
	r.ParseForm()
	password := r.PostFormValue("password")
	
	//Create new password entry in map
	if len(password) > 0 {

		//Get new entry ID using mutex for thread-safe
		data.lock.Lock()
		data.idCount++
		newID := data.idCount
		data.lock.Unlock()

		//Create hash entry & store in sync map
		entry := Entry{}
		entry.ID = newID
		entry.Done = 0
		data.DB.Store(entry.ID, &entry)

		//Send ID to client
		//log.Printf("%v",entry.ID)
		f, err := io.WriteString(w, strconv.Itoa(entry.ID))
		if checkError(f, err) == -1 {
			return
		}

		//Encode the password (and save to sync map) as a goroutine
		go func() {
			entry.EncodeValue = encode.Encode(password)
			//log.Printf("Hash: %v",entry.EncodeValue)

			//Sleep 5s as per requirement
			time.Sleep(5 * time.Second)
			
			//Only mark it done after waiting the 5 sec to simulate contention
			entry.Done = 1
		}()

	} else {
		sendError(w, "404 Page Not Found")
		return
	}
}

//Retrieve a hash from database given ID
func hashGet(data * Db, w http.ResponseWriter, r *http.Request){
	//Get ID from URL Path (id can ONLY be a positive number)
	idRequested, err := strconv.Atoi(path.Base(r.URL.Path))
	if idRequested < 1 || err != nil {
		sendError(w, "404 Page Not Found")
		return
	}

	//Check database for ID
	entry, found := data.DB.Load(idRequested)
	if found == false || entry == nil {
		sendError(w, "404 Page Not Found")
		return
	}

	//Check if entry has completed writing
	//If not, wait and keep checking every 1s
	//Timeout after 10s
	timeout := 10
	hash := entry.(*Entry)
	for hash.Done != 1 {
		time.Sleep(1 * time.Second)
		timeout--
		if (timeout == 0){
			sendError(w, "404 Not Found: Database write timeout")
			return
		}
	}

	//Send hash to client
	f, err:= io.WriteString(w, hash.EncodeValue)
	checkError(f, err)
}