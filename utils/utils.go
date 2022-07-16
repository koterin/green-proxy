package utils

import(
        "net/http"
        "net/http/httputil"
        "time"
        "net/url"
        "log"
)

var PublicUrl string
var AuthServerUrl string
var AuthApiUrl string
var hClient = &http.Client{Timeout: 10 * time.Second}

func ProxyRedirect(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        var auth bool

        token := r.URL.Query().Get("greenToken")
        if token != "" {
            auth = checkToken(token)
            if auth {
                setCookie(w, token)
                redirectWithoutToken(w, r)

                return
            }
        } else {
            auth = checkSessionCookie(r)
        }

        if !auth {
            redirectToAuthServer(w, r)
        } else {
            serveContent(w, r, proxy)
        }
    }
}

func redirectToAuthServer(w http.ResponseWriter, r *http.Request) {
    redirect := r.URL.Scheme + PublicUrl + r.URL.Path
    link := AuthServerUrl + "?redirect=" + redirect
    http.Redirect(w, r, link, 302)
}

func serveContent(w http.ResponseWriter, r *http.Request, proxy *httputil.ReverseProxy) {
    proxy.ServeHTTP(w, r)
}

func setCookie(w http.ResponseWriter, token string) {
    cookie := createCookie(token)
    http.SetCookie(w, cookie)
}

func redirectWithoutToken(w http.ResponseWriter, r *http.Request) {
    deleteTokenFromUrl(r.URL)

    newLink := r.URL.Scheme + r.URL.Host + r.URL.RawQuery
    http.Redirect(w, r, newLink, 302)
}

func deleteTokenFromUrl(u *url.URL) {
    q := u.Query()
	q.Del("greenToken")
	u.RawQuery = q.Encode()
}

func checkToken(token string) bool {
    req, err := http.NewRequest("GET", AuthApiUrl, nil)
    if err != nil {
        return false
    }

    req.AddCookie(&http.Cookie{Name: "sessionId", Value: token})

    AddBasicReqHeaders(req)
    log.Println("req host is ", req.Host)

    resp, err := hClient.Do(req)
    if err != nil {
        return false
    }

    if resp.StatusCode == 201 {
        return true
    }

    return false
}

func checkSessionCookie(req *http.Request) bool {
    sessionCookie, err := req.Cookie("sessionId")
    if err != nil {
        return false
    }

    return checkToken(sessionCookie.Value)
}

func CorsHandler(w http.ResponseWriter, req *http.Request) {
    AddBasicHeaders(w)
    w.WriteHeader(http.StatusOK)

    return
}

func AddBasicReqHeaders(req *http.Request) {
    req.Header.Set("Method", "GET")
    req.Header.Set("Accept", "*/*")
    req.Header.Set("Access-Control-Allow-Origin", PublicUrl)
    req.Header.Set("Access-Control-Allow-Credentials", "true")
    req.Header.Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
    req.Header.Set("Content-Type", "application/json")
}

func AddBasicHeaders(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", PublicUrl)
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
    w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
    w.Header().Set("Content-Type", "application/json")
}

func createCookie(token string) *http.Cookie {
    return &http.Cookie{
        Name:   "sessionId",
        Value:  token,
        Expires: (time.Now().Add(time.Duration(168) * time.Hour)),
        Path: "/",
        HttpOnly: true,
        Secure: true,
    }
}
