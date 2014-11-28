package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"poke/controllers"
)

func init() {
	beego.Get("/", func(c *context.Context) { c.Redirect(302, beego.UrlFor("BasicController.Index")) })
	beego.Router("/", &controllers.BasicController{}, "get:Index")
	beego.Router("/uwants", &controllers.BasicController{}, "get:Index")
	beego.Router("/uwants/login", &controllers.LoginController{}, "get,post:Login")
	beego.Router("/uwants/reply/rest", &controllers.AjaxController{}, `get:UwantsReply`)
	beego.Router("/uwants/topic/rest", &controllers.AjaxController{}, `get:UwantsTopic`)
	beego.Router("/uwants/reply", &controllers.Uwants{}, "get:GetReply;post:PostReply")
	beego.Router("/uwants/topic", &controllers.Uwants{}, "get:GetTopic;post:PostTopic")

	beego.Router("/uwants/records", &controllers.TopicRecords{}, "get:Records")
	beego.Router(`/uwants/records/search/`, &controllers.TopicRecords{}, "get:Search")

	beego.Router("/uwants/replyrecords", &controllers.ThreadRecords{}, "get:Records")
	beego.Router(`/uwants/replyrecords/search/`, &controllers.ThreadRecords{}, "get:Search")

	beego.Router("/uwants/users/single", &controllers.UserController{}, "get:Single")
	beego.Router("/uwants/users/bunch", &controllers.UserController{}, "get:Bunch")
	beego.Router("/uwants/users/bunch", &controllers.UserController{}, "post:BunchAdd")
	beego.Router(`/uwants/users/update/:id(\d+)`, &controllers.UserController{}, "post:UpdateUser")
	beego.Router("/uwants/users/create", &controllers.UserController{}, "post:CreateUser")
	beego.Router(`/uwants/users/delete/:id(\d+)`, &controllers.UserController{}, "get:DeleteUser")

	beego.Router("/uwants/threads/single", &controllers.ThreadController{}, "get:Single")
	beego.Router("/uwants/threads/bunch", &controllers.ThreadController{}, "get:Bunch")
	beego.Router("/uwants/threads/bunch", &controllers.ThreadController{}, "post:BunchAdd")
	beego.Router(`/uwants/threads/update/:id(\d+)`, &controllers.ThreadController{}, "post:UpdateThread")
	beego.Router("/uwants/threads/create", &controllers.ThreadController{}, "post:CreateThread")
	beego.Router(`/uwants/threads/delete/:id(\d+)`, &controllers.ThreadController{}, "get:DeleteThread")

}
