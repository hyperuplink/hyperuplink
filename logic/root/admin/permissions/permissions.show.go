package permissions

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
)

func Show(rt *runtime.Runtime) (view *View, err error) {
	cats, err := gh.Repositories(rt).Category.All(common.QueryOptions{
		OrderBy: "position",
		Order:   common.Ascending,
	})
	if err != nil {
		return nil, err
	}

	groups, err := gh.Repositories(rt).Group.All(common.QueryOptions{
		OrderBy: "name",
		Order:   common.Ascending,
	})
	if err != nil {
		return nil, err
	}

	perms, err := gh.Repositories(rt).Permission.All()
	if err != nil {
		return nil, err
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

	return &View{
		DefaultLevel:  defaultLevel,
		Groups:        groupViews,
		HasCategories: len(*cats) > 0,
	}, nil
}
