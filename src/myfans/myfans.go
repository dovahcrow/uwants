package myfans

import (
	"client"
	//"code.google.com/p/go.net/html"
	//"code.google.com/p/mahonia"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"math/rand"
	"net/url"
	"time"
	//"strings"
	"regexp"
)

var root = `http://myfans.info/`

type Myfans struct {
	*client.Client
	username string
	password string
}

func New(username, password string) *Myfans {
	u := &Myfans{}
	u.Client = client.New()
	u.UseProxy(`http://127.0.0.1:8087`)
	u.UseEncoder(`utf8`)
	u.username = username
	u.password = password
	return u
}

func (this *Myfans) Login() (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			return
		}
	}()
	re, err := this.Get(`http://myfans.info/logging.php?action=login`)
	e(`获取登陆地址失败`, err)
	defer re.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(re)
	e(`解析首页失败`, err)

	loginaddr, ok := doc.Find(`form#loginform`).Attr(`action`)
	e(`解析登陆地址失败`, ok)
	v := url.Values{}
	v.Add(`username`, this.username)
	v.Add(`password`, this.password)
	re, err = this.PostForm(root+loginaddr, v)
	e(`登陆失败`, err)
	doc, err = goquery.NewDocumentFromResponse(re)
	return nil
}

func (this *Myfans) SendReply(tid string, text string) (returl string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			return
		}
	}()
	re, err := this.Get(`http://myfans.info/viewthread.php?tid=` + tid + `&extra=page%3D1`)
	e(`获取帖子首页失败`, err)
	defer re.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(re)
	e(`创建页面索引失败`, err)

	hashvalue, ok := doc.Find(`#fastpostform`).Find(`input[name="formhash"]`).Attr(`value`)
	e(`获取回复哈希值失败,formhash`, ok)

	dhash, ok := doc.Find(`#fastpostform`).Find(`input[name="dhash"]`).Attr(`value`)
	e(`获取回复哈希值失败,dhash`, ok)

	html, _ := doc.Html()
	adhash := regexp.MustCompile(`document\.getElementById\("dhash"\)\.value="(.+)";`).FindStringSubmatch(html)
	if len(adhash) > 1 {
		dhash = adhash[1]
	}

	actionvalue, ok := doc.Find(`#fastpostform`).Attr(`action`)
	e(`获取回复地址失败`, ok)

	retv := url.Values{}
	//retv.Add(`subject`, title)
	retv.Add(`message`, text)
	retv.Add(`formhash`, hashvalue)
	retv.Add(`dhash`, dhash)

	re, err = this.PostForm(root+actionvalue, retv)
	e(`回复失败`, err)
	defer re.Body.Close()

	return ``, nil
	//target, ok := doc.Find(`meta[http-equiv="refresh"]`).Attr(`content`)
	//e(`获取回复地址失败`, ok)
	//targeturl := regexp.MustCompile(`url=(.+)`).FindStringSubmatch(target)
	//if len(targeturl) < 2 {
	//	panic(fmt.Errorf(`获取回复地址失败`))
	//}
	//return root + targeturl[1], nil

}

func (this *Myfans) NewTopic(fid string, title string, text string) (topicaddr string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			return
		}
	}()

	re, err := this.Get(`http://myfans.info/post.php?action=newthread&fid=` + fid)
	e(`获取帖子板块失败`, err)
	defer re.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(re)
	e(`解析板块页面失败`, err)
	postaddr, ok := doc.Find(`form#postform`).Attr(`action`)
	e(`获取发帖地址失败`, ok)
	hashvalue, ok := doc.Find(`form#postform`).Find(`input[name="formhash"]`).Attr(`value`)
	e(`获取哈希失败`, ok)

	dhash, ok := doc.Find(`form#postform`).Find(`input[name="dhash"]`).Attr(`value`)
	e(`获取回复哈希值失败,dhash`, ok)

	html, _ := doc.Html()
	adhash := regexp.MustCompile(`document\.getElementById\("dhash"\)\.value="(.+)";`).FindStringSubmatch(html)
	if len(adhash) > 1 {
		dhash = adhash[1]
	}

	postv := url.Values{}
	postv.Add(`formhash`, hashvalue)
	postv.Add(`subject`, title)
	postv.Add(`message`, text)
	postv.Add(`dhash`, dhash)

	if classes := doc.Find(`select[name="typeid"]`).Children().Length(); classes != 0 {
		radint := rand.New(rand.NewSource(time.Now().Unix())).Intn(classes)
		id, ok := doc.Find(`select[name="typeid"]`).Children().Eq(radint).Attr(`value`)
		e(`获取分类id失败`, ok)
		postv.Add(`typeid`, id)
	}

	re, err = this.PostForm(root+postaddr, postv)
	e(`发帖失败`, err)
	defer re.Body.Close()
	doc, err = goquery.NewDocumentFromReader(re.Body)
	return ``, nil

	//doc, err = goquery.NewDocumentFromReader(mahonia.NewDecoder(`big5`).NewReader(re.Body))
	//e(`发帖失败`, err)
	//target, ok := doc.Find(`meta[http-equiv="refresh"]`).Attr(`content`)
	//e(`获取发帖地址失败`, ok)
	//targeturl := regexp.MustCompile(`url=(.+)`).FindStringSubmatch(target)
	//if len(targeturl) < 2 {
	//	panic(fmt.Errorf(`解析发帖地址失败`))
	//}
	//return root + targeturl[1], nil
}
func e(desc string, err interface{}) {
	switch fault := err.(type) {
	case error:
		{
			if fault != nil {
				panic(fmt.Errorf(desc+": %v", err))
			}
		}
	case bool:
		{
			if !fault {
				panic(fmt.Errorf(desc))
			}
		}
	}

}

func ReadAll(i io.Reader) string {
	b, err := ioutil.ReadAll(i)
	e(`ReadAll失败`, err)
	return string(b)
}
