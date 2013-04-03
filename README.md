gae-channel
===========

Google App Engine  Channel Service Client Implemented in Golang

# Setup a channel service on GAE

hello.go
```
package hello

import (
    "fmt"
    "net/http"

    "appengine"
    "appengine/channel"
)

const (
    TOPIC = "counter"
)

func init() {
    http.HandleFunc("/new_token", handler_new_token)    
}

func handler_new_token(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    tok, err := channel.Create(c, TOPIC)    
    callback := r.FormValue("callback") 
    if err != nil {
        http.Error(w, "Couldn't create Channel", http.StatusInternalServerError)
        c.Errorf("channel.Create: %v", err)
        return
    }
    if callback == "" {
        w.Header().Set("Content-type", "text/javascript")
        fmt.Fprintf(w, "%s", tok)
    } else {
        fmt.Fprintf(w, callback+"('%s')", tok)
    }
}
```

app.yaml

```
application: <your_app_name>
version: 1
runtime: go
api_version: go1

handlers:
- url: (/.*)
  script: _go_app
```

# Sample Client code
```
package main

import (
    "encoding/json"
    . "github.com/hugozhu/gae-channel"
    "log"
)

func main() {
    log.Println("started")
    stop_chan := make(chan bool)

    channel := NewChannel("http://<your_app_id>.appspot.com/new_token")
    socket := channel.Open()
    socket.OnOpened = func() {
        log.Println("socket opened!")
    }

    socket.OnClose = func() {
        log.Println("socket closed!")
        stop_chan <- true
    }

    socket.OnMessage = func(msg *Message) {
        if msg.Level() >= 3 && msg.Child.Key == "c" {
            v1 := *msg.Child.Child.Val
            if len(v1) > 0 {
                s := "[" + v1[0].Key + "]"
                var v []string
                json.Unmarshal([]byte(s), &v)
                if len(v) == 2 && v[0] == "ae" {
                    s = v[1]
                    log.Println(s) //what really intersting
                }
            }
        }
    }

    socket.OnError = func(err error) {
        log.Println("error:", err)
    }

    <-stop_chan
}
```

# See also

1. http://hugozhu.myalert.info/2013/04/03/24-google-channel-service.html