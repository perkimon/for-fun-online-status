# for-fun-online-status
An online status service supporting UDP and TCP connections.  Written as an expression of a test challenge.

V1
* TCP on http://0.0.0.0:2000
* UDP on port 2000

Usage

go run *.go

In 4 different terminal windows run the following commands:
```

//To test UDP run  (UDP will timeout after 30 seconds without sending another request)
nc -u localhost 2000
{"user_id": 1, "friends": [2, 3, 4]}

//to avoid timeout send a "Wave" before 30 seconds is up.
{"user_id": 1, "action":3}

//To test TCP (no wave needed)
nc localhost 2000
{"user_id": 2, "friends": [1, 3, 4]}


```

TCP is stateful so detects disconnections, UDP clients are deemed offline after 30 seconds with no updates.
