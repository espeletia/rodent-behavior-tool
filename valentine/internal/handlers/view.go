package handlers

import (
	"ghiaccio/models"
	"log"
	"net/http"
	"path/filepath"
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
	fetcher      *usecases.RestFetcher
}

func NewViewHandler(users *usecases.UserUsecase, cages *usecases.CageUsecase, video *usecases.VideoUsecase, fetcher *usecases.RestFetcher) *ViewHandler {
	return &ViewHandler{
		client:       http.Client{},
		userUsecase:  users,
		cageUsecase:  cages,
		videoUsecase: video,
		fetcher:      fetcher,
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

func (vh *ViewHandler) MessageView(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	cageId := mux.Vars(r)["id"]
	log.Printf("|%s|", cageId)
	messageId := mux.Vars(r)["message"]
	log.Printf("|%s|", messageId)

	cageMessage, err := vh.cageUsecase.GetCageMessage(ctx, cageId, messageId)
	if err != nil {
		return err
	}

	view.MessageView(*cageMessage).Render(ctx, w)
	return nil
}

func (vh *ViewHandler) Login(w http.ResponseWriter, r *http.Request) error {
	view.LoginView().Render(r.Context(), w)
	return nil
}

func (vh *ViewHandler) UploadAnalysis(w http.ResponseWriter, r *http.Request) error {
	view.UploadFormView().Render(r.Context(), w)
	return nil
}

func (vh *ViewHandler) RegisterCage(w http.ResponseWriter, r *http.Request) error {
	view.LinkCageView().Render(r.Context(), w)
	return nil
}

func (vh *ViewHandler) HandleUploadAnalysis(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	// Parse multipart form
	err := r.ParseMultipartForm(32 << 20) // 32 MB
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return err
	}

	// Get text fields
	name := r.FormValue("name")
	var description *string = nil
	descriptionForm := r.FormValue("description")
	if descriptionForm != "" {
		description = &descriptionForm
	}

	// Get the uploaded file
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "MP4 file is required", http.StatusBadRequest)
		return err
	}
	http.Redirect(w, r, "/app", http.StatusSeeOther)

	go func() {
		defer file.Close()

		// Optional: Validate extension
		if filepath.Ext(header.Filename) != ".mp4" {
			http.Error(w, "Only .mp4 files are allowed", http.StatusBadRequest)
			return
		}
		zap.L().Info("redirecting")

		zap.L().Info("uploading")
		url, err := vh.fetcher.GenericFileUpload(ctx, file, header.Filename)
		if err != nil {
			return
		}

		videoDto := models.CreateVideoDto{
			Name:        name,
			Description: description,
			VideoUrl:    url,
		}

		zap.L().Info("creting video")
		err = vh.videoUsecase.CreateVideo(ctx, videoDto)
		if err != nil {
			return
		}
	}()

	return nil
}

func (vh *ViewHandler) HandleCageLink(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	// Get text fields
	code := r.FormValue("code")

	zap.L().Info("creting video")
	err := vh.cageUsecase.ActivateCage(ctx, code)
	if err != nil {
		return err
	}

	http.Redirect(w, r, "/app", http.StatusSeeOther)

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
