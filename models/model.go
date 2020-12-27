package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/gorilla/feeds"
	"github.com/lxneng/refeed/config"
	"github.com/mitchellh/mapstructure"
)

type Feed struct {
	Name  string
	Url   string
	Xpath string
}

type ZhiJsonObject struct {
	ZhihuPosts []ZhihuPost `json:"data"`
}

type ZhihuPost struct {
	Title   string
	Url     string
	Content string
	Excerpt string
	Type    string
}

func GetAtomFeed(slug string) (atom string, err error) {
	c := config.GetConfig()
	feeds := c.Get("feeds").(map[string]interface{})

	var feed Feed
	val, ok := feeds[slug]
	if !ok {
		err = errors.New("feed not found!")
		return
	}

	if err = mapstructure.Decode(val, &feed); err != nil {
		return
	}
	atom, err = feed.GenerateFeed()
	return
}

func (f *Feed) GenerateFeed() (atom string, err error) {
	if strings.HasPrefix(f.Url, "https://www.zhihu.com/column") {
		atom, err = f.GenerateFeedZhiHu()
		return
	}
	doc, err := htmlquery.LoadURL(f.Url)
	if err != nil {
		return
	}

	nodes, err := htmlquery.QueryAll(doc, f.Xpath)
	if err != nil {
		return
	}

	var items []*feeds.Item
	now := time.Now()

	fd := &feeds.Feed{
		Title:       f.Name,
		Link:        &feeds.Link{Href: f.Url},
		Description: f.Name,
		Created:     now,
	}

	u, err := url.Parse(f.Url)
	if err != nil {
		return
	}
	rootURL := fmt.Sprintf("%s://%s", u.Scheme, u.Hostname())

	for _, node := range nodes {
		url := htmlquery.SelectAttr(node, "href")
		title := htmlquery.InnerText(node)
		if strings.HasPrefix(url, "/") {
			url = fmt.Sprintf("%s%s", rootURL, url)
		}
		items = append(items, &feeds.Item{
			Title:   title,
			Link:    &feeds.Link{Href: url},
			Id:      url,
			Created: now,
		})
	}
	fd.Items = items
	atom, err = fd.ToAtom()
	if err != nil {
		return
	}
	return
}

func (f *Feed) GenerateFeedZhiHu() (atom string, err error) {
	slist := strings.Split(f.Url, "/")
	zlid := slist[len(slist)-1]
	apiUrl := fmt.Sprintf("https://www.zhihu.com/api/v4/columns/%s/items", zlid)
	res, err := http.Get(apiUrl)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return
	}

	var jcontent ZhiJsonObject
	json.Unmarshal(body, &jcontent)

	var items []*feeds.Item
	now := time.Now()

	fd := &feeds.Feed{
		Title:   f.Name,
		Link:    &feeds.Link{Href: f.Url},
		Created: now,
	}

	for _, post := range jcontent.ZhihuPosts {
		if post.Type == "answer" {
			continue
		}
		items = append(items, &feeds.Item{
			Title:       post.Title,
			Link:        &feeds.Link{Href: post.Url},
			Id:          post.Url,
			Description: post.Excerpt,
			Content:     post.Content,
			Created:     now,
		})
	}
	fd.Items = items
	atom, err = fd.ToAtom()
	if err != nil {
		return
	}
	return
}
