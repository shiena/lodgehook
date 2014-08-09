package lodgehook

import (
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var newPostLocation = regexp.MustCompile("/articles/[0-9]+$")

type LodgeArticle struct {
	UserId string
	Title  string
	Tags   []string
	Body   string
	Loc    string
}

type LodgeHook interface {
	Hook(article *LodgeArticle)
}

type LodgeHookTransport struct {
	LodgeHook []LodgeHook
	Transport http.RoundTripper
}

func getArticle(req *http.Request) (article *LodgeArticle, err error) {
	if strings.ToUpper(req.Method) == "POST" {
		if req.RequestURI == "/articles" {
			var post url.Values
			post, err = dumpPostForm(req)
			if err == nil {
				article = &LodgeArticle{}
				article.UserId = post.Get("article[user_id]")
				article.Title = post.Get("article[title]")
				article.Body = post.Get("article[body]")
				tag_list := post.Get("article[tag_list]")
				if tag_list != "" {
					article.Tags = strings.Split(tag_list, ", ")
				}
			}
		}
	}
	return
}

func getLoc(resp *http.Response) (location string) {
	if resp.StatusCode == 302 {
		loc := resp.Header.Get("Location")
		if newPostLocation.MatchString(loc) {
			location = loc
		}
	}
	return
}

func (t *LodgeHookTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {

	article, err := getArticle(req)
	if err != nil {
		log.Println(err)
	}

	resp, err = t.Transport.RoundTrip(req)
	log.Printf("%s %s", req.Method, req.RequestURI)

	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("%d %s", resp.StatusCode, resp.Status)
	if article == nil {
		return
	}

	loc := getLoc(resp)
	if loc == "" {
		return
	}

	article.Loc = loc
	for _, hook := range t.LodgeHook {
		go hook.Hook(article)
	}
	return
}

func NewLodgeHookTransport(hook ...LodgeHook) *LodgeHookTransport {
	return &LodgeHookTransport{
		LodgeHook: hook,
		Transport: http.DefaultTransport,
	}
}
