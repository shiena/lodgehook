package hook

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/shiena/lodgehook"
)

const (
	idobataURL = "https://idobata.io/hook/"
)

type IdobataHook struct {
	token string
}

func NewIdobataHook(token string) *IdobataHook {
	return &IdobataHook{token: token}
}

func (i *IdobataHook) post(values url.Values) (*http.Response, error) {
	if len(i.token) <= 0 {
		return nil, nil
	}
	return http.PostForm(idobataURL+i.token, values)
}

func (i *IdobataHook) PostText(message string) (*http.Response, error) {
	values := url.Values{}
	values.Add("source", message)
	return i.post(values)
}

func (i *IdobataHook) PostHtml(html string) (*http.Response, error) {
	values := url.Values{}
	values.Add("format", "html")
	values.Add("source", html)
	return i.post(values)
}

func (i *IdobataHook) Hook(article *lodgehook.LodgeArticle) {
	message := formatMessage(article)
	log.Println(message)
	i.PostHtml(message)
}

func formatMessage(article *lodgehook.LodgeArticle) string {
	var tag string
	if len(article.Tags) > 0 {
		tag = " [" + strings.Join(article.Tags, "] [") + "]"
	}
	return fmt.Sprintf(`[New] <a href="%s">%s</a>%s`, article.Loc, article.Title, tag)
}
