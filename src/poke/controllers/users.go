package controllers

import (
	//"fmt"
	"poke/base"
	"poke/models"
	"strconv"
)

type UserController struct {
	base.BaseController
}

func (this *UserController) Index() {
	us, err := models.GetAllUsers()
	if err != nil {
		this.Data[`error`] = err
		this.TplNames = `error.html`
		return
	}
	this.Data[`users`] = us
	this.Data[`position`] = "user"
	this.TplNames = `users.html`
}
func (this *UserController) UpdateUser() {
	ids := this.Ctx.Input.Param(`:id`)
	id, _ := strconv.ParseInt(ids, 10, 0)
	usnm := this.GetString(`username`)
	pswd := this.GetString(`password`)
	err := models.UpdateUser(int(id), usnm, pswd)
	if err != nil {
		this.Data[`json`] = map[string]interface{}{`error`: err}
		this.ServeJson()
		return
	}
	this.Data[`json`] = map[string]interface{}{`succ`: `ok`}
	this.ServeJson()
	return

}
func (this *UserController) CreateUser() {
	usnm := this.GetString(`username`)
	pswd := this.GetString(`password`)
	err := models.CreateUser(usnm, pswd)
	if err != nil {
		this.Data[`json`] = map[string]interface{}{`error`: err}
		this.ServeJson()
		return
	}
	this.Data[`json`] = map[string]interface{}{`succ`: `ok`}
	this.ServeJson()
	return
}
func (this *UserController) DeleteUser() {
	ids := this.Ctx.Input.Param(`:id`)
	id, _ := strconv.ParseInt(ids, 10, 0)
	err := models.DeleteUser(int(id))
	if err != nil {
		this.Data[`json`] = map[string]interface{}{`error`: err}
		this.ServeJson()
		return
	}
	this.Data[`json`] = map[string]interface{}{`succ`: `ok`}
	this.ServeJson()
	return
}
