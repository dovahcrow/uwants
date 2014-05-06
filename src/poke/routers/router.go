package routers

import (
	"github.com/astaxie/beego"
	"poke/controllers"
)

func init() {
	beego.Router("/uwants/reply/rest", &controllers.AjaxController{}, `get:UwantsReply`)
	beego.Router("/uwants/topic/rest", &controllers.AjaxController{}, `get:UwantsTopic`)
	beego.Router("/uwants/reply", &controllers.Uwants{}, "get:GetReply;post:PostReply")
	beego.Router("/uwants/topic", &controllers.Uwants{}, "get:GetTopic;post:PostTopic")

	beego.Router("/uwants/users", &controllers.UserController{}, "get:Index")
	beego.Router(`/uwants/users/update/:id(\d+)`, &controllers.UserController{}, "post:UpdateUser")
	beego.Router("/uwants/users/create", &controllers.UserController{}, "post:CreateUser")
	beego.Router(`/uwants/users/delete/:id(\d+)`, &controllers.UserController{}, "get:DeleteUser")

	beego.Router("/uwants/threads", &controllers.ThreadController{}, "get:Index")
	beego.Router(`/uwants/threads/update/:id(\d+)`, &controllers.ThreadController{}, "post:UpdateThread")
	beego.Router("/uwants/threads/create", &controllers.ThreadController{}, "post:CreateThread")
	beego.Router(`/uwants/threads/delete/:id(\d+)`, &controllers.ThreadController{}, "get:DeleteThread")
}
