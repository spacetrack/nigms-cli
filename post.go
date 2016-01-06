/*
 * nigms-cli - nun-ist-genug-mit-schnee command line interface
 *
 * post
 *
 */

package main

import (
    "net/url"
	"strings"
	"time"
)

type Post struct {
	Id       string   `yaml: "id"`
	Type     string   `yaml: "type"`
	Status   string   `yaml: "status"` // draft, published
	Title    string   `yaml: "title"`
	Body     string   `yaml: "body"`
	Tags     []string `yaml: "tags"`
	time_str string   `yaml: time`
	Time     time.Time
}

func (p *Post) GetTumblrApiValues() url.Values {
	values := url.Values{}

	if len(p.Id) > 0 {
		values.Set("id", p.Id)
	}

	values.Set("type", p.Type)
	values.Set("state", p.Status)
	values.Set("title", p.Title)
	values.Set("body", p.Body)
	values.Set("tags", strings.Join(p.Tags, ","))
	values.Set("date", p.time_str)

	return values
}
