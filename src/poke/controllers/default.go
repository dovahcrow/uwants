package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"loginpool"
	"math/rand"
	"poke/base"
	"poke/models"
	//"regexp"
	"strconv"
	"time"
)

var Waittime time.Duration = 10

type rets struct {
	Username string
	Ret      interface{}
	URL      string
	Title    string
	Text     string
}
type AjaxController struct {
	base.BaseController
}

func (this *AjaxController) UwantsReply() {
	username := this.GetString(`username`)
	password := this.GetString(`password`)
	fmt.Println(``)
	cl := loginpool.GetUwants(username, password)

	tid := this.GetString(`tid`)
	title := this.GetString(`title`)
	text := this.GetString(`text`)
	addr, err := cl.SendReply(tid, title, text)
	this.Data[`addr`] = addr
	this.Data[`err`] = err
	this.TplNames = `ret.html`

}
func (this *AjaxController) UwantsTopic() {
	username := this.GetString(`username`)
	password := this.GetString(`password`)
	cl := loginpool.GetUwants(username, password)

	fid := this.GetString(`fid`)
	title := this.GetString(`title`)
	text := this.GetString(`text`)
	addr, err := cl.NewTopic(fid, title, text)
	this.Data[`addr`] = addr
	this.Data[`err`] = err
	this.TplNames = `ret.html`
}

type Uwants struct {
	base.BaseController
}

func (this *Uwants) GetReply() {
	threads, _ := models.GetAllThreads()
	users, _ := models.GetAllUsers()
	this.Data[`threads`] = threads
	this.Data[`users`] = users
	this.Data[`position`] = "reply"
	this.Data[`subp`] = `reply`
	this.TplNames = `uwants-reply.html`
}
func (this *Uwants) GetTopic() {
	threads, _ := models.GetAllThreads()
	users, _ := models.GetAllUsers()
	this.Data[`threads`] = threads
	this.Data[`users`] = users
	this.Data[`position`] = "topic"
	this.Data[`subp`] = `topic`
	this.TplNames = `uwants-topic.html`

}
func (this *Uwants) PostReply() {
	users := this.GetStrings(`users`)
	beego.Trace("POST reply user ids:", users)

	threads := this.GetStrings(`threads`)
	beego.Trace("POST reply thread ids:", threads)

	waittime, _ := this.GetInt(`waittime`)

	if waittime < 20 {
		waittime = 20
	}
	beego.Trace("set waittime", waittime)

	tids, err := models.GetAllTids()
	if err != nil {
		this.Ctx.Output.Body([]byte(`{"error":"获取可用帖子id失败"}`))
		beego.Critical(`get all tids fail!`, err)
		return
	}
	if len(tids) == 0 {
		this.Ctx.Output.Body([]byte(`{"error":"帖子id数据库为空,请先采集"}`))
		beego.Warn(`tids is empty`)
		return
	}
	go func() {

		rd := rand.New(rand.NewSource(time.Now().Unix()))

		for _, v := range threads {
			time.Sleep(time.Duration(waittime) * time.Second)

			useri, _ := strconv.Atoi(users[rd.Intn(len(users))])
			threadi, _ := strconv.Atoi(v)

			u, _ := models.GetUserById(useri)
			t, _ := models.GetThreadById(threadi)

			cl := loginpool.GetUwants(u.Username, u.Password)

			tid := tids[rd.Intn(len(tids))].Tids

			beego.Trace(`User "`, u.Username, `" send "`, t.Title, `" to tid: "`, tid, `"`)

			rt, err := cl.SendReply(tid, t.Title, t.Text)
			if err != nil {
				beego.Warn(`send reply to`, tid, `fail`)
			} else {
				models.AddThreadSends(threadi)
				beego.Info(`reply ret`, rt)
			}
		}
	}()

	this.Ctx.Output.Body([]byte(`{"ok":"任务已提交"}`))
}

func (this *Uwants) PostTopic() {
	users := this.GetStrings(`users`)
	beego.Trace("POST topic user ids:", users)

	threads := this.GetStrings(`threads`)
	beego.Trace("POST topic thread ids:", threads)
	waittime, _ := this.GetInt(`waittime`)

	if waittime < 20 {
		waittime = 20
	}
	fids, err := models.GetAllFids()

	if err != nil {
		this.Ctx.Output.Body([]byte(`{"error":"获取可用板块id失败"}`))
		beego.Critical(`get all fids fail!`, err)
		return
	}
	if len(fids) == 0 {
		this.Ctx.Output.Body([]byte(`{"error":"板块id数据库为空,请先采集"}`))
		beego.Warn(`tids is empty`)
		return
	}

	beego.Trace("set waittime", waittime)
	go func() {

		rd := rand.New(rand.NewSource(time.Now().Unix()))

		for _, v := range threads {
			time.Sleep(time.Duration(waittime) * time.Second) //sleep some times

			useri, _ := strconv.Atoi(users[rd.Intn(len(users))])
			threadi, _ := strconv.Atoi(v)

			u, _ := models.GetUserById(useri)
			t, _ := models.GetThreadById(threadi)

			cl := loginpool.GetUwants(u.Username, u.Password)

			fid := fids[rd.Intn(len(fids))].Fids

			beego.Trace("User", u.Username, "send", t.Title, `to fid:`, fid)

			//send the topic
			rt, err := cl.NewTopic(fid, t.Title, t.Text)
			if err != nil {
				models.AddRecord(u.Username, err.Error(), t.Title, false)
				beego.Warn(`post topic error`, err)
				continue
			} else {
				models.AddThreadSends(threadi)
				models.AddRecord(u.Username, rt, t.Title, true)
			}
		}
	}()

	this.Ctx.Output.Body([]byte(`{"ok":"任务已提交"}`))
}
func (this *Uwants) Records() {
	pageI, _ := this.GetInt("page")
	page := int(pageI)
	count, _ := models.GetRecordsCounts()

	pages := func() int {
		if int(count)%10 == 0 {
			return int(count) / 10
		} else {
			return int(count)/10 + 1
		}
	}()
	if pages < page {
		page = 0
	}
	re, err := models.GetRecordsRange(page*10, (page+1)*10-1)
	if err != nil {
		beego.Warn(`get page range error`, err)
	}

	this.Data[`page3`] = page
	if page-1 > 0 {
		this.Data[`page2`] = page - 1
	}
	if page-2 > 0 {
		this.Data[`page1`] = page - 2
	}
	if page+1 < pages {
		this.Data[`page4`] = page + 1
	}
	if page+2 < pages {
		this.Data[`page2`] = page + 2
	}
	succn, err := models.GetRecordsCountSuccOrNot(true)
	if err != nil {
		beego.Warn(`get succ count error`, err)
	}
	failn, err := models.GetRecordsCountSuccOrNot(false)
	if err != nil {
		beego.Warn(`get fail count error`, err)
	}
	this.Data[`succ`] = succn
	this.Data[`fail`] = failn
	this.Data[`records`] = re
	this.TplNames = `records.html`
}
