package main

import(
        "log"
        "net/url"
        "net/http"
        "net/http/httputil"
        "flag"

        "github.com/koterin/green-proxy/tree/master/utils"
)

var localhost string
var proxyPort string

func main() {
    log.SetPrefix("[LOG] ")
    log.SetFlags(3)

    flag.StringVar(&localhost,
                    "host",
                    "localhost:8080",
                    "address of the resource being proxied")
    flag.StringVar(&proxyPort,
                    "port",
                    ":3000",
                    "port to run proxy on")
    flag.StringVar(&utils.PublicUrl,
                    "url",
                    "https://superset.berizaryad.ru",
                    "public URL of the resource to be proxied")
    flag.StringVar(&utils.AuthServerUrl,
                    "authUrl",
                    "https://password.berizaryad.ru",
                    "URL of the AuthServer")
    flag.StringVar(&utils.AuthApiUrl,
                    "authapi",
                    "https://password.berizaryad.ru/api/auth",
                    "URL of the /api/auth handler")
    flag.Parse()

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
