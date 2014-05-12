package uwants

import (
	"client"
	//"github.com/astaxie/beego"
	//"code.google.com/p/go.net/html"
	"code.google.com/p/mahonia"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"math/rand"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var Proxy string = `http://127.0.0.1:8087`
var root = `http://www.uwants.com/`

type Uwants struct {
	*client.Client
	Decoder  mahonia.Decoder
	username string
	password string
}

func New(username, password string) *Uwants {
	u := &Uwants{}
	u.Client = client.New()
	u.Decoder = mahonia.NewDecoder(`big5`)
	if Proxy != `` {
		u.UseProxy(Proxy)
	}
	u.UseEncoder(`big5`)
	u.username = strings.TrimSpace(username)
	u.password = strings.TrimSpace(password)
	return u
}

func (this *Uwants) Login() (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			return
		}
	}()
	re, err := this.Get(root + `index.php`)
	e(`获取登陆地址失败`, err)
	defer re.Body.Close()

	doc, err := goquery.NewDocumentFromReader(mahonia.NewDecoder(`big5`).NewReader(re.Body))
	e(`解析首页失败`, err)

	loginaddr, ok := doc.Find(`form[name="loginform"]`).Attr(`action`)
	e(`解析登陆地址失败`, ok)
	v := url.Values{}
	v.Add(`username`, this.username)
	v.Add(`password`, this.password)
	_, err = this.PostForm(root+loginaddr, v)
	e(`登陆失败`, err)
	return nil
}

func (this *Uwants) SendReply(tid string, title string, text string) (returl string, err error) {
	var doc *goquery.Document
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			return
		}
	}()
	re, err := this.Get(`http://www.uwants.com/viewthread.php?tid=` + strings.TrimSpace(tid) + `&extra=page%3D1`)
	e(`获取帖子首页失败`, err)
	defer re.Body.Close()

	doc, err = goquery.NewDocumentFromReader(this.Decoder.NewReader(re.Body))
	e(`创建页面索引失败`, err)

	ht, _ := doc.Html()

	hashvalue, ok := doc.Find(`#postform`).Find(`input[name="formhash"]`).Attr(`value`)
	if !ok {
		if strings.Contains(ht, `被禁用`) {
			e(`账号被禁用`, false)
		}
	}
	e(`获取回复哈希值失败`, ok)

	actionvalue, ok := doc.Find(`#postform`).Attr(`action`)
	e(`获取回复地址失败`, ok)

	retv := url.Values{}
	retv.Add(`subject`, strings.TrimSpace(title))
	retv.Add(`message`, strings.TrimSpace(text))
	retv.Add(`formhash`, hashvalue)

	re, err = this.PostForm(`http://www.uwants.com/`+actionvalue, retv)
	e(`回复失败`, err)
	defer re.Body.Close()
	doc, err = goquery.NewDocumentFromReader(mahonia.NewDecoder(`big5`).NewReader(re.Body))
	e(`回复失败`, err)

	target, ok := doc.Find(`meta[http-equiv="refresh"]`).Attr(`content`)
	e(`获取回复地址失败`, ok)
	targeturl := regexp.MustCompile(`url=(.+)`).FindStringSubmatch(target)
	if len(targeturl) < 2 {
		panic(fmt.Errorf(`获取回复地址失败`))
	}
	return root + targeturl[1], nil

}

func (this *Uwants) NewTopic(fid string, title string, text string) (topicaddr string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			return
		}
	}()

	re, err := this.Get(`http://www.uwants.com/forumdisplay.php?fid=` + strings.TrimSpace(fid))
	e(`获取帖子板块失败`, err)
	defer re.Body.Close()

	doc, err := goquery.NewDocumentFromReader(mahonia.NewDecoder(`big5`).NewReader(re.Body))
	e(`解析板块页面失败`, err)

	ht, _ := doc.Html()
	postaddr, ok := doc.Find(`form#postform`).Attr(`action`)
	//if !ok {
	//	beego.Debug(doc.Html())
	//}
	if !ok {
		if strings.Contains(ht, `被禁用`) {
			e(`账号被禁用`, false)
		}
	}
	e(`获取发帖地址失败`, ok)

	hashvalue, ok := doc.Find(`form#postform`).Find(`input[name="formhash"]`).Attr(`value`)
	e(`获取哈希失败`, ok)
	postv := url.Values{}
	postv.Add(`formhash`, hashvalue)
	postv.Add(`subject`, strings.TrimSpace(title))
	postv.Add(`message`, strings.TrimSpace(text))
	if classes := doc.Find(`select[name="typeid"]`).Children().Length(); classes != 0 {
		radint := rand.New(rand.NewSource(time.Now().Unix())).Intn(classes)
		id, ok := doc.Find(`select[name="typeid"]`).Children().Eq(radint).Attr(`value`)
		e(`获取分类id失败`, ok)
		postv.Add(`typeid`, id)
	}

	re, err = this.PostForm(root+postaddr, postv)
	e(`发帖失败`, err)
	defer re.Body.Close()
	doc, err = goquery.NewDocumentFromReader(mahonia.NewDecoder(`big5`).NewReader(re.Body))
	e(`发帖失败`, err)
	target, ok := doc.Find(`meta[http-equiv="refresh"]`).Attr(`content`)
	//if !ok {
	//	beego.Debug(doc.Html())
	//}
	e(`获取发帖后地址失败`, ok)
	targeturl := regexp.MustCompile(`url=(.+)`).FindStringSubmatch(target)
	if len(targeturl) < 2 {
		panic(fmt.Errorf(`解析发帖后地址失败`))
	}
	return root + targeturl[1], nil
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

//func main() {
//	cl := http.Client{}
//	cl.Transport = &http.Transport{Proxy: func(*http.Request) (*url.URL, error) {
//		v, e := url.Parse(`http://127.0.0.1:8087`)
//		fmt.Println(v, e)
//		return v, e
//	}}
//	j, _ := cookiejar.New(nil)
//	cl.Jar = j
//	v := url.Values{}
//	v.Add(`username`, `doomsplayer`)
//	v.Add(`password`, `1cd3599df`)
//	re, err := cl.PostForm(`http://www.uwants.com/logging.php?action=login&loginsubmit=true`, v)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	b, _ := ioutil.ReadAll(re.Body)
//	s := mahonia.NewDecoder(`big5`).ConvertString(string(b))
//	re, _ = cl.Get(`http://www.uwants.com/viewthread.php?tid=17318841&extra=page%3D1`)
//	b, _ = ioutil.ReadAll(re.Body)
//	reg := regexp.MustCompile(`<input type="hidden" name="formhash" value="(.*)" /?>`).FindStringSubmatch(string(b))
//	fmt.Println(reg)
//	v = url.Values{}
//	v.Add(`message`, `好可爱哟`)
//	v.Add(`formhash`, reg[1])
//	re, _ = cl.PostForm(`http://www.uwants.com/post.php?action=reply&fid=35&tid=17318841&extra=page%3D1&replysubmit=yes`, v)
//	b, _ = ioutil.ReadAll(re.Body)
//	s = mahonia.NewDecoder(`big5`).ConvertString(string(b))
//	fmt.Println(s)
//}
