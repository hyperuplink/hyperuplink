package topics

import (
	"reflect"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/reply"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

type TopicCreateForm struct {
	Text    string `form:"text" validate:"required,min=1"`
	TopicID string `form:"topic_id" validate:"required,uuid"`
	ReplyID string `form:"reply_id" validate:"omitempty,uuid"`
}

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

	frm := new(TopicCreateForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute.Fill(
			map[string]string{
				"categories": c.Params("categories"),
				"forums":     c.Params("forums"),
				"topics":     c.Params("topics"),
			},
		))
	}

	r.Runtime.Debug("form", frm)

	rep := new(reply.Reply)
	rep.ShortID = shortuuid.New()
	rep.TopicID, err = uuid.Parse(frm.TopicID)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	top, err := r.Runtime.Repositories.Topic.GetByUUID(
		rep.TopicID,
		common.QueryOptions{Limit: 1},
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}
	fum, err := r.Runtime.Repositories.Forum.GetByUUID(
		top.ForumID,
		common.QueryOptions{Limit: 1},
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}
	if !req.Perms().CanWriteID(fum.CategoryID) {
		return req.RedirectToRoot()
	}

	if frm.ReplyID != "" {
		var replyUUID uuid.UUID
		replyUUID, err = uuid.Parse(frm.ReplyID)
		if ret, rerr := req.RespondOnError(err); ret == true {
			return rerr
		}
		rep.ReplyID = uuid.NullUUID{
			UUID:  replyUUID,
			Valid: true,
		}
	} else {
		rep.ReplyID = uuid.NullUUID{
			UUID:  uuid.Nil,
			Valid: false,
		}
	}
	rep.AuthorID, _ = req.Session.GetUserUUID()
	rep.Text = frm.Text
	rep.HTML, err = r.Runtime.Markdown.Convert(rep.Text)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	rep.AttachmentIDs, err = helpers.ProcessAttachments(r.Runtime, c, rep.AuthorID)
	if err != nil {
		req.Flash.SetError(err)
		return req.RedirectToRoute(myRoute.Fill(
			map[string]string{
				"categories": c.Params("categories"),
				"forums":     c.Params("forums"),
				"topics":     c.Params("topics"),
			},
		))
	}

	rep.ID, err = r.Runtime.Repositories.Reply.Create(rep)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	r.notifyReply(c, req, rep, top, fum)

	var total int64
	total, err = r.Runtime.Repositories.Reply.VAllCountForTopicUUID(rep.TopicID, common.QueryOptions{})

	pages := helpers.GetNumberOfPages(total, req.System.GetPostsPerPage())

	return req.RedirectToRouteWithQuery(myRoute.Fill(
		map[string]string{
			"categories": c.Params("categories"),
			"forums":     c.Params("forums"),
			"topics":     c.Params("topics"),
		},
	), "page", strconv.Itoa(pages))
}
