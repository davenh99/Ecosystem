package tools

import (
	"net/http"
)

func SetCookie(w http.ResponseWriter, value string)  {
	cookie := http.Cookie{
		Name: "jwt_cookie",
		Value: value,
		Secure: false, // https, false for prod make sure to set to true later...
		HttpOnly: true, // this one to disallow javascript?
		SameSite: http.SameSiteLaxMode, // again, just for prod, set to strict
		Path: "/",
		// Domain: "localhost", // I hope this works
		MaxAge: 31536000000, // 1 year
	}
	
	http.SetCookie(w, &cookie)
}

// func GetCookie(r *http.Request) (string, error) {
// 	cookie, err := r.Cookie("jwt_refresh")
// 	// cookieStr := cookie.Value
// 	return cookie, err
// }
