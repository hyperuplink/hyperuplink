package board

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/view"
	logicroot "xn--gckvb8fzb.com/hyperuplink/logic/root"
)

// @Summary		Show the board
// @Description	The categories the viewer may read, each with its forums, and the
// @Description	five most recently active topics among them.
// @Tags			board
// @Produce		json
// @Success		200	{object}	object{categories=[]view.Category,recent_topics=[]view.Topic}
// @Failure		401	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/board [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	board, err := logicroot.BoardView(r.Runtime, req.Perms())
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	m := req.Views()

	cats := make([]view.Category, 0, len(board.CategoriesForums))
	for i := range board.CategoriesForums {
		cf := &board.CategoriesForums[i]
		cats = append(cats, m.Category(&cf.Category, m.Forums(cf.Forums, req.Perms()), req.Perms()))
	}

	return req.Respond(fiber.Map{
		"categories":    cats,
		"recent_topics": m.Topics(board.RecentTopics),
	})
}
