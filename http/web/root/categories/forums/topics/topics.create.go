package topics

import (
	"errors"
	"reflect"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logictopics "xn--gckvb8fzb.com/hyperuplink/logic/root/categories/forums/topics"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Create(c fiber.Ctx) (err error) {
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

	topicRoute := myRoute.Fill(
		map[string]string{
			"categories": c.Params("categories"),
			"forums":     c.Params("forums"),
			"topics":     c.Params("topics"),
		},
	)

	in := new(logictopics.CreateReplyInput)

	if ok := req.ValidateForm(in, reflect.TypeOf(*in)); !ok {
		return req.RedirectToRoute(topicRoute)
	}

	r.Runtime.Debug("form", in)

	authorID, _ := req.Session.GetUserUUID()

	var created *logictopics.CreatedReply
	created, err = logictopics.CreateReply(
		r.Runtime,
		authorID,
		req.Perms(),
		in,
		func(authorID uuid.UUID) ([]uuid.UUID, error) {
			return helpers.ProcessAttachments(r.Runtime, c, authorID)
		},
	)
	if errors.Is(err, errs.ErrForbidden) {
		return req.RedirectToRoot()
	}
	if errors.Is(err, errs.ErrAttachmentTooLarge) ||
		errors.Is(err, errs.ErrAttachmentFormatNotAllowed) ||
		errors.Is(err, errs.ErrAttachmentHookFailed) ||
		errors.Is(err, errs.ErrAttachmentDuplicate) ||
		errors.Is(err, errs.ErrAttachmentUploadFailed) {
		req.Flash.SetError(err)
		return req.RedirectToRoute(topicRoute)
	}
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	r.notifyReply(c, req, created.Reply, created.Topic, created.Forum)

	var pages int
	pages, err = logictopics.ReplyPages(
		r.Runtime,
		created.Reply.TopicID,
		req.System.GetPostsPerPage(),
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRouteWithQuery(topicRoute, "page", strconv.Itoa(pages))
}
