package main

import(
        "log"
        "net/url"
        "net/http"
        "net/http/httputil"

        "github.com/alexflint/go-arg"

        "berizaryad/green-proxy/utils"
)

var Args struct {
	Host    string  `help:"Service to be proxied (on localhost)" default:"http://localhost:8080"`
    Port    string  `help:"Port to serve proxy from" default:":3000"`
    Url     string  `help:"URL of the service to be proxied" default:"https://superset.berizaryad.ru"`
	ApiKey  string  `help:"API-key for /api/auth" default:"1234"`
    AuthUrl string  `help:"URL of the Authorization server" default:"https://password.berizaryad.ru"`
    AuthApi string  `help:"URL of the /api/auth handler" default:"https://password.berizaryad.ru/api/auth"`
}

func main() {
    log.SetPrefix("[LOG] ")
    log.SetFlags(3)

    p := arg.MustParse(&Args)

    localServer, err := url.Parse(Args.Host)
    if err != nil {
	    p.Fail("local resource must be in form of http://localhost:8080")
    }

    utils.InitConfig(Args.Url, Args.AuthUrl, Args.AuthApi, Args.ApiKey)

    log.Println("green-proxy started successfully on port ", Args.Port)
    log.Println("serving content from ", Args.Host)

    proxy := httputil.NewSingleHostReverseProxy(localServer)

    http.HandleFunc("/", utils.ProxyRedirect(proxy))

    log.Fatal(http.ListenAndServe(Args.Port, nil))
}
