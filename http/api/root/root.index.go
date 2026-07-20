package root

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicroot "xn--gckvb8fzb.com/hyperuplink/logic/root"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary		Show the board
// @Description	Lists every category the caller may see together with the forums
// @Description	inside it and the board-wide counters.
// @Tags			board
// @Produce		json
// @Success		200	{object}	logicroot.Board
// @Failure		403	{object}	request.ErrorResponse
// @Failure		500	{object}	request.ErrorResponse
// @Security		BearerAuth
// @Security		APIKeyAuth
// @Router			/ [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	var board *logicroot.Board
	board, err = logicroot.BoardView(r.Runtime, req.Perms())
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.Respond(board)
}
