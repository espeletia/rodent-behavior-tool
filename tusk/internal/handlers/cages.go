package handlers

import (
	"encoding/json"
	"ghiaccio/models"
	"net/http"
	"time"
	"tusk/internal/domain"
	"tusk/internal/usecases"

	"tusk/internal/util"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/uuid"
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
	Path: "/v1/cages",
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
	activationCode, secretToken, err := ch.cagesUsecase.CreateNewCage(ctx)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(models.CageCreationResponse{
		ActivationCode: activationCode,
		SecretToken:    secretToken,
	})
	if err != nil {
		return err
	}
	return nil
}

var RegisterCageOp = openapi3Struct.Path{
	Path: "/v1/activate/{code}",
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
	Path: "/v1/cages",
	Item: openapi3.PathItem{
		Get: &openapi3.Operation{
			Tags:        []string{"Cages"},
			OperationID: "userGetCages",
			Description: "fetch all cages that belong to an user",
			Responses: map[string]*openapi3.ResponseRef{
				"200": {
					Value: &openapi3.Response{
						Description: util.ToPointer("Cages response"),
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
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(models.Cages{Data: result})
	if err != nil {
		return err
	}
	return nil
}

var CageGetSelfOp = openapi3Struct.Path{
	Path: "/internal/cage",
	Item: openapi3.PathItem{
		Get: &openapi3.Operation{
			Tags:        []string{"CagesInternal"},
			OperationID: "cageGetSelg",
			Description: "cage fetches itself to get information from server",
			Responses: map[string]*openapi3.ResponseRef{
				"200": {
					Value: &openapi3.Response{
						Description: util.ToPointer("Cage response"),
						Content: map[string]*openapi3.MediaType{
							"application/json": {
								Schema: openapi3.NewSchemaRef("#/components/schemas/Cage", nil),
							},
						},
					},
				},
			},
		},
	},
}

func (ch *CagesHandler) CageSelf(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	cage, err := ch.cagesUsecase.CageSelf(ctx)
	if err != nil {
		return err
	}
	mappedCage := mapCageToModel(*cage)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(mappedCage)
	if err != nil {
		return err
	}
	return nil
}

var CageSendMessageOp = openapi3Struct.Path{
	Path: "/internal/cage/message",
	Item: openapi3.PathItem{
		Post: &openapi3.Operation{
			Tags:        []string{"CagesInternal"},
			OperationID: "cageSendMessage",
			Description: "cage sends a message to server",
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Description: "Message payload from cage",
					Required:    true,
					Content: map[string]*openapi3.MediaType{
						"application/json": {
							Schema: &openapi3.SchemaRef{
								Ref: "#/components/schemas/CageMessageRequest",
							},
						},
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

func (ch *CagesHandler) ProcessMessage(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var messageData models.CageMessageRequest
	err := json.NewDecoder(r.Body).Decode(&messageData)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	err = ch.cagesUsecase.ProcessCageMessage(ctx, mapCageMessageToDomain(messageData))
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

	return nil
}

var GetCageMessagesOp = openapi3Struct.Path{
	Path: "/v1/cages/{id}/messages",
	Item: openapi3.PathItem{
		Get: &openapi3.Operation{
			Tags:        []string{"Cages"},
			OperationID: "getCageMessages",
			Description: "fetch cage messages with cursor",
			Parameters: openapi3.Parameters{
				{
					Value: &openapi3.Parameter{
						Name:        "id",
						In:          "path",
						Description: "cage id",
						Required:    true,
						Schema:      openapi3.NewSchemaRef("#/components/schemas/ID", nil),
					},
				},

				{
					Ref: "#/components/parameters/beforeQuery",
				},
				{
					Ref: "#/components/parameters/afterQuery",
				},
				{
					Ref: "#/components/parameters/limitQuery",
				},
			},
			Responses: map[string]*openapi3.ResponseRef{
				"200": {
					Value: &openapi3.Response{
						Description: util.ToPointer("Cage messages response"),
						Content: map[string]*openapi3.MediaType{
							"application/json": {
								Schema: openapi3.NewSchemaRef("#/components/schemas/CageMessagesCursored", nil),
							},
						},
					},
				},
			},
		},
	},
}

func (ch *CagesHandler) FetchMessages(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	offsetlimit := defaultOffsetLimit()
	cursor, err := getCursorFromRequest(r)
	if err == nil {
		offsetlimit, err = mapCursorInput(ctx, cursor)
		if err != nil {
			return err
		}
	}
	data, err := ch.cagesUsecase.GetCageMessages(ctx, parsedId, *offsetlimit)
	if err != nil {
		return err
	}
	result := []models.CageMessage{}
	for _, message := range data.Data {
		result = append(result, mapCageMessageToModels(message))
	}
	messagesCursored := models.CageMessagesCursored{
		Data:   result,
		Cursor: models.Cursor(data.Cursor),
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(messagesCursored)
	if err != nil {
		return err
	}
	return nil
}

func mapCageMessageToDomain(message models.CageMessageRequest) domain.CageMessageData {
	timestamp := time.Unix(message.Timestamp, 0)
	return domain.CageMessageData{
		Revision:  message.Revision,
		Water:     message.Water,
		Food:      message.Food,
		Light:     message.Light,
		Temp:      message.Temp,
		Humidity:  message.Humidity,
		VideoUrl:  message.VideoUrl,
		Timestamp: timestamp,
	}
}

func mapCageMessageToModels(message domain.CageMessage) models.CageMessage {
	cageId := message.CageID.String()
	return models.CageMessage{
		CageID:    cageId,
		Revision:  message.Revision,
		Water:     message.Water,
		Food:      message.Food,
		Light:     message.Light,
		Temp:      message.Temp,
		Humidity:  message.Humidity,
		VideoUrl:  message.VideoUrl,
		Timestamp: message.Timestamp.Unix(),
	}
}

func mapCageToModel(cage domain.Cage) models.Cage {
	var userID *string = nil
	if cage.UserID != nil {
		userIDStr := cage.UserID.String()
		userID = &userIDStr
	}
	return models.Cage{
		ID:          cage.ID.String(),
		UserID:      userID,
		Description: cage.Description,
		Register:    cage.Register.Unix(),
		Name:        cage.Name,
	}
}
