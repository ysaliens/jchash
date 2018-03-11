package handlers

import "net/http"
import "io"
import "os"
import "syscall"

func MakeHandleShutdown(Stop chan os.Signal) func(http.ResponseWriter, *http.Request){
	return func(w http.ResponseWriter, r *http.Request){
		Stop <- syscall.SIGTERM
		//Send shutdown notice to requestor (not specified by requirement but nice to have)
		w.WriteHeader(http.StatusOK)
		f, err := io.WriteString(w, "Server shutdown initiated")
		if checkError(f, err) == -1 {
			return
		}
	}
}