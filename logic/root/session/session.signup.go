package session

import (
	"fmt"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"

	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob/confirmation/signupconfirmation"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
	repoSetting "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

type SignUpInput struct {
	// TODO: https://github.com/go-playground/validator/issues/807
	Username       string `json:"username" form:"username" validate:"required,min=2,max=32"`
	Email          string `json:"email" form:"email" validate:"required,max=254"`
	EmailIsJID     bool   `json:"email_is_jid" form:"email_is_jid"`
	Password       string `json:"password" form:"password" validate:"required,min=8,max=64"`
	PasswordRepeat string `json:"password_repeat" form:"password_repeat" validate:"omitempty,eqcsfield=Password"`
	Language       string `json:"language" form:"-"`
}

type SignUpOptions struct {
	Activate bool
}

type FieldError struct {
	Field string
	Err   error
}

func (e *FieldError) Error() string {
	return e.Err.Error()
}

func (e *FieldError) Unwrap() error {
	return e.Err
}

func isValidJID(jid string) bool {
	at := strings.IndexByte(jid, '@')
	if at <= 0 || at == len(jid)-1 {
		return false
	}
	if strings.Count(jid, "@") != 1 {
		return false
	}
	if strings.ContainsAny(jid, " \t\r\n") {
		return false
	}
	return true
}

func SignUp(
	rt *runtime.Runtime,
	in *SignUpInput,
	opts SignUpOptions,
) (*user.User, error) {
	if err := validator.New().Struct(in); err != nil {
		return nil, err
	}

	settingAuth, err := repoSetting.GetByID[setting.Auth](rt.Repositories.Setting, "auth")
	if err != nil {
		return nil, err
	}

	var isJID bool
	switch settingAuth.JSONValue.AddressType {
	case setting.JID:
		isJID = true
	case setting.EmailAndJID:
		isJID = in.EmailIsJID
	default:
		isJID = false
	}

	if isJID {
		if !isValidJID(in.Email) {
			return nil, &FieldError{
				Field: "email",
				Err:   fmt.Errorf("%w_email_jid", errs.ErrValidation),
			}
		}
	} else if verr := validator.New().Var(in.Email, "email"); verr != nil {
		return nil, &FieldError{
			Field: "email",
			Err:   fmt.Errorf("%w_email_email", errs.ErrValidation),
		}
	}

	usr := new(user.User)
	usr.Username = in.Username
	usr.Role = user.UserRole
	usr.Email = in.Email
	if err = usr.SetEmailForConfirmation(in.Email); err != nil {
		return nil, err
	}
	if isJID {
		usr.SetEmailIsJID()
	}
	usr.Language = in.Language
	if usr.Language == "" {
		usr.Language = "en"
	}

	promoteAdmin := rt.Config.UsersPromoteAdmin()
	if slices.Index(promoteAdmin, strings.ToLower(usr.Email)) > -1 {
		usr.Role = user.AdminRole
	}

	if err = usr.SetPassword(in.Password); err != nil {
		return nil, err
	}

	usr.ID, err = rt.Repositories.User.Create(usr)
	if err != nil {
		return nil, err
	}

	if !opts.Activate {
		return usr, nil
	}

	activated, err := rt.Repositories.User.GetByUUID(usr.ID, common.QueryOptions{Limit: 1})
	if err != nil {
		return nil, err
	}
	if err = activated.ConfirmEmail(activated.EmailConfirmationToken); err != nil {
		return nil, err
	}
	if err = rt.Repositories.User.Update(activated); err != nil {
		return nil, err
	}

	return activated, nil
}

func SendSignupConfirmation(
	rt *runtime.Runtime,
	usr *user.User,
	subject string,
	signupURL string,
) error {
	settingSystem, err := repoSetting.GetByID[setting.System](rt.Repositories.Setting, "system")
	if err != nil {
		return err
	}

	settingCommsEmail, err := repoSetting.GetByID[setting.CommsEmail](rt.Repositories.Setting, "comms_email")
	if err != nil {
		return err
	}

	settingCommsXMPP, err := repoSetting.GetByID[setting.CommsXMPP](rt.Repositories.Setting, "comms_xmpp")
	if err != nil {
		return err
	}

	sc, err := signupconfirmation.New(
		&settingSystem.JSONValue,
		usr,
		subject,
		usr.EmailConfirmationToken,
		signupURL,
	)
	if err != nil {
		return err
	}

	targetID := settingCommsEmail.JSONValue.TargetID
	if usr.EmailIsJID {
		targetID = settingCommsXMPP.JSONValue.TargetID
	}

	return rt.Dispatch.SignupConfirmation(targetID, sc)
}
