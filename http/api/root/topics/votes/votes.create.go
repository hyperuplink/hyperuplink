package votes

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/view"
	topicsfeed "xn--gckvb8fzb.com/hyperuplink/logic/root/topics"
)

// @Summary		Vote in a poll
// @Description	Casts the viewer's vote for the option at index selection. A vote
// @Description	is final, so a second one is a conflict, and a poll that has ended
// @Description	takes no more votes. The viewer needs write permission in the
// @Description	topic's category.
// @Tags			board
// @Accept			json
// @Produce		json
// @Param			topic	path		string			true	"The topic UUID or short id"
// @Param			request	body		votes.VoteBody	true	"The option to vote for"
// @Success		201		{object}	object{poll=view.Poll}
// @Failure		401		{object}	problem.Problem
// @Failure		403		{object}	problem.Problem
// @Failure		404		{object}	problem.Problem
// @Failure		409		{object}	problem.Problem
// @Failure		422		{object}	problem.Problem
// @Security		BearerAuth
// @Router			/topics/{topic}/votes [post]
func (r *Route) Create(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireUser(); ret {
		return rerr
	}

	top, err := topicsfeed.Read(req.Context(), r.Runtime, c.Params("topic"), req.Perms())
	if ret, rerr := req.RespondOnReadError(err); ret {
		return rerr
	}

	in := new(VoteBody)
	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	voted, err := topicsfeed.Vote(req.Context(), r.Runtime, &topicsfeed.VoteInput{
		Topic:     top,
		AuthorID:  req.User().ID,
		Selection: in.Selection,
	}, req.Perms())
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.RespondCreated(fiber.Map{
		"poll": view.PollFrom(voted.Poll, voted.Topic.EndedAt, true),
	})
}
