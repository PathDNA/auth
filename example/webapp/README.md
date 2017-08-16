# webapp example

## Run:
```
➤ go build
➤ ./webapp -s -d # then open a new terminal


```

## Output:

### Server terminal:
```
apiserv: mw.go:45: [reqID:1] [127.0.0.1] [200] POST /api/v1/signup [134.665597ms]
        Headers: {"Accept-Encoding":["gzip"],"Content-Length":["100"],"User-Agent":["Go-http-client/1.1"]}
        Request (100): {"username":"one","password":"1234","profile":{"name":"name","phone":"phone","agency":{"fee":0.15}}}
apiserv: mw.go:45: [reqID:2] [127.0.0.1] [200] POST /api/v1/login [138.040469ms]
        Headers: {"Accept-Encoding":["gzip"],"Content-Length":["36"],"User-Agent":["Go-http-client/1.1"]}
        Request (36): {"username":"one","password":"1234"}
apiserv: mw.go:45: [reqID:3] [127.0.0.1] [200] GET /api/v1/profile [30.602µs]
apiserv: mw.go:45: [reqID:4] [127.0.0.1] [400] POST /api/v1/signup [143.12865ms]
        Headers: {"Accept-Encoding":["gzip"],"Content-Length":["100"],"User-Agent":["Go-http-client/1.1"]}
        Request (100): {"username":"one","password":"1234","profile":{"name":"name","phone":"phone","agency":{"fee":0.15}}}
apiserv: mw.go:45: [reqID:5] [127.0.0.1] [200] POST /api/v1/login [135.602234ms]
        Headers: {"Accept-Encoding":["gzip"],"Content-Length":["36"],"User-Agent":["Go-http-client/1.1"]}
        Request (36): {"username":"one","password":"1234"}
apiserv: mw.go:45: [reqID:6] [127.0.0.1] [200] GET /api/v1/profile [20.078µs]
apiserv: mw.go:45: [reqID:7] [127.0.0.1] [200] POST /api/v1/signup [136.859646ms]
        Headers: {"Accept-Encoding":["gzip"],"Content-Length":["112"],"User-Agent":["Go-http-client/1.1"]}
        Request (112): {"username":"one-adv","password":"1234","profile":{"name":"name","phone":"phone","advertiser":{"agencyID":"1"}}}
apiserv: mw.go:45: [reqID:8] [127.0.0.1] [200] POST /api/v1/login [134.285167ms]
        Headers: {"Accept-Encoding":["gzip"],"Content-Length":["40"],"User-Agent":["Go-http-client/1.1"]}
        Request (40): {"username":"one-adv","password":"1234"}
apiserv: mw.go:45: [reqID:9] [127.0.0.1] [200] GET /api/v1/profile [22.531µs]
```

### Client Terminal:
```
➤ ./webapp -c -signup
{"code":200,"data":"1","success":true}

➤ ./webapp -login "one:1234" -profile
{"code":200,"data":"1","success":true}
{"code":200,"data":{"id":"1","username":"one","status":1,"created":1502909125,"profile":{"name":"name","phone":"phone","agency":{"fee":0.15}}},"success":true}

➤ ./webapp -c -signup -signupInfo "one-adv:1234:name:phone:1"
{"code":200,"data":"2","success":true}

➤ ./webapp -login "one-adv:1234" -profile
{"code":200,"data":"2","success":true}
{"code":200,"data":{"id":"2","username":"one-adv","status":1,"created":1502909125,"profile":{"name":"name","phone":"phone","advertiser":{"agencyID":"1"}}},"success":true}
```
