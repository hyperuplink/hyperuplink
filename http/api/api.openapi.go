package api

import (
	"io/fs"

	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/swaggo/swag"
	"xn--gckvb8fzb.com/glides/http/route"
)

const swaggerSpecPath = "docs/swagger.json"

func (srv *API) loadOpenAPI() {
	embedded, ok := srv.rt.GetEmbedOk("docs")
	if !ok {
		srv.rt.Warn("openapi", "no embedded documentation")
		return
	}

	doc, err := fs.ReadFile(embedded, swaggerSpecPath)
	if err != nil {
		srv.rt.Warn("openapi", "error", "error", err)
		return
	}

	if swag.GetSwagger(swag.Name) == nil {
		version, _, _, _ := srv.rt.GetBuild()
		swag.Register(swag.Name, &swag.Spec{
			Version:         version,
			SwaggerTemplate: string(doc),
		})
	}

	docPath := "/" + route.For("APIV1OpenAPI").AsURL()
	srv.app.Get(docPath, func(c fiber.Ctx) error {
		spec, rerr := swag.ReadDoc()
		if rerr != nil {
			return rerr
		}

		return c.Type("json").SendString(spec)
	})

	uiPath := "/" + route.For("APIV1Docs").AsURL()
	srv.app.Get(uiPath+"/*", swaggo.New(swaggo.Config{
		URL:   docPath,
		Title: "Hyperuplink API",
	}))
}
