package main

import(
        "log"
        "net/url"
        "net/http"
        "net/http/httputil"
//      "os"

        "ktrn.com/greenProxy"
)

//var authServerUrl = os.Getenv("HOST_URL")
var Localhost = "http://localhost:8080"

func main() {
    log.SetPrefix("[LOG] ")
    log.SetFlags(3)


    localServer, err := url.Parse(Localhost)
    if err != nil {
        panic(err)
    }

    log.Println("green-proxy started successfully on port :3000")

    proxy := httputil.NewSingleHostReverseProxy(localServer)

    http.HandleFunc("/", greenProxy.ProxyRedirect(proxy))

    log.Fatal(http.ListenAndServe(":3000", nil))
}
