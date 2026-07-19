package topics

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logictopics "xn--gckvb8fzb.com/hyperuplink/logic/root/categories/forums/topics"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) PollVote(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	in := new(TopicPollVoteBody)

	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	authorID, ok := req.UserUUID()
	if !ok {
		return req.RespondError(errs.ErrForbidden)
	}

	err = logictopics.PollVoteBySlugs(r.Runtime, req.Perms(), &logictopics.PollVoteBySlugsInput{
		ForumSlug: c.Params("forums"),
		TopicSlug: c.Params("topics"),
		AuthorID:  authorID,
		Selection: in.Selection,
	})
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondOK()
}
