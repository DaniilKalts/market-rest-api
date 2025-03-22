package jwt

import (
	"net/http"
	"strconv"

	"github.com/DaniilKalts/market-rest-api/internal/config"
)

func SetAuthCookies(w http.ResponseWriter, userID int) (accessToken, refreshToken string, err error) {
	accessToken, err = GenerateJWT(config.Config.Server.BaseURL, strconv.Itoa(userID), 15)
	if err != nil {
		return "", "", err
	}

	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		Domain:   config.Config.Server.Domain,
		MaxAge:   900,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, accessCookie)

	refreshToken, err = GenerateJWT(config.Config.Server.BaseURL, strconv.Itoa(userID), 1440)
	if err != nil {
		return "", "", err
	}

	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		Domain:   config.Config.Server.Domain,
		MaxAge:   86400,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, refreshCookie)

	return accessToken, refreshToken, nil
}

func DeleteAuthCookies(w http.ResponseWriter) error {
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Domain:   config.Config.Server.Domain,
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, accessCookie)

	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Domain:   config.Config.Server.Domain,
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, refreshCookie)

	return nil
}
