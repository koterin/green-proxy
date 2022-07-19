package utils

import(
        "net/http"
        "net/http/httputil"
        "time"
        "net/url"
        "log"
)

var (
        PublicUrl       string
        AuthServerUrl   string
        AuthApiUrl      string
        API_KEY         string
        hClient = &http.Client{Timeout: 10 * time.Second}
)

func ProxyRedirect(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {

        token := r.URL.Query().Get("greenToken")
        if token == "" {
            if sessionValid := checkSessionCookie(r); !sessionValid {
                redirectToAuthServer(w, r)
                return
            }
            serveContent(w, r, proxy)

            return
        }

        if tokenIsValid := checkToken(token); !tokenIsValid {
			redirectToAuthServer(w, r)
            return
        }
        setCookie(w, token)
        redirectWithoutToken(w, r)
    }
}

func redirectToAuthServer(w http.ResponseWriter, r *http.Request) {
    redirect := r.URL.Scheme + r.URL.Host + r.URL.Path
    link := AuthServerUrl + "?redirect=" + redirect
    log.Println("host is", r.URL.Host)
    log.Println("link to resirect is", link)
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
    http.Redirect(w, r, newLink, http.StatusSeeOther)
}

func deleteTokenFromUrl(u *url.URL) {
    q := u.Query()
	q.Del("greenToken")
	u.RawQuery = q.Encode()
}

func checkToken(token string) bool {
    req, err := http.NewRequest("GET", AuthApiUrl, nil)
    if err != nil {
        log.Println("Error while creating request to /auth/api: ", err)
        return false
    }

    setupRequest(req, token)

    resp, err := hClient.Do(req)
    if err != nil {
        logRequest(req)
        return false
    }

    if resp.StatusCode == http.StatusCreated {
        return true
    }

    return false
}

func setupRequest(req *http.Request, token string) {
    addBasicReqHeaders(req)
    addOrigin(req)
    addApiKey(req)
    req.AddCookie(createCookie(token))
}

func logRequest(req *http.Request) {
    log.Println("Error while sending request to /auth/api")
    log.Println("---DEBUG---")

    dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
        log.Fatal(err)
    }
    log.Println(dump)
}

func checkSessionCookie(req *http.Request) bool {
    sessionCookie, err := req.Cookie("sessionId")
    if err != nil {
        return false
    }

    return checkToken(sessionCookie.Value)
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

func addBasicReqHeaders(req *http.Request) {
    req.Header.Set("Method", "GET")
    req.Header.Set("Accept", "*/*")
    req.Header.Set("Access-Control-Allow-Origin", PublicUrl)
    req.Header.Set("Access-Control-Allow-Credentials", "true")
    req.Header.Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
    req.Header.Set("Content-Type", "application/json")
}

func addOrigin(req *http.Request) {
    originUrl, _ := url.Parse(PublicUrl)
    req.Header.Set("X-Green-Origin", originUrl.Host)
}

func addApiKey(req *http.Request) {
    req.Header.Set("Api-Key", API_KEY)
}

func InitConfig(url string, authUrl string, authApi string, apiKey string) {
    PublicUrl = url
    AuthServerUrl = authUrl
    AuthApiUrl = authApi
    API_KEY = apiKey
}
