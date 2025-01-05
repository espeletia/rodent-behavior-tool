package handlers

import (
	"encoding/json"
	"net/http"
	"tusk/internal/domain"
	"tusk/internal/handlers/models"
	"tusk/internal/usecases"

	"tusk/internal/util"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
	"github.com/nextap-solutions/openapi3Struct"
)

type CagesHandler struct {
	cagesUsecase *usecases.CagesUsecase
}

func NewCagesHandler(cagesUsecase *usecases.CagesUsecase) *CagesHandler {
	return &CagesHandler{
		cagesUsecase: cagesUsecase,
	}
}

var CreateCageOp = openapi3Struct.Path{
	Path: "/cages",
	Item: openapi3.PathItem{
		Post: &openapi3.Operation{
			Tags:        []string{"Cages"},
			OperationID: "createCage",
			Description: "create a new inactive cage that will be pending activation",
			Responses: map[string]*openapi3.ResponseRef{
				"200": {
					Value: &openapi3.Response{
						Description: util.ToPointer("Cage creation response"),
						Content: map[string]*openapi3.MediaType{
							"application/json": {
								Schema: openapi3.NewSchemaRef("#/components/schemas/CageCreationResponse", nil),
							},
						},
					},
				},
			},
		},
	},
}

func (ch *CagesHandler) CreateCage(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	activationCode, err := ch.cagesUsecase.CreateNewCage(ctx)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(models.CageCreationResponse{
		ActivationCode: activationCode,
	})
	if err != nil {
		return err
	}
	return nil
}

var RegisterCageOp = openapi3Struct.Path{
	Path: "/activate/{code}",
	Item: openapi3.PathItem{
		Get: &openapi3.Operation{
			Tags:        []string{"Cages"},
			OperationID: "registerCage",
			Description: "register a cage under a user",
			Parameters: openapi3.Parameters{
				{
					Value: &openapi3.Parameter{
						Name:        "code",
						In:          "path",
						Description: "activation code for the cage",
						Required:    true,
						Schema:      openapi3.NewSchemaRef("#/components/schemas/Code", nil),
					},
				},
			},
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

func (ch *CagesHandler) RegisterCage(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	activationCode := mux.Vars(r)["code"]

	err := ch.cagesUsecase.RegisterCage(ctx, activationCode)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

	return err
}

var UserGetCagesOp = openapi3Struct.Path{
	Path: "/cages",
	Item: openapi3.PathItem{
		Get: &openapi3.Operation{
			Tags:        []string{"Cages"},
			OperationID: "userGetCages",
			Description: "fetch all cages that belong to an user",
			Responses: map[string]*openapi3.ResponseRef{
				"200": {
					Value: &openapi3.Response{
						Description: util.ToPointer("Cage creation response"),
						Content: map[string]*openapi3.MediaType{
							"application/json": {
								Schema: openapi3.NewSchemaRef("#/components/schemas/Cages", nil),
							},
						},
					},
				},
			},
		},
	},
}

func (ch *CagesHandler) GetCagesForUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	cages, err := ch.cagesUsecase.GetCagesForUser(ctx)
	if err != nil {
		return err
	}
	result := []models.Cage{}
	for _, cage := range cages {
		result = append(result, mapCageToModel(cage))
	}
	err = json.NewEncoder(w).Encode(models.Cages{Data: result})
	if err != nil {
		return err
	}
	return nil
}

func mapCageToModel(cage domain.Cage) models.Cage {
	return models.Cage(cage)
}
