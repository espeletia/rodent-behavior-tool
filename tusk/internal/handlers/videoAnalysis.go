package handlers

import (
	"encoding/json"
	"ghiaccio/models"
	"net/http"
	"tusk/internal/domain"
	"tusk/internal/usecases"
	"tusk/internal/util"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/nextap-solutions/openapi3Struct"

	"github.com/getkin/kin-openapi/openapi3"
)

type VideoAnalysisHandler struct {
	videoUsecase *usecases.VideoUsecase
}

func NewVideoAnalysisHandler(videoUsecase *usecases.VideoUsecase) *VideoAnalysisHandler {
	return &VideoAnalysisHandler{
		videoUsecase: videoUsecase,
	}
}

var CreateVideoOp = openapi3Struct.Path{
	Path: "/v1/video",
	Item: openapi3.PathItem{
		Put: &openapi3.Operation{
			Tags:        []string{"VideoAnalysis"},
			OperationID: "createVideoAnalysis",
			Description: "create a new video for analysis",
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Description: "Add video data",
					Required:    true,
					Content: map[string]*openapi3.MediaType{
						"application/json": {
							Schema: &openapi3.SchemaRef{
								Ref: "#/components/schemas/CreateVideoDto",
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

func (vah *VideoAnalysisHandler) CreateVideoAnalysis(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	// Decode the JSON body to the `Viewport` struct
	var videoData models.CreateVideoDto
	err := json.NewDecoder(r.Body).Decode(&videoData)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	err = vah.videoUsecase.CreateNewVideo(ctx, mapVideoDTOToDomain(videoData))
	if err != nil {
		return err
	}

	// Convert the result to JSON and write to the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

var GetVideoByIDOp = openapi3Struct.Path{
	Path: "/v1/video/{id}",
	Item: openapi3.PathItem{
		Get: &openapi3.Operation{
			Tags:        []string{"VideoAnalysis"},
			OperationID: "getVideoAnalysis",
			Description: "get video analysis by id",
			Parameters: openapi3.Parameters{
				{
					Value: &openapi3.Parameter{
						Name:        "id",
						In:          "path",
						Description: "video id",
						Required:    true,
						Schema:      openapi3.NewSchemaRef("#/components/schemas/ID", nil),
					},
				},
			},
			Responses: map[string]*openapi3.ResponseRef{
				"200": {
					Value: &openapi3.Response{
						Description: util.ToPointer("VideoAnalysis"),
						Content: map[string]*openapi3.MediaType{
							"application/json": {
								Schema: openapi3.NewSchemaRef("#/components/schemas/VideoAnalysis", nil),
							},
						},
					},
				},
			},
		},
	},
}

func (vah *VideoAnalysisHandler) GetVideoAnalysisByID(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	videoId := mux.Vars(r)["id"]
	videoUuid, err := uuid.Parse(videoId)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	video, err := vah.videoUsecase.GetByID(ctx, videoUuid)
	if err != nil {
		return err
	}

	// Convert the result to JSON and write to the response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(mapVideoToModel(*video))
	if err != nil {
		return err
	}
	return nil
}

var GetVideosCursoredOp = openapi3Struct.Path{
	Path: "/v1/videos",
	Item: openapi3.PathItem{
		Get: &openapi3.Operation{
			Tags:        []string{"VideoAnalysis"},
			OperationID: "getVideosCursored",
			Description: "fetch video analysis with cursor",
			Parameters: openapi3.Parameters{
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
						Description: util.ToPointer("videos cursored response"),
						Content: map[string]*openapi3.MediaType{
							"application/json": {
								Schema: openapi3.NewSchemaRef("#/components/schemas/CursoredVideoAnalysis", nil),
							},
						},
					},
				},
			},
		},
	},
}

func (vah *VideoAnalysisHandler) GetVideosCursored(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	offsetlimit := defaultOffsetLimit()
	cursor, err := getCursorFromRequest(r)
	if err == nil {
		offsetlimit, err = mapCursorInput(ctx, cursor)
		if err != nil {
			return err
		}
	}
	data, err := vah.videoUsecase.GetVideosCursored(ctx, *offsetlimit)
	if err != nil {
		return err
	}
	result := []models.VideoAnalysis{}
	for _, video := range data.Data {
		result = append(result, mapVideoToModel(video))
	}
	videosCursored := models.CursoredVideoAnalysis{
		Data:   result,
		Cursor: models.Cursor(data.Cursor),
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(videosCursored)
	if err != nil {
		return err
	}

	return nil
}

func mapVideoDTOToDomain(data models.CreateVideoDto) domain.CreateVideoDto {
	return domain.CreateVideoDto{
		VideoUrl:    data.VideoUrl,
		Description: data.Description,
		Name:        data.Name,
	}
}

func mapVideoToModel(v domain.Video) models.VideoAnalysis {
	var analysedVideo *models.MediaFile = nil
	if v.AnalysedVideo != nil {
		video := mapMediaFileToModel(*v.AnalysedVideo)
		analysedVideo = &video
	}
	return models.VideoAnalysis{
		ID:            v.ID.String(),
		Video:         mapMediaFileToModel(v.Video),
		OwnerId:       v.OwnerId.String(),
		Description:   v.Description,
		Name:          v.Name,
		AnalysedVideo: analysedVideo,
	}
}

func mapMediaFileToModel(m domain.MediaFile) models.MediaFile {
	var masterId *string = nil
	if m.MasterID != nil {
		stringId := m.MasterID.String()
		masterId = &stringId
	}
	return models.MediaFile{
		ID:         m.ID.String(),
		MimeType:   m.MimeType,
		Variant:    m.Variant,
		EntityType: m.EntityType,
		MasterID:   masterId,
		Url:        m.Url,
		Duration:   m.Duration,
		Size:       m.Size,
		Width:      m.Width,
		Height:     m.Height,
		Master:     m.Master,
	}
}
