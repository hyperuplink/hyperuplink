package topics

import (
	"reflect"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/helpers"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/reply"
	"github.com/mrusme/hyperuplink/models/user"
	"github.com/mrusme/hyperuplink/services/repositories/common"
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
		user.GuestRole, // TODO: Remove!
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
	rep.TopicID, err = uuid.Parse(frm.TopicID)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
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
	rep.Kind = reply.Regular
	rep.Text = frm.Text
	// rep.PollVote =
	// rep.RSVPResponse =
	_, err = r.Runtime.Repositories.Reply.Create(rep)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

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
