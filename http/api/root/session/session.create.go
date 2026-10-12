package session

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/api/problem"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/view"
	logicaccount "xn--gckvb8fzb.com/hyperuplink/logic/root/account"
	logicapi "xn--gckvb8fzb.com/hyperuplink/logic/root/account/api"
	logicsession "xn--gckvb8fzb.com/hyperuplink/logic/root/session"
)

const keyNameLayout = "2006-01-02 15:04"

// @Summary		Sign in
// @Description	Exchanges a username and a password for an API key. An account
// @Description	with two-factor authentication has to send the current code in
// @Description	otp_code; without it the problem code is err_otp_required, so a
// @Description	client can ask for the code in a second step. The key is named
// @Description	after name, or after the time of the sign-in, and its secret is
// @Description	returned once, in token. Accounts that signed up through an
// @Description	OAuth provider have no password and obtain their keys on the
// @Description	board instead.
// @Tags			session
// @Accept			json
// @Produce		json
// @Param			request	body		session.SignInBody	true	"The credentials and an optional key name"
// @Success		201		{object}	object{token=string,key=view.Key,user=view.Account}
// @Failure		401		{object}	problem.Problem
// @Failure		422		{object}	problem.Problem
// @Failure		429		{object}	problem.Problem
// @Router			/session [post]
func (r *Route) Create(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	in := new(SignInBody)
	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	usr, err := logicsession.SignIn(req.Context(), r.Runtime, &in.SignInInput)
	if errors.Is(err, errs.ErrOTPCodeWrong) {
		return req.RespondProblem(problem.WithStatus(fiber.StatusUnauthorized, err))
	}
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	name := in.Name
	if name == "" {
		name = "API " + time.Now().UTC().Format(keyNameLayout)
	}

	key, secret, err := logicapi.Create(req.Context(), r.Runtime, usr.ID, &logicapi.CreateInput{KeyName: name})
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	acc, err := logicaccount.Show(req.Context(), r.Runtime, usr)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.RespondCreated(fiber.Map{
		"token": secret,
		"key":   view.KeyFrom(key),
		"user":  req.Views().Account(usr, acc.Groups, acc.Profile),
	})
}
