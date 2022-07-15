package main

import(
        "log"
        "net/url"
        "net/http"
        "net/http/httputil"

        "github.com/koterin/green-proxy/tree/master/utils"
)

// TODO: this has to be env
var Localhost = "http://localhost:8080"
var proxyPort = ":3000"

func main() {
    log.SetPrefix("[LOG] ")
    log.SetFlags(3)


    localServer, err := url.Parse(Localhost)
    if err != nil {
        panic(err)
    }

    log.Println("green-proxy started successfully on port ", proxyPort)
    log.Println("serving content from ", Localhost)

    proxy := httputil.NewSingleHostReverseProxy(localServer)

    http.HandleFunc("/", utils.ProxyRedirect(proxy))

    log.Fatal(http.ListenAndServe(proxyPort, nil))
}
