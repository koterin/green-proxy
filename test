server () {
  + 1. user asks for supernet
    
    2. request is being sent with the cookie header

    3. proxy sees it
    
    4. if it redirect from password.berizaryad.ru, it sets the cookie sent via X-Green-Token

    5. checks cookie - sends api/Auth, copies the cookie

    6. password answers 401 or 200

    7. proxy recieves 200 and serves the page

    8. proxy receives 401 and redirects user to password

    9.with the header X-Redirect-to: requested path

    10. password check for the header X-Redirect-to; if it exists, it puts it into the Session Storage

    11. user completes authorization successdfully:
    
    12. after authenticate.html js checks if there's redirect-to in Session Storage
    
    13. if there is, sessionId is not set to the password, it is being added to the header of redirect
    
    14. redirect link is deleted from the Session Store

    15. header X-Green-Token is being sent with the redirect
}
