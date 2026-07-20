package topics

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	logictopics "xn--gckvb8fzb.com/hyperuplink/logic/root/categories/forums/topics"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
)

type Route struct {
	route.RouteController
}

type TopicCreateBody struct {
	logictopics.CreateReplyInput
	AttachmentIDs []string `json:"attachment_ids" validate:"omitempty,dive,uuid"`
}

type TopicPollVoteBody struct {
	Selection int `json:"selection" form:"selection"`
}

func New(
	rt *runtime.Runtime,
	router fiber.Router,
) (*Route, error) {
	r := new(Route)

	r.Runtime = rt
	r.Router = router
	r.Path = route.For("CategoriesForumsTopics").Pathname()
	r.Env = route.NewEnv()

	r.Router.Route("/"+r.Path, func(base fiber.Router) {
		base.Get("", r.Show).Name("show")
		base.Post("", r.Create).Name("create")
		base.Post("/poll", r.PollVote).Name("poll")
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
