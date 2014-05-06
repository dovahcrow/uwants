package controllers

import (
	//"fmt"
	"poke/base"
	"poke/models"
	"strconv"
)

type ThreadController struct {
	base.BaseController
}

func (this *ThreadController) Index() {
	us, err := models.GetAllThreads()
	if err != nil {
		this.Data[`error`] = err
		this.TplNames = `error.html`
		return
	}
	this.Data[`threads`] = us
	this.Data[`position`] = "thread"
	this.TplNames = `threads.html`
}
func (this *ThreadController) UpdateThread() {
	ids := this.Ctx.Input.Param(`:id`)
	id, _ := strconv.ParseInt(ids, 10, 0)
	title := this.GetString(`title`)
	text := this.GetString(`text`)
	err := models.UpdateThread(int(id), title, text)
	if err != nil {
		this.Data[`json`] = map[string]interface{}{`error`: err}
		this.ServeJson()
		return
	}
	this.Data[`json`] = map[string]interface{}{`succ`: `ok`}
	this.ServeJson()
	return

}
func (this *ThreadController) CreateThread() {
	title := this.GetString(`title`)
	text := this.GetString(`text`)
	err := models.CreateThread(title, text)
	if err != nil {
		this.Data[`json`] = map[string]interface{}{`error`: err}
		this.ServeJson()
		return
	}
	this.Data[`json`] = map[string]interface{}{`succ`: `ok`}
	this.ServeJson()
	return
}
func (this *ThreadController) DeleteThread() {
	ids := this.Ctx.Input.Param(`:id`)
	id, _ := strconv.ParseInt(ids, 10, 0)
	err := models.DeleteThread(int(id))
	if err != nil {
		this.Data[`json`] = map[string]interface{}{`error`: err}
		this.ServeJson()
		return
	}
	this.Data[`json`] = map[string]interface{}{`succ`: `ok`}
	this.ServeJson()
	return
}
