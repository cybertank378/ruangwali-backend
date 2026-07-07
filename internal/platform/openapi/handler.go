// =========================================================
// File: internal/platform/openapi/handler.go
// =========================================================

package openapi

import (
	"bytes"
	"net/http"

	openapispec "github.com/ruangwali/cmd/api/openapi"
)

const swaggerUIHTML = `<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8" />

    <meta
        name="viewport"
        content="width=device-width, initial-scale=1"
    />

    <title>RuangWali API Documentation</title>

    <link
        rel="stylesheet"
        href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css"
    />
</head>

<body>
    <div id="swagger-ui"></div>

    <script
        src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js">
    </script>

    <script
        src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js">
    </script>

    <script>
        window.onload = function () {
            window.ui = SwaggerUIBundle({
                url: "/openapi.yaml",

                dom_id: "#swagger-ui",

                deepLinking: true,

                displayRequestDuration: true,

                filter: true,

                persistAuthorization: true,

                tryItOutEnabled: true,

                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],

                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Specification(
	writer http.ResponseWriter,
	_ *http.Request,
) {
	writer.Header().Set(
		"Content-Type",
		"application/yaml; charset=utf-8",
	)

	writer.Header().Set(
		"Cache-Control",
		"no-store",
	)

	writer.WriteHeader(
		http.StatusOK,
	)

	_, _ = bytes.NewReader(
		openapispec.Specification,
	).WriteTo(
		writer,
	)
}

func (h *Handler) Documentation(
	writer http.ResponseWriter,
	_ *http.Request,
) {
	writer.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)

	writer.Header().Set(
		"Cache-Control",
		"no-store",
	)

	writer.WriteHeader(
		http.StatusOK,
	)

	_, _ = writer.Write(
		[]byte(
			swaggerUIHTML,
		),
	)
}
