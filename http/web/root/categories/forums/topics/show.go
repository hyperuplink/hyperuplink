package topics

import (
	"github.com/gofiber/fiber/v3"
)

func (r *Route) Show(c fiber.Ctx) (err error) {
	// myRoute := route.For("CategoriesForumsTopics")
	// req := request.New(r, c, myRoute,
	// 	[]string{"base"}, "categories/forums/topics/show",
	// 	"")
	//
	// if ret, rerr := req.AccessControl(user.GuestRole); ret {
	// 	return rerr
	// }
	//
	// var cat *topic.Topic
	// cat, err = r.Runtime.Repositories.Topic.GetBySlug(
	// 	c.Params("topics"), // TODO: Abstract into req, automatic err handling
	// 	common.QueryOptions{
	// 		Limit: 1,
	// 	},
	// )
	// if ret, rerr := req.RespondOnError(err); ret == true {
	// 	return rerr
	// }
	//
	// req.UpdateTitle(cat.Name)
	// req.SetData("topic", cat)
	//
	// var fums *[]vforum.VForum
	// fums, err = r.Runtime.Repositories.Forum.VAllForTopicUUID(
	// 	cat.ID,
	// 	common.QueryOptions{
	// 		OrderBy: "position",
	// 		Order:   common.Ascending,
	// 	},
	// )
	// if ret, rerr := req.RespondOnError(err); ret == true {
	// 	return rerr
	// }
	//
	// req.SetData("forums", fums)
	//
	// return req.Respond()
	return nil
}
