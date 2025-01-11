package handlers

import (
	"encoding/json"
	"ghiaccio/models"
	"net/http"
	"tusk/internal/domain"
	"tusk/internal/usecases"
	"tusk/internal/util"

	"github.com/nextap-solutions/openapi3Struct"

	"github.com/getkin/kin-openapi/openapi3"
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

var PingOp = openapi3Struct.Path{
	Path: "/",
	Item: openapi3.PathItem{
		Get: &openapi3.Operation{
			Tags:        []string{"Ping"},
			OperationID: "ping",
			Description: "ping the API to test connectivity",
			Responses: map[string]*openapi3.ResponseRef{
				"204": {
					Value: &openapi3.Response{
						Description: util.ToPointer("No content"),
					},
				},
			},
		},
	},
}

func (ph *UserHandler) Ping(w http.ResponseWriter, r *http.Request) error {
	response := struct {
		Message string `json:"message"`
	}{
		Message: "PONG!",
	}
	err := json.NewEncoder(w).Encode(response)
	return err
}

var MeOp = openapi3Struct.Path{
	Path: "/v1/me",
	Item: openapi3.PathItem{
		Get: &openapi3.Operation{
			Tags:        []string{"Users"},
			OperationID: "me",
			Description: "get information about the user currently logged in",
			Responses: map[string]*openapi3.ResponseRef{
				"200": {
					Value: &openapi3.Response{
						Description: util.ToPointer("User"),
						Content: map[string]*openapi3.MediaType{
							"application/json": {
								Schema: openapi3.NewSchemaRef("#/components/schemas/User", nil),
							},
						},
					},
				},
			},
		},
	},
}

func (uu *UserHandler) Me(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	user, err := uu.userUsecase.Me(ctx)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	mappedUser := mapDomainUserToModelUser(*user)
	err = json.NewEncoder(w).Encode(mappedUser)
	if err != nil {
		return err
	}
	return nil
}

var LoginOp = openapi3Struct.Path{
	Path: "/v1/login",
	Item: openapi3.PathItem{
		Post: &openapi3.Operation{
			Tags:        []string{"Users"},
			OperationID: "login",
			Description: "request an access token based on user credentials",
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Description: "Add login credentials for the user",
					Required:    true,
					Content: map[string]*openapi3.MediaType{
						"application/json": {
							Schema: &openapi3.SchemaRef{
								Ref: "#/components/schemas/LoginCreds",
							},
						},
					},
				},
			},
			Responses: map[string]*openapi3.ResponseRef{
				"200": {
					Value: &openapi3.Response{
						Description: util.ToPointer("Token"),
						Content: map[string]*openapi3.MediaType{
							"application/json": {
								Schema: openapi3.NewSchemaRef("#/components/schemas/LoginResp", nil),
							},
						},
					},
				},
			},
		},
	},
}

func (uu *UserHandler) Login(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var creds models.LoginCreds
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	usr, err := uu.authUsecase.Login(ctx, mapModelUserCreds(creds))
	if err != nil {
		return err
	}
	// Convert the result to JSON and write to the response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(models.LoginResp{
		Token: usr,
	})
	if err != nil {
		return err
	}
	return nil
}

var CreateUserOp = openapi3Struct.Path{
	Path: "/v1/register",
	Item: openapi3.PathItem{
		Put: &openapi3.Operation{
			Tags:        []string{"Users"},
			OperationID: "createUser",
			Description: "create a new user",
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Description: "Add user data",
					Required:    true,
					Content: map[string]*openapi3.MediaType{
						"application/json": {
							Schema: &openapi3.SchemaRef{
								Ref: "#/components/schemas/UserData",
							},
						},
					},
				},
			},
			Responses: map[string]*openapi3.ResponseRef{
				"200": {
					Value: &openapi3.Response{
						Description: util.ToPointer("User"),
						Content: map[string]*openapi3.MediaType{
							"application/json": {
								Schema: openapi3.NewSchemaRef("#/components/schemas/User", nil),
							},
						},
					},
				},
			},
		},
	},
}

func (uu *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	// Decode the JSON body to the `Viewport` struct
	var userData models.UserData
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	usr, err := uu.userUsecase.CreateUser(ctx, mapModelUserDataToDomainUserData(userData), userData.Password)
	if err != nil {
		return err
	}
	// Convert the result to JSON and write to the response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(mapDomainUserToModelUser(*usr))
	if err != nil {
		return err
	}
	return err
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
