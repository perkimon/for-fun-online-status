# for-fun-online-status
An online status service supporting UDP and TCP connections.  Written as an expression of a test challenge.

V1
* HTTP on http://0.0.0.0:2000/status (Using curl it is possible to stream updates on when your friends are online)
* UDP on port 2000
* Accepts POST requests to /status and uses text/event-stream to stream json messages back to the client 
* Status Tracker has tests

Usage

go run *.go

In 4 different terminal windows run the following commands:
```
curl -X POST -d '{"user_id": 1, "friends": [2, 3, 4]}' http://localhost:2000/status
curl -X POST -d '{"user_id": 2, "friends": [1, 3, 4]}' http://localhost:2000/status
curl -X POST -d '{"user_id": 3, "friends": [1, 2, 4]}' http://localhost:2000/status
curl -X POST -d '{"user_id": 4, "friends": [1, 2, 3]}' http://localhost:2000/status

//To test UDP run 
nc -u localhost 2000
{"user_id": 1, "friends": [2, 3, 4]}

```