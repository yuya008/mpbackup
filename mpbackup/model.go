package mpbackup

import (
	"net/url"
	"github.com/yuya008/mpbackup/downloader"
	"path"
	"github.com/PuerkitoBio/goquery"
	"bytes"
	"strings"
)

// 用户表
type MepaiUser struct {
	Id		int
	Avatar	string
	Cover	string
	Sn 		string
}

// 作品附件表
type MepaiWorksAppendix struct {
	Id		int64
	WorksId int64
	Src 	string
}

// 活动表
type MepaiActivity struct {
	Id int
	Cover string
	Content string
}

func (mp *Mpbackup) pushUserToDownLoader() {
	for i := 1; ; i++ {
		users := make([]MepaiUser, 0)
		err := mp.dbEngine.Limit(1000, (i - 1) * 1000).Asc("id").Find(&users)
		if err != nil {
			runlog.Panicln(err.Error())
		}
		if len(users) == 0 {
			return
		}
		for _, user := range users {
			avatarURL, err := url.Parse(URLPrefix + user.Avatar)
			if err != nil {
				runlog.Println(err.Error())
				continue
			}
			coverURL, err := url.Parse(URLPrefix + user.Cover)
			if err != nil {
				runlog.Println(err.Error())
				continue
			}
			avatarDownTask := &downloader.Task{
				Url : avatarURL.String(),
				To : path.Clean(mp.backupPath + "/" + avatarURL.EscapedPath()),
			}
			runlog.Println("推送 -> ", avatarDownTask)
			mp.downloader.Push(avatarDownTask)
			coverTask := &downloader.Task{
				Url : coverURL.String(),
				To : path.Clean(mp.backupPath + "/" + coverURL.EscapedPath()),
			}
			runlog.Println("推送 -> ", coverTask)
			mp.downloader.Push(coverTask)
		}
	}
}

func (mp *Mpbackup) pushWorksAppendixToDownLoader() {
	for i := 1; ; i++ {
		worksAppendixs := make([]MepaiWorksAppendix, 0)
		err := mp.dbEngine.Limit(1000, (i - 1) * 1000).Asc("id").Find(&worksAppendixs)
		if err != nil {
			runlog.Panicln(err.Error())
		}
		if len(worksAppendixs) == 0 {
			return
		}
		for _, worksAppendix := range worksAppendixs {
			worksAppendixURL, err := url.Parse(URLPrefix + worksAppendix.Src)
			if err != nil {
				runlog.Println(err.Error())
				continue
			}
			worksAppendixTask := &downloader.Task{
				Url : worksAppendixURL.String(),
				To : path.Clean(mp.backupPath + "/" + worksAppendixURL.EscapedPath()),
			}
			runlog.Println("推送 -> ", worksAppendixTask)
			mp.downloader.Push(worksAppendixTask)
		}
	}
}

func (mp *Mpbackup) pushActivityToDownLoader() {
	for i := 1; ; i++ {
		activitys := make([]MepaiActivity, 0)
		err := mp.dbEngine.Limit(1000, (i - 1) * 1000).Asc("id").Find(&activitys)
		if err != nil {
			runlog.Panicln(err.Error())
		}
		if len(activitys) == 0 {
			return
		}
		for _, activity := range activitys {
			coverURL, err := url.Parse(URLPrefix + activity.Cover)
			if err != nil {
				runlog.Println(err.Error())
				continue
			}
			coverTask := &downloader.Task{
				Url : coverURL.String(),
				To : path.Clean(mp.backupPath + "/" + coverURL.EscapedPath()),
			}
			runlog.Println("推送 -> ", coverTask)
			mp.downloader.Push(coverTask)
			// DOM 分析
			imgSrc, err := htmlImgSrcParser([]byte(activity.Content))
			if err != nil {
				runlog.Println(err.Error())
				continue
			}
			for _, imgS := range imgSrc {
				if !strings.Contains(imgS, "mepai") {
					continue
				}
				imgSrcURL, err := url.Parse(imgS)
				if err != nil {
					runlog.Println(err.Error())
					continue
				}
				imgSrcTask := &downloader.Task{
					Url : imgSrcURL.String(),
					To : path.Clean(mp.backupPath + "/" + imgSrcURL.EscapedPath()),
				}
				runlog.Println("推送 -> ", imgSrcTask)
				mp.downloader.Push(imgSrcTask)
			}
		}
	}
}

func htmlImgSrcParser(html []byte) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, err
	}
	retval := make([]string, 0)
	for _, node := range doc.Find("img").Nodes {
		for _, attr := range node.Attr {
			if attr.Key == "src" {
				retval = append(retval, attr.Val)
			}
		}
	}
	return retval, nil
}
