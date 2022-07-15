package main

import(
        "log"
        "net/url"
        "net/http"
        "net/http/httputil"
        "os"

        "github.com/koterin/green-proxy/tree/master/utils"
)

func main() {
    log.SetPrefix("[LOG] ")
    log.SetFlags(3)

    localhost := os.Getenv("PROXY_LOCALHOST")
    proxyPort := os.Getenv("PROXY_PORT")

    localServer, err := url.Parse(localhost)
    if err != nil {
        panic(err)
    }

    log.Println("green-proxy started successfully on port ", proxyPort)
    log.Println("serving content from ", localhost)

    proxy := httputil.NewSingleHostReverseProxy(localServer)

    http.HandleFunc("/", utils.ProxyRedirect(proxy))

    log.Fatal(http.ListenAndServe(proxyPort, nil))
}
