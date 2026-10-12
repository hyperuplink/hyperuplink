package topics

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/topics/replies"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/topics/votes"
	logicnewpost "xn--gckvb8fzb.com/hyperuplink/logic/root/newpost"
)

type Route struct {
	route.RouteController
}

type TopicCreateBody struct {
	logicnewpost.CreateInput
	AttachmentIDs []string `json:"attachment_ids" validate:"omitempty,dive,uuid"`
}

func New(
	rt *runtime.Runtime,
	router fiber.Router,
) (*Route, error) {
	r := new(Route)

	r.Runtime = rt
	r.Router = router
	r.Path = route.For("APIV1Topics").Pathname()
	r.Env = route.NewEnv()

	r.Router.Route("/"+r.Path, func(base fiber.Router) {
		base.Get("", r.Index).Name("index")
		base.Post("", r.Create).Name("create")

		base.Route("/"+route.For("APIV1Topic").Pathname(), func(member fiber.Router) {
			member.Get("", r.Show).Name("show")

			repliesRoute, err := replies.New(r.Runtime, member)
			r.Runtime.NilOrDie(err)
			r.Routes = append(r.Routes, repliesRoute)

			votesRoute, err := votes.New(r.Runtime, member)
			r.Runtime.NilOrDie(err)
			r.Routes = append(r.Routes, votesRoute)
		}, "topic.")
	}, r.Path+".")

	return r, nil
}
