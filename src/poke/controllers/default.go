package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"loginpool"
	"math/rand"
	"poke/base"
	"poke/models"
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
	this.TplNames = `uwants-reply.html`
}
func (this *Uwants) PostReply() {
	users := this.GetStrings(`users`)
	beego.Trace("POST reply user ids:", users)
	threads := this.GetStrings(`threads`)
	beego.Trace("POST reply thread ids:", threads)
	ret := []rets{}
	rd := rand.New(rand.NewSource(time.Now().Unix()))
	for _, v := range threads {
		time.Sleep(Waittime * time.Second)
		useri, _ := strconv.Atoi(users[rd.Intn(len(users))])
		threadi, _ := strconv.Atoi(v)
		u, _ := models.GetUserById(useri)
		t, _ := models.GetThreadById(threadi)
		cl := loginpool.GetUwants(u.Username, u.Password)
		r := rets{}
		tid := fmt.Sprint(rd.Intn(16904811))
		beego.Trace(`User "`, u.Username, `" send "`, t.Title, `" to tid: "`, tid, `"`)
		rt, err := cl.SendReply(fmt.Sprint(rd.Intn(16904811)), t.Title, t.Text)
		r.Ret = err
		r.Text = t.Text
		r.Title = t.Title
		r.URL = rt
		r.Username = u.Username
		ret = append(ret, r)

	}
	this.Data[`rets`] = ret
	this.TplNames = `ret.html`
}

func (this *Uwants) PostTopic() {
	users := this.GetStrings(`users`)
	beego.Trace("POST topic user ids:", users)
	threads := this.GetStrings(`threads`)
	beego.Trace("POST topic thread ids:", threads)
	ret := []rets{}
	rd := rand.New(rand.NewSource(time.Now().Unix()))
	for _, v := range threads {
		time.Sleep(Waittime * time.Second)
		useri, _ := strconv.Atoi(users[rd.Intn(len(users))])
		threadi, _ := strconv.Atoi(v)
		u, _ := models.GetUserById(useri)
		t, _ := models.GetThreadById(threadi)
		cl := loginpool.GetUwants(u.Username, u.Password)
		r := rets{}
		fid := fmt.Sprint(rd.Intn(400))
		beego.Trace("User", u.Username, "send", t.Title, `to fid:`, fid)
		rt, err := cl.NewTopic(fid, t.Title, t.Text)
		r.Ret = err
		r.Text = t.Text
		r.Title = t.Title
		r.URL = rt
		r.Username = u.Username
		ret = append(ret, r)

	}

	this.Data[`rets`] = ret
	this.TplNames = `ret.html`
}
func (this *Uwants) GetTopic() {
	threads, _ := models.GetAllThreads()
	users, _ := models.GetAllUsers()
	this.Data[`threads`] = threads
	this.Data[`users`] = users
	this.Data[`position`] = "topic"
	this.TplNames = `uwants-topic.html`

}
