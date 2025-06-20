package openapi

import _ "embed"

//go:embed openapi.yaml
var OpenAPIYaml []byte

//go:embed swaggerui.html
var SwaggerUIHtml []byte
