package transport

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"auth/pkg/app/provider"
	"auth/pkg/app/service"
)

type Handler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	UpdateToken(w http.ResponseWriter, r *http.Request)
}

func NewHandler(
	tokenService service.TokenService,
	userService *service.UserService,
	tokenProvider provider.TokenProvider,
) Handler {
	return &handler{
		tokenService:  tokenService,
		userService:   userService,
		tokenProvider: tokenProvider,
	}
}

type handler struct {
	tokenService  service.TokenService
	userService   *service.UserService
	tokenProvider provider.TokenProvider
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	req := registerBody{}
	err = json.Unmarshal(body, &req)

	login := req.Email
	password := req.Password

	err = h.userService.CreateUser(login, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.login(w, login)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	req := loginBody{}
	err = json.Unmarshal(body, &req)

	login := req.Email
	password := req.Password

	isAuth, err := h.userService.Authenticate(login, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isAuth {
		http.Error(w, errors.New("password not matched").Error(), http.StatusUnauthorized)
		return
	}

	h.login(w, login)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	h.logout(w, r)
}

func (h *handler) UpdateToken(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "token not found", http.StatusUnauthorized)
		return
	}

	login, err := h.tokenService.ParseToken(tokenCookie.Value)
	if err != nil {
		http.Error(w, "token is invalid", http.StatusUnauthorized)
		return
	}

	savedToken, err := h.tokenProvider.GetTokenByLogin(login)
	if err != nil {
		http.Error(w, "token not exists", http.StatusUnauthorized)
		return
	}

	if savedToken != tokenCookie.Value {
		http.Error(w, "token is invalid", http.StatusUnauthorized)
		return
	}

	h.login(w, login)
	w.WriteHeader(http.StatusOK)
}

func (h *handler) logout(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("token")
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	deleteCookie(w, "token")

	login, err := h.tokenService.ParseToken(tokenCookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	err = h.tokenService.DeleteToken(login)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}
}

func (h *handler) login(w http.ResponseWriter, login string) {
	accessToken, accessExpirationTime, err := h.tokenService.CreateToken(login, 5*time.Hour)
	if err != nil {
		http.Error(w, errors.New("could not create token").Error(), http.StatusUnauthorized)
		return
	}
	err = h.tokenService.SaveToken(accessToken, login)
	if err != nil {
		http.Error(w, errors.New("could save token").Error(), http.StatusUnauthorized)
		return
	}

	setCookie(w, "token", accessToken, accessExpirationTime)
}

func setCookie(w http.ResponseWriter, name, value string, expirationTime time.Time) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expirationTime,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}

func deleteCookie(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}
