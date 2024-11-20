package handlers

import (
	"encoding/json"
	"net/http"
	"tusk/internal/domain"
	"tusk/internal/handlers/models"
	"tusk/internal/usecases"

	"go.uber.org/zap"
)

type UserHandler struct {
	userUsecase *usecases.UserUsecase
	authUsecase *usecases.AuthUsecase
}

func NewUserHandler(userusecase *usecases.UserUsecase, authUsecase *usecases.AuthUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userusecase,
		authUsecase: authUsecase,
	}
}

func (ph *UserHandler) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := struct {
			Message string `json:"message"`
		}{
			Message: "PONG!",
		}
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(400)
			return
		}
		return
	}
}

func (uu *UserHandler) Me() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, err := uu.userUsecase.Me(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		mappedUser := mapDomainUserToModelUser(*user)
		err = json.NewEncoder(w).Encode(mappedUser)
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}

func (uu *UserHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var creds models.LoginCreds
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		usr, err := uu.authUsecase.Login(ctx, mapModelUserCreds(creds))
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		zap.L().Info("HIT", zap.Any("places", usr))
		// Convert the result to JSON and write to the response
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(models.LoginResp{
			Token: usr,
		})
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}

	}
}

func (uu *UserHandler) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// Decode the JSON body to the `Viewport` struct
		var userData models.UserData
		err := json.NewDecoder(r.Body).Decode(&userData)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		zap.L().Info("userData", zap.Any("userData", userData))

		defer r.Body.Close()

		usr, err := uu.userUsecase.CreateUser(ctx, mapModelUserDataToDomainUserData(userData), userData.Password)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		zap.L().Info("HIT", zap.Any("places", usr))
		// Convert the result to JSON and write to the response
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(mapDomainUserToModelUser(*usr))
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}

	}
}

func mapModelUserDataToDomainUserData(usr models.UserData) domain.UserData {
	return domain.UserData{
		Username:    usr.Username,
		Email:       usr.Email,
		DisplayName: usr.DisplayName,
	}
}

func mapDomainUserToModelUser(usr domain.User) models.User {
	return models.User{
		ID:          usr.ID,
		DisplayName: usr.DisplayName,
		Email:       usr.Email,
		Username:    usr.Username,
	}
}

func mapModelUserCreds(usr models.LoginCreds) domain.LoginCreds {
	return domain.LoginCreds{
		Email:    usr.Email,
		Password: usr.Password,
	}
}
