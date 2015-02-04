package main

import (
    "net/http"
    "fmt"
    "log"
    "os"
    "io"
    "runtime"

    // go get github.com/gorilla/mux
    "github.com/gorilla/mux"
)

const (
    HOST = "localhost"
)

func Handler_404(w http.ResponseWriter, r *http.Request){
    fmt.Fprint(w, "Oops, something went wrong!")
}

func Handler_www(w http.ResponseWriter, r *http.Request){
    fmt.Fprint(w, "Hello world :)")
}

func Handler_api(w http.ResponseWriter, r *http.Request){
    fmt.Fprint(w, "This is the API")
}

func Handler_secure(w http.ResponseWriter, r *http.Request){
    fmt.Fprint(w, "This is Secure")
}

func redirect(r *mux.Router, from string, to string){
    r.Host(from).Subrouter().HandleFunc("/", func (w http.ResponseWriter, r *http.Request){
        http.Redirect(w, r, to, 301)
    })
}

func main(){
    port := 9000
    ssl_port := 443

    runtime.GOMAXPROCS(runtime.NumCPU())

    http_r := mux.NewRouter()
    https_r := mux.NewRouter()

    //  HTTP 404
    http_r.NotFoundHandler = http.HandlerFunc(Handler_404)

    //  Redirect "http://HOST" => "http://www.HOST"
    redirect(http_r, HOST, fmt.Sprintf("http://www.%s:%d", HOST, port))

    //  Redirect "http://secure.HOST" => "https://secure.HOST"
    redirect(http_r, "secure."+HOST, fmt.Sprintf("https://secure.%s", HOST))

    www := http_r.Host("www."+HOST).Subrouter()
    www.HandleFunc("/", Handler_www)

    api := http_r.Host("api."+HOST).Subrouter()
    api.HandleFunc("/", Handler_api)

    secure := https_r.Host("secure."+HOST).Subrouter()
    secure.HandleFunc("/", Handler_secure)

    //  Start HTTP
    err_http := http.ListenAndServe(fmt.Sprintf(":%d", port), http_r)
    if err_http != nil {
        log.Fatal("Web server (HTTP): ", err_http)
    }

    //  Start HTTPS
    err_https := http.ListenAndServeTLS(fmt.Sprintf(":%d", ssl_port), "D:/Go/src/www/ssl/public.crt", "D:/Go/src/www/ssl/private.key", https_r)
    if err_https != nil {
        log.Fatal("Web server (HTTPS): ", err_https)
    }
}