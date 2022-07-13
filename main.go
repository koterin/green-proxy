package main

import(
        "log"
        "net/url"
        "net/http"
        "net/http/httputil"
  //      "os"
)

//var authServerUrl = os.Getenv("HOST_URL")
var authServerUrl = "https://password.berizaryad.ru"
var resourcePubHost = "http://supernet.ktrn.com"
var resourceHost = "51.250.27.68:8080/"

func main() {
    log.SetPrefix("[LOG] ")
    log.SetFlags(3)


    authServer, err := url.Parse(authServerUrl)
    if err != nil {
        panic(err)
    }

    log.Println("green-proxy started successfully")

    proxy := httputil.NewSingleHostReverseProxy(authServer)

    http.HandleFunc("/", proxyRedirect(proxy))

    log.Fatal(http.ListenAndServe(":3000", nil))
}

func proxyRedirect(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        redirect := r.URL.Scheme + resourcePubHost + r.URL.Path

        log.Println("SEND BACK TO ", redirect)

        auth := false
        if !auth {
            w.Header().Set("X-Redirect-To", redirect)
            http.Redirect(w, r, authServerUrl, http.StatusSeeOther)
        } else {
            r.URL.Scheme = "https"
            r.Host = resourceHost
            w.Header().Set("X-Redirect-To", redirect)
            w.Header().Set("X-Session", "1234")
            p.ServeHTTP(w, r)
        }
    }
}
