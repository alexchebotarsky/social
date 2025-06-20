package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alexchebotarsky/social/social-media-aggregator/openapi"
)

func OpenAPIYaml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(openapi.OpenAPIYaml)))
	_, err := w.Write(openapi.OpenAPIYaml)
	if err != nil {
		log.Printf("Error writing OpenAPI YAML: %v", err)
	}
}

func SwaggerUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(openapi.SwaggerUIHtml)))
	_, err := w.Write(openapi.SwaggerUIHtml)
	if err != nil {
		log.Printf("Error writing Swagger UI: %v", err)
	}
}
