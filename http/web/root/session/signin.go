package session

import (
	"errors"
	"reflect"

	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/asyncjob/notification/reply"
	"github.com/mrusme/hyperuplink/models/user"
)

type SignInForm struct {
	// TODO: https://github.com/go-playground/validator/issues/807
	Username string `form:"username" validate:"required,min=2,max=32"`
	Password string `form:"password" validate:"required,min=8,max=64"`
}

func (r *Route) SignInShow(c fiber.Ctx) (err error) {
	myRoute := route.For("SessionSignin")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL(),
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(user.GuestRole); ret {
		return rerr
	}

	//---------- [ TEST CODE ]--------------------------------------------------//
	var re []reply.Reply
	re = append(re, reply.Reply{
		Recipient: reply.Recipient{
			Username: "John",
			Address:  "john@doe.org",
			Lang:     "en",
		},
		Category: reply.Category{
			Title: "Cloud",
		},
		Forum: reply.Forum{
			Title: "General",
		},
		Topic: reply.Topic{
			Title: "Is Hetzner good?",
		},
	})
	re = append(re, reply.Reply{
		Recipient: reply.Recipient{
			Username: "Tom",
			Address:  "tom@ato.org",
			Lang:     "en",
		},
		Category: reply.Category{
			Title: "Cloud",
		},
		Forum: reply.Forum{
			Title: "General",
		},
		Topic: reply.Topic{
			Title: "Is Google Cloud still bad?",
		},
	})
	err = r.Runtime.Dispatch.ReplyNotifications(
		"notifications",
		re,
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}
	//---------- [/ TEST CODE ]-------------------------------------------------//

	return req.Respond()
}

func (r *Route) SignInCreate(c fiber.Ctx) (err error) {
	myRoute := route.For("SessionSignin")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL(),
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(user.GuestRole); ret {
		return rerr
	}

	frm := new(SignInForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.Respond()
	}

	r.Runtime.Debug("form", frm)

	var usr *user.User
	usr, err = r.Runtime.Repositories.User.GetByUsername(frm.Username)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	var match bool = false
	if match, _, err = usr.CheckPassword(frm.Password); !match {
		if err == nil {
			err = errors.New("username_password_wrong")
		}
	}
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	if err := req.Session.Set("local", usr.ID.String()); err != nil {
		req.Session.Reset()
		return req.RespondError(err)
	}

	return req.RedirectToRoot()
}
