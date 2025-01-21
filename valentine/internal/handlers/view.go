package handlers

import (
	"ghiaccio/models"
	"net/http"
	"time"
	"valentine/internal/usecases"
	"valentine/view"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ViewHandler struct {
	client       http.Client
	userUsecase  *usecases.UserUsecase
	cageUsecase  *usecases.CageUsecase
	videoUsecase *usecases.VideoUsecase
}

func NewViewHandler(users *usecases.UserUsecase, cages *usecases.CageUsecase, video *usecases.VideoUsecase) *ViewHandler {
	return &ViewHandler{
		client:       http.Client{},
		userUsecase:  users,
		cageUsecase:  cages,
		videoUsecase: video,
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

	videos, err := vh.videoUsecase.GetVideos(ctx)
	if err != nil {
		return err
	}
	view.AppView(*usr, cages, videos.Data).Render(r.Context(), w)
	return nil
}

func (vh *ViewHandler) CageView(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	cageMessages, err := vh.cageUsecase.GetCageMessages(ctx, id)
	if err != nil {
		return err
	}

	view.CageView(*cageMessages).Render(ctx, w)
	return nil
}

func (vh *ViewHandler) Login(w http.ResponseWriter, r *http.Request) error {
	view.LoginView().Render(r.Context(), w)
	return nil
}

func (vh *ViewHandler) About(w http.ResponseWriter, r *http.Request) error {
	view.AboutView().Render(r.Context(), w)
	return nil
}

func (vh *ViewHandler) Register(w http.ResponseWriter, r *http.Request) error {
	view.RegisterView().Render(r.Context(), w)
	return nil
}

func (vh *ViewHandler) HandleRegisterForm(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	err := r.ParseForm()
	if err != nil {
		return err
	}
	email := r.Form.Get("email")
	username := r.Form.Get("username")
	displayName := r.Form.Get("display_name")
	password := r.Form.Get("password")

	zap.L().Info("data",
		zap.String("username", username),
		zap.String("password", password),
		zap.String("display_name", displayName),
		zap.String("email", email))
	err = vh.userUsecase.Register(ctx, models.UserData{
		Email:       email,
		Username:    username,
		DisplayName: displayName,
		Password:    password,
	})
	if err != nil {
		return err
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return nil
}

func (vh *ViewHandler) HandleLoginForm(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	email := r.Form.Get("username")
	password := r.Form.Get("password")
	zap.L().Info("Parsing form")
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
