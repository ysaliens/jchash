# jchash

A Go encryption HTTP server

## About
jchash is a multi-threaded HTTP web server capable of encrypting passwords using SHA-512. It can
* Encrypt a password given a "password=value" string with a POST HTTP request and store it
* Retrieve a hashed password given an ID and a GET HTTP request
* Display total hash requests served and average hash request response time in ms given a /stats command
* Shutdown remotely given a /shutdown command

## Setup
Install [Go](https://golang.org/)

Clone git
```
mkdir /drives/c/Users/Steve/go/src/github.com/ysaliens/jchash
cd /drives/c/Users/Steve/go/src/github.com/ysaliens/jchash
git clone https://github.com/ysaliens/jchash.git .
```
Build project
```
cd /drives/c/Users/Steve/go/src/github.com/ysaliens/jchash
go build
```
## Running the server
The server can be started by running the executable ./jchash.exe. Port can be specified by `--port=$PORT`. If no port is specified, the server will default to 8080. 
```
➤ ./jchash.exe --port=8000
2018/03/11 16:04:30 Listening on http://localhost:8000
```

## Hashing a password
To hash a password, send a `POST` request such as `curl --data "password=$PASSWORD" http://localhost:8080/hash` where $PASSWORD is the password value. The `/hash` at the end instructs the server to hash the password using SHA-512 and store it. If the is processed, the server will send status code 200 (OK) and return a numeric ID that can be used to retrieve the hash later. A hash command should return immediately, however the password will not be written to the database for 5 seconds.
```
➤ curl --data "password=test2" http://localhost:8000/hash
1  
```
## GET a hashed password
To retrieve a hash, send a `GET` such as `curl http://localhost:8000/hash/$ID` where $ID is the numeric ID returned after a hash operation. A successful `GET` will return a status code 200 (OK) and a string of the hashed password. If a `GET` command is issued for an ID that is still being written to the database, the server will re-try reading the hash every second for up to 10 seconds (double the time of the slowest hash write operation). 
```
➤ curl http://localhost:8000/hash/1
bSAb7u+1ibCO8GctrII1PQy9mtmeFkLIOhYB89ZHvMoAMle16PMb3B1z++yE+whcedbiZ3t/+SfoI6VOeJFA2Q==
```
## Shut down server
There are two ways to shut the server down - local or remote. Locally, a Control+C or SIGTERM signal will tell the server to finish write operations and stop. To shut the server down remotely, issue a /shutdown such as 
```
Client:
➤ curl http://localhost:8000/shutdown
Server shutdown initiated

Server log:
2018/03/11 18:19:35 Received shutdown request, shutting down.
2018/03/11 18:19:41 Server gracefully shut down.
```
In both cases, the server will wait 6 seconds to allow for any database requests to finish writing.

## Performance Stats
To check how many hash requests (hash a password or retrieve a stored hash) the server has processed, issue a /stats command
```
➤ curl http://localhost:8000/stats
{"total":5,"average":0}
```
The server will return a JSON object with the total number of hash requests processed thus far (including errors) and the average time to respond to a hash request in milliseconds. The server does not count /stats and /shutdown commands in these statistics. It also does not count the 5 seconds it takes for a hash `POST` to be written to the database as that is done in the background and does not affect response time from a client perspective. A quick note - the 0 ms average request time most often seen is correct - most hash commands take a fraction of a millisecond. Putting sleep, log, or print statements in the hash handlers will increase the average request time.


