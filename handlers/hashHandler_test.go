package handlers

import "testing"
import "time"
import "net/http"
import "net/http/httptest"
import "net/url"
import "strings"

//Create an HTTP request and execute. Returns HTTP response.
func sendCmd(testServer *HTTPServer, pass string, cmdType string, cmd string) *httptest.ResponseRecorder {
	//If password isn't empty, encode it to HTTP form
	form := url.Values{}
	if pass != "" {
		form.Add("password", pass)
	}
	encodedForm := strings.NewReader(form.Encode())

	//Populate HTTP request
	r, _ := http.NewRequest(cmdType, cmd, encodedForm)

	//Add header to request - omitting param=value in header will skip encoding of password data
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	//Send request
	w := httptest.NewRecorder()
	handler := *testServer.Handler
	handler.ServeHTTP(w, r)
	return w
}

//Test ALL handlers (hashing POST, hashing GET, shutdown, stats, errors)
func TestAll(t *testing.T){
	//Make a new server
	testServer := Create(":8080")

	//Allow time to start server
	time.Sleep(2 * time.Second)

	//Invalid command type
	a := sendCmd(testServer, "", "POST", "/WHAT")
	if a.Code != http.StatusNotFound {
		t.Errorf("Response Code: %v", a.Code)
		t.Errorf("Correct: %v", http.StatusNotFound)
	}

	//Invalid method type
	a = sendCmd(testServer, "", "WHAT", "/hash/1")
	if a.Code != http.StatusNotFound {
		t.Errorf("Response Code: %v", a.Code)
		t.Errorf("Correct: %v", http.StatusNotFound)
	}

	//Get a missing hash
	a = sendCmd(testServer, "", "GET", "/hash/1")
	if a.Code != http.StatusNotFound {
		t.Errorf("Response Code: %v", a.Code)
		t.Errorf("Correct: %v", http.StatusNotFound)
	}

	//Negative ID
	a = sendCmd(testServer, "", "GET", "/hash/-5")
	if a.Code != http.StatusNotFound {
		t.Errorf("Response Code: %v", a.Code)
		t.Errorf("Correct: %v", http.StatusNotFound)
	}

	//Zero ID
	a = sendCmd(testServer, "", "GET", "/hash/0")
	if a.Code != http.StatusNotFound {
		t.Errorf("Response Code: %v", a.Code)
		t.Errorf("Correct: %v", http.StatusNotFound)
	}

	//Weird string
	a = sendCmd(testServer, "", "GET", "/hash/hash/hash/3")
	if a.Code != http.StatusNotFound {
		t.Errorf("Response Code: %v", a.Code)
		t.Errorf("Correct: %v", http.StatusNotFound)
	}

	//Try to hash an empty password
	a = sendCmd(testServer, "", "POST", "/hash/")
	if a.Code != http.StatusNotFound {
		t.Errorf("Response Code: %v", a.Code)
		t.Errorf("Correct: %v", http.StatusNotFound)
	}

	//Hash password 1 ("angryMonkey")
	pass := "angryMonkey"
	a = sendCmd(testServer, pass, "POST", "/hash")
	hashID1 := a.Body.String()
	if a.Code != http.StatusOK && hashID1 != "1" {
		t.Errorf("Response Code: %v, Response ID: %v", a.Code, hashID1)
		t.Errorf("Correct: %v, Response ID: 1", http.StatusOK)
	}

	//Get the hash of 'angryMonkey' 
	//(ID should be 1, checked above)
	a = sendCmd(testServer, "", "GET", "/hash/1")
	correct := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
	if a.Body.String() != correct {
		t.Errorf("Server Response: %v", a.Body.String())
		t.Errorf("Correct: %v", correct)
	}

	//Hash password 2 ("test2")
	pass = "test2"
	a = sendCmd(testServer, pass, "POST", "/hash")
	hashID2 := a.Body.String()
	if a.Code != http.StatusOK  && hashID2 != "2" {
		t.Errorf("Response Code: %v, Response ID: %v", a.Code, hashID2)
		t.Errorf("Correct: %v, Response ID: 2", http.StatusOK)
	}

	//Get the hash of "test2" (ID=2, checked above)
	a = sendCmd(testServer, "", "GET", "/hash/2")
	correct = "bSAb7u+1ibCO8GctrII1PQy9mtmeFkLIOhYB89ZHvMoAMle16PMb3B1z++yE+whcedbiZ3t/+SfoI6VOeJFA2Q=="
	if a.Body.String() != correct {
		t.Errorf("Response: %v", a.Body.String())
		t.Errorf("Correct: %v", correct)
	}

	//Stats
	a = sendCmd(testServer, "", "GET", "/stats")
	//t.Errorf("%v",a.Body.String())	//Debug to see the JSON object
	if a.Code != http.StatusOK {
		t.Errorf("Response Code: %v, Response: %v", a.Code, a.Body.String())
		t.Errorf("Correct Code: %v, Response: Server shutdown initiated", http.StatusOK)
	}

	//Shutdown
	a = sendCmd(testServer, "", "GET", "/shutdown")
	if a.Code != http.StatusOK && a.Body.String() != "Server shutdown initiated" {
		t.Errorf("Response Code: %v, Response: %v", a.Code, a.Body.String())
		t.Errorf("Correct Code: %v, Response: Server shutdown initiated", http.StatusOK)
	}
}