package topics

import (
	"errors"
	"reflect"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logictopics "xn--gckvb8fzb.com/hyperuplink/logic/root/categories/forums/topics"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

type TopicPollVoteForm struct {
	Selection string `form:"selection" validate:"required,number"`
}

func (r *Route) PollVote(c fiber.Ctx) (err error) {
	myRoute := route.For("CategoriesForumsTopics")
	req := request.New(r, c, myRoute,
		[]string{"base"}, "categories/forums/topics/show",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	topicRoute := myRoute.Fill(map[string]string{
		"categories": c.Params("categories"),
		"forums":     c.Params("forums"),
		"topics":     c.Params("topics"),
	})

	frm := new(TopicPollVoteForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(topicRoute)
	}

	selection, err := strconv.Atoi(frm.Selection)
	if ret, rerr := req.RedirectToRouteOnError(err, topicRoute); ret == true {
		return rerr
	}

	authorID, ok := req.Session.GetUserUUID()
	if !ok {
		return req.RedirectToRoot()
	}

	err = logictopics.PollVoteBySlugs(r.Runtime, req.Perms(), &logictopics.PollVoteBySlugsInput{
		ForumSlug: c.Params("forums"),
		TopicSlug: c.Params("topics"),
		AuthorID:  authorID,
		Selection: selection,
	})
	if errors.Is(err, errs.ErrForbidden) {
		return req.RedirectToRoot()
	}
	if errors.Is(err, errs.ErrUniqueViolationOn) {
		req.Flash.SetInfo("poll_already")
		return req.RedirectToRoute(topicRoute)
	}
	if ret, rerr := req.RedirectToRouteOnError(err, topicRoute); ret == true {
		return rerr
	}

	return req.RedirectToRoute(topicRoute)
}
