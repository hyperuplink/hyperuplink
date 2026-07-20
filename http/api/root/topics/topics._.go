package topics

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
)

type Route struct {
	route.RouteController
}

type ReplyCreateBody struct {
	Text          string   `json:"text" validate:"required,min=1"`
	ReplyID       string   `json:"reply_id" validate:"omitempty,uuid"`
	AttachmentIDs []string `json:"attachment_ids" validate:"omitempty,dive,uuid"`
}

func New(
	rt *runtime.Runtime,
	router fiber.Router,
) (*Route, error) {
	r := new(Route)

	r.Runtime = rt
	r.Router = router
	r.Path = route.For("Topics").Pathname()
	r.Env = route.NewEnv()

	r.Router.Route("/"+r.Path, func(base fiber.Router) {
		base.Get("", r.Index).Name("index")
		base.Get("/:id", r.Show).Name("show")
		base.Post("/:id/replies", r.CreateReply).Name("replies.create")
	}, r.Path+".")

	return r, nil
}

func (r *Route) GetRuntime() *runtime.Runtime {
	return r.Runtime
}

func (r *Route) GetPath() string {
	return r.Path
}

func (r *Route) GetEnv() *route.Environment {
	return r.Env
}
