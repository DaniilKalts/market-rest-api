package jwt

import (
	"net/http"

	"github.com/DaniilKalts/market-rest-api/internal/config"
)

func SetCookie(w http.ResponseWriter, name, value, domain string, maxAge int, secure, httpOnly bool, sameSite http.SameSite) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   domain,
		MaxAge:   maxAge,
		Secure:   secure,
		HttpOnly: httpOnly,
		SameSite: sameSite,
	}
	http.SetCookie(w, cookie)
}

func SetAuthCookies(w http.ResponseWriter, accessToken, refreshToken string) error {
	SetCookie(w, "access_token", accessToken, config.Config.Server.Domain, 900, true, true, http.SameSiteLaxMode)
	SetCookie(w, "refresh_token", refreshToken, config.Config.Server.Domain, 86400, true, true, http.SameSiteLaxMode)

	return nil
}

func DeleteAuthCookies(w http.ResponseWriter) error {
	SetCookie(w, "access_token", "", config.Config.Server.Domain, -1, true, true, http.SameSiteLaxMode)
	SetCookie(w, "refresh_token", "", config.Config.Server.Domain, -1, true, true, http.SameSiteLaxMode)

	return nil
}
