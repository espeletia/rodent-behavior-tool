package handlers

import (
	"context"
	"embed"
	"net/http"
	"path"

	"tusk/internal/domain"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/runtime/middleware"
)

//go:embed swagger
var f embed.FS

const basePath = "/"

func HandleSwaggerFile() (http.Handler, []domain.ApiPath, error) {
	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}

	bytes, err := f.ReadFile("swagger/swagger.yaml")
	if err != nil {
		return nil, nil, err
	}
	newDoc, err := loader.LoadFromData(bytes)
	if err != nil {
		return nil, nil, err
	}
	unSupportedPaths := []domain.ApiPath{}

	jsonBytes, err := newDoc.MarshalJSON()
	if err != nil {
		return nil, nil, err
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "application/json")
		// #nosec: Errors unhandled
		w.Write(jsonBytes)
	}), unSupportedPaths, nil
}

func HandleSwaggerUI() http.Handler {
	return middleware.SwaggerUI(middleware.SwaggerUIOpts{
		BasePath: basePath,
		SpecURL:  path.Join(basePath, "swagger.json"),
		Path:     "/swagger",
	}, http.NotFoundHandler())
}
