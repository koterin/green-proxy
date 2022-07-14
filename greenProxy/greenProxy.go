package greenProxy

import(
        "log"
        "net/http"
        "net/http/httputil"
        "time"
        "net/url"
        //"os"
)

// TODO: this has to be env
var AuthServerUrl = "https://password.berizaryad.ru"
var AuthApiUrl = AuthServerUrl + "/api/auth/"
var PublicUrl = "https://swagger.berizaryad.ru"
var hClient = &http.Client{Timeout: 10 * time.Second}

func ProxyRedirect(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        var auth bool

        token := r.URL.Query().Get("greenToken")
        if token != "" {
            log.Println("24: token is present")
            auth = checkToken(token)
            if auth {
                setCookie(w, token)
                redirectWithoutToken(w, r)

                return
            }
        } else {
            log.Println("no GreenToken, sending with browser cookies to checkSessionCookie")
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
    http.Redirect(w, r, link, http.StatusSeeOther)
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
    http.Redirect(w, r, newLink, 301)
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

    log.Println("making checkToken request to ", req.URL)

    req.AddCookie(&http.Cookie{Name: "sessionId", Value: token})
    log.Println("sending cookie ", token)

    resp, err := hClient.Do(req)
    if err != nil {
        return false
    }

    log.Println("code: ", resp.StatusCode)
    log.Println("status: ", resp.Status)
    if resp.StatusCode == 201 {
        return true
    }

    return false
}

func checkSessionCookie(req *http.Request) bool {
    sessionCookie, err := req.Cookie("sessionId")
    if err != nil {
        log.Println("checking browser cookie - it's not here")
        return false
    }
    log.Println("checking browser cookie - it's here and it's ", sessionCookie.Value)

    return checkToken(sessionCookie.Value)
}

func CorsHandler(w http.ResponseWriter, req *http.Request) {
    AddBasicHeaders(w)
    w.WriteHeader(http.StatusOK)

    return
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
