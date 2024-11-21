package main

import (
	"context"
	"log"

	"tusk/internal/handlers"

	"github.com/nextap-solutions/openapi3Struct"

	"github.com/getkin/kin-openapi/openapi3"
)

func main() {
	t := openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:       "TUSK",
			Description: `Documentation of ratt API.`,
			Version:     "0.0.1",
		},
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{
				"ID": openapi3.NewSchemaRef("", &openapi3.Schema{
					Type: "string",
				}),
			},
			Parameters: map[string]*openapi3.ParameterRef{
				"beforeQuery": {
					Value: &openapi3.Parameter{
						Name:        "before",
						In:          "query",
						Description: "Cursor before",
						Schema: openapi3.NewSchemaRef("", &openapi3.Schema{
							Type: "string",
						}),
					},
				},
				"afterQuery": {
					Value: &openapi3.Parameter{
						Name:        "after",
						In:          "query",
						Description: "Cursor after",
						Schema: openapi3.NewSchemaRef("", &openapi3.Schema{
							Type: "string",
						}),
					},
				},
				"limitQuery": {
					Value: &openapi3.Parameter{
						Name:        "limit",
						In:          "query",
						Description: "Query limit",
						Schema: openapi3.NewSchemaRef("", &openapi3.Schema{
							Type: "string",
						}),
					},
				},
				"shareToken": {
					Value: &openapi3.Parameter{
						Name:        "share_token",
						In:          "query",
						Description: "A share token of the share url",
						Required:    false,
						Schema: openapi3.NewSchemaRef("", &openapi3.Schema{
							Type: "string",
						}),
					},
				},
			},
			SecuritySchemes: openapi3.SecuritySchemes{
				"bearerAuth": {
					Value: openapi3.NewJWTSecurityScheme(),
				},
			},
		},
		Security: []openapi3.SecurityRequirement{
			{
				"bearerAuth": {},
			},
		},
	}

	p := openapi3Struct.NewParser(t, openapi3Struct.WithPackagePaths([]string{"../internal/handlers/models/"}))
	err := p.ParseSchemasFromStructs()
	if err != nil {
		log.Fatalf("ParseSchemasFromStructs %v", err)
	}

	p.AddPath(handlers.PingOp)
	p.AddPath(handlers.MeOp)
	p.AddPath(handlers.LoginOp)
	p.AddPath(handlers.CreateUserOp)

	err = p.Validate(context.Background())
	if err != nil {
		log.Fatalf("Error validating %v", err)
	}

	err = p.SaveYamlToFile("../internal/handlers/swagger/swagger.yaml")
	if err != nil {
		log.Fatal(err)
	}
}