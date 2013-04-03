gae-channel
===========

Google App Engine  Channel Service Client Implemented in Golang

# Setup a channel service on GAE

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


# Sample Client code

```
package main

import (
    . "github.com/hugozhu/gae-channel"
    "log"
)

func main() {
    log.Println("started")
    stop_chan := make(chan bool)
    channel := NewChannel("http://<your_app_name>.appspot.com/new_token")
    socket := channel.Open()
    socket.OnOpened = func() {
        log.Println("socket opened!")
    }

    socket.OnClose = func() {
        log.Println("socket closed!")
        stop_chan <- true
    }

    socket.OnMessage = func(msg *Message) {
        log.Println(msg.ToString())
    }

    socket.OnError = func(err error) {
        log.Println("error:", err)
    }

    <-stop_chan
}
```