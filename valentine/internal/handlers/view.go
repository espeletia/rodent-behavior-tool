package handlers

import (
	"errors"
	"net/http"
	"time"
	"valentine/internal/usecases"
	"valentine/view"
)

type ViewHandler struct {
	client      http.Client
	userUsecase *usecases.UserUsecase
	cageUsecase *usecases.CageUsecase
}

func NewViewHandler(users *usecases.UserUsecase, cages *usecases.CageUsecase) *ViewHandler {
	return &ViewHandler{
		client:      http.Client{},
		userUsecase: users,
		cageUsecase: cages,
	}
}

func (vh *ViewHandler) Render(w http.ResponseWriter, r *http.Request) error {
	view.Thing().Render(r.Context(), w)
	return nil
}

func (vh *ViewHandler) App(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	usr, err := vh.userUsecase.Me(ctx)
	if err != nil {
		return err
	}

	cages, err := vh.cageUsecase.GetCages(ctx)
	if err != nil {
		return err
	}

	view.AppView(*usr, cages).Render(r.Context(), w)
	return nil
}

func (vh *ViewHandler) Login(w http.ResponseWriter, r *http.Request) error {
	view.LoginView().Render(r.Context(), w)
	return nil
}

func (vh *ViewHandler) HandleLoginForm(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	email := r.Form.Get("username")
	password := r.Form.Get("password")
	token, err := vh.userUsecase.Login(r.Context(), email, password)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    *token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(10 * time.Minute),
	})
	http.Redirect(w, r, "/app", http.StatusSeeOther)
	return nil
}

func GetAuthToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", errors.New("no auth token found")
		}
		return "", err
	}
	return cookie.Value, nil
}
