package main

import (
	"errors"
	"html"
	"net/http"
	"time"
)

func newArticle(w http.ResponseWriter, r *http.Request) {
	var feedback string
	article := &Article{}
	if r.Method == "POST" {
		feedback = "Article sent"
		r.ParseForm()
		var err error
		if article, err = genArticle(r); err != nil {
			feedback = "Oops...! " + err.Error()
		} else {
			// EventStart: newArticle
			artMgr.setArticle(article)
			db.sync(articlePrefix, article)
			go updateIndexAndFeed()
			// EventEnd: newArticle
		}
	}
	if err := tmpl.ExecuteTemplate(w, "new", map[string]interface{}{
		"config":   config,
		"feedback": feedback,
		"form":     article,
		"header":   "new",
	}); err != nil {
		logger.Println("new:", err.Error())
	}
}

func genArticle(r *http.Request) (*Article, error) {
	if !checkKeyExist(r.Form, "author", "email", "content", "title") {
		logger.Println("new:", "required field not exists")
		return nil, errors.New("required field not exists")
	}
	// may we need a filter?
	return &Article{
		Id:         artMgr.allocArticleId(),
		Author:     html.EscapeString(r.Form.Get("author")),
		Email:      html.EscapeString(r.Form.Get("email")),
		Website:    genWebsite(r.Form.Get("website")),
		RemoteAddr: r.RemoteAddr,
		Date:       time.Now(),
		Title:      html.EscapeString(r.Form.Get("title")),
		Content:    r.Form.Get("content"),
		QuoteNotif: r.Form.Get("notify") == "on",
	}, nil
}

func init() {
	if config["AdminUrl"][len(config["AdminUrl"])-1] == '/' {
		config["AdminUrl"] = config["AdminUrl"][:len(config["AdminUrl"])-1]
	}
	http.HandleFunc(config["RootUrl"]+config["AdminUrl"], newArticle)
	// http.HandleFunc(config["AdminUrl"] + "modify", modifyArticle)
}
