package controllers

import (
	//"github.com/astaxie/beego"
	"poke/base"
)

type BasicController struct {
	base.BaseController
}

func (this *BasicController) Index() {
	this.Data[`version`] = `版本 v1.2.0`
	this.TplNames = "index.html"
}

type LoginController struct {
	base.BaseController
}

func (this *LoginController) Prepare() {
	this.Layout = `layout.html`
	this.Data[`position`] = ""
	this.Data[`subp`] = ""
}
func (this *LoginController) Exit() {
	this.DelSession(`login`)
	this.Redirect(`/uwants`, 302)
}
func (this *LoginController) Login() {
	switch this.Ctx.Input.Method() {
	case "POST":
		{
			if this.GetString(`password`) == `uwants` {
				this.SetSession(`login`, true)
			}
			this.Redirect(`/uwants`, 302)
		}
	case "GET":
		{
			this.TplNames = "login.html"
		}
	}

}
