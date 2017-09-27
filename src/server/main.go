package main

import (
    "log"
    "flag"
    "net/http"
    "go/build"
    "path"
)

var addr = flag.String("addr", ":8080", "listen port")

func httpHandler(w http.ResponseWriter, r *http.Request) {
    log.Println(r.URL)

    if r.URL.Path != "/" {
        http.Error(w, "Not found", 404)
        return
    }

    if r.Method != "GET" {
        http.Error(w, "Method not allowed", 405)
        return
    }

    http.ServeFile(w, r, path.Join(build.Default.GOPATH, "www/index.html"))

}

func main() {
    flag.Parse()
    http.HandleFunc("/", httpHandler)
    http.ListenAndServe(*addr, nil)
}