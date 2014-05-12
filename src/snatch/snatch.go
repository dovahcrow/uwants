package snatch

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"loginpool"
	//"strings"
	"html"
	"log"
	"regexp"
	"strconv"
	. "tools"
)

var root = `http://www.uwants.com/`

func GetFid() (fids []string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	cl := loginpool.GetUwants(`doomsplayer`, `1cd3599df`)
	re, err := cl.Get(root + `index2.php`)
	E(`爬虫:获取首页失败`, err)
	defer re.Body.Close()
	doc, err := goquery.NewDocumentFromReader(cl.Decoder.NewReader(re.Body))
	E(`爬虫:分析首页失败`, err)
	as := doc.Find("a")

	content := as.Map(func(i int, q *goquery.Selection) string {
		s, _ := q.Attr("href")
		return s
	})

	reg := regexp.MustCompile(`forumdisplay\.php\?fid=(\d+)`)
	dumpmap := map[string]struct{}{}

	for _, v := range content {
		if t := reg.FindStringSubmatch(v); len(t) > 1 {
			dumpmap[t[1]] = struct{}{}
		}
	}

	for k := range dumpmap {
		fids = append(fids, k)
	}
	return

}

func ChkFidAv(fids []string) (fidsav []string, err error) {
	cl := loginpool.GetUwants(`doomsplayer`, `1cd3599df`)
	for _, fid := range fids {
		re, err := cl.Get(root + `forumdisplay.php?fid=` + fid)
		if err != nil {
			continue
		}
		doc, err := goquery.NewDocumentFromResponse(re)
		if err != nil {
			continue
		}
		_, ok := doc.Find(`form#postform`).Attr(`action`)
		if ok {
			fidsav = append(fidsav, fid)
		}
		re.Body.Close()
	}
	return
}

func GetTid(fid string) (tids []string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	cl := loginpool.GetUwants(`doomsplayer`, `1cd3599df`)

	re, err := cl.Get(root + `forumdisplay.php?fid=` + fid)
	E(fmt.Sprintf("爬虫:获取板块页 %v 失败", fid), err)
	defer re.Body.Close()

	doc, err := goquery.NewDocumentFromReader(cl.Decoder.NewReader(re.Body))
	E(fmt.Sprintf("爬虫:分析板块页 %v 失败", fid), err)

	flast, ok := doc.Find(".pages>a.last").Attr(`href`)
	if !ok {
		log.Print(fmt.Sprintf("爬虫:分析板块页最后页 %v 失败,尝试分析最后一页", fid))
	}
	flast, ok = doc.Find(".pages>a").Filter(".next").Last().Attr(`href`)
	if !ok {
		log.Print(fmt.Sprintf("爬虫:分析板块页最后页 %v 失败,采用一页", fid))
	}
	flast = fid + `&page=1`

	reg := regexp.MustCompile(`forumdisplay\.php\?fid=\d+&page=(\d+)`)

	slastpage := reg.FindStringSubmatch(html.UnescapeString(flast))

	if len(slastpage) < 2 {
		E(fmt.Sprintf("爬虫:分析板块页最后页 %v 失败", fid), false)
	}

	lastpage, err := strconv.Atoi(slastpage[1])
	E(fmt.Sprintf("爬虫:分析板块页最后页 %v 失败", fid), err)

	reg = regexp.MustCompile(`thread_(\d+)`)
	for i := 1; i < lastpage+1; i++ {
		re, err := cl.Get(root + fid + `&page=` + fmt.Sprint(i))
		if err != nil {
			//TODO handle err
			continue
		}
		doc, err := goquery.NewDocumentFromReader(cl.Decoder.NewReader(re.Body))
		if err != nil {
			//TODO error handle
			continue
		}

		doc.Find(`span.tsubject`).Each(
			func(i int, q *goquery.Selection) {
				ids, _ := q.Attr(`id`)

				th := reg.FindStringSubmatch(ids)
				if len(th) < 2 {
					//TODO error handle
					return
				}
				tids = append(tids, th[1])
			})
	}
	return

}

func ChkTidAv(tids []string) (tidsav []string, err error) {
	cl := loginpool.GetUwants(`doomsplayer`, `1cd3599df`)
	for _, tid := range tids {
		re, err := cl.Get(root + `viewthread.php?tid=` + tid)
		if err != nil {
			continue
		}
		doc, err := goquery.NewDocumentFromResponse(re)
		if err != nil {
			continue
		}
		_, ok := doc.Find(`form#postform`).Attr(`action`)
		if ok {
			tidsav = append(tidsav, tid)
		}
		re.Body.Close()
	}
	return
}
