package permissions

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

type PermissionRow struct {
	CategoryID   string
	CategoryName string
	Level        string
}

type GroupView struct {
	ID          string
	Name        string
	Rows        []PermissionRow
	AddableCats []category.Category
}

func (r *Route) Index(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminPermissions")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	cats, err := r.Runtime.Repositories.Category.All(common.QueryOptions{
		OrderBy: "position",
		Order:   common.Ascending,
	})
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	groups, err := r.Runtime.Repositories.Group.All(common.QueryOptions{
		OrderBy: "name",
		Order:   common.Ascending,
	})
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	perms, err := r.Runtime.Repositories.Permission.All()
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	defaultLevel := permission.LevelToString(permission.ReadWrite)
	groupPerms := map[string]map[uuid.UUID]byte{}
	for i := range *perms {
		p := (*perms)[i]
		if p.GroupID.Valid == false {
			if p.CategoryID.Valid == false {
				defaultLevel = p.LevelString()
			}
			continue
		}
		if p.CategoryID.Valid == false {
			continue
		}
		if groupPerms[p.GroupID.String] == nil {
			groupPerms[p.GroupID.String] = map[uuid.UUID]byte{}
		}
		groupPerms[p.GroupID.String][uuid.UUID(p.CategoryID.Bytes)] = p.Level()
	}

	groupViews := []GroupView{}
	for _, g := range *groups {
		gv := GroupView{ID: g.ID, Name: g.Name}
		mapped := groupPerms[g.ID]
		for _, cat := range *cats {
			if lvl, ok := mapped[cat.ID]; ok {
				gv.Rows = append(gv.Rows, PermissionRow{
					CategoryID:   cat.ID.String(),
					CategoryName: cat.Name,
					Level:        permission.LevelToString(lvl),
				})
			} else {
				gv.AddableCats = append(gv.AddableCats, cat)
			}
		}
		groupViews = append(groupViews, gv)
	}

	req.SetData("default_level", defaultLevel)
	req.SetData("groups", groupViews)
	req.SetData("has_categories", len(*cats) > 0)

	return req.Respond()
}
