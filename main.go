package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Har struct {
	Log *HarLog
}

type HarLog struct {
	Version string
	Pages   []*Page
	Entries []*Entry
}

type Page struct {
	StartedDateTime string
	Id              string
	Title           string
}

type Entry struct {
	Pageref  string
	Request  *Request
	Response *Response
}

type Request struct {
	Url         string
	Method      string
	Headers     []*Header `json:"headers,omitempty"`
	QueryString []*Header `json:"queryString,omitempty"`
	Cookies     []*Cookie `json:"cookies,omitempty"`
	PostData    *PostData `json:"postData,omitempty"`
	SlimHeaders []string
	Body        string `json:"body,omitempty"`
}

type Response struct {
	Status      int
	Headers     []*Header `json:"headers,omitempty"`
	Cookies     []*Cookie `json:"cookies,omitempty"`
	Content     *PostData `json:"content,omitempty"`
	SlimHeaders []string
	Body        string
}

type Header struct {
	Name  string
	Value string
}

type Cookie struct {
	Name     string
	Value    string
	Path     string
	HttpOnly bool
	Secure   bool
}

type PostData struct {
	MimeType string
	Text     string
}

func main() {
	content, err := os.ReadFile("/Users/danielgabay/repos/har-converter/hdfc2.har")
	if err != nil {
		fmt.Printf("error reading file: %v", err)
		return
	}
	originLog := Har{}
	if err := json.Unmarshal(content, &originLog); err != nil {
		fmt.Printf("error unmarshalling file: %v", err)
		return
	}

	origin := originLog.Log
	if len(origin.Version) == 0 {
		fmt.Printf("error: missing version in HAR file")
		return
	}
	excludeMime := []string{"font", "image", "video", "javascript"}

	modified := *origin
	for i := 0; i < len(modified.Entries); i++ {
		req := modified.Entries[i].Request
		res := origin.Entries[i].Response
		if inArray(res.Content.MimeType, excludeMime) {
			modified.Entries[i] = nil
			continue
		}

		req.SlimHeaders = slimHeaders(req.Headers)
		req.Headers = nil //make([]*Header, 0)
		req.Cookies = nil //make([]*Cookie, 0)
		req.QueryString = nil
		if req.PostData != nil && len(req.PostData.Text) > 0 {
			req.Body = req.PostData.Text
		}
		req.PostData = nil

		res.SlimHeaders = slimHeaders(res.Headers)
		res.Headers = nil //make([]*Header, 0)
		res.Cookies = nil //make([]*Cookie, 0)
		if res.Content != nil && len(res.Content.Text) > 0 {
			res.Body = res.Content.Text
		}
		res.Content = nil
	}

	updated, err := json.MarshalIndent(modified, "", "    ")
	if err != nil {
		return
	}
	if err := os.WriteFile("/Users/danielgabay/repos/har-converter/AuthenticationX.txt", updated, 666); err != nil {
		fmt.Printf("error writing file: %v", err)
		return
	}
}

func slimHeaders(headers []*Header) []string {
	slimHeaders := []string{}
	for i := 0; i < len(headers); i++ {
		header := headers[i]
		slimHeaders = append(slimHeaders, fmt.Sprintf("%s:%s", header.Name, header.Value))
	}
	return slimHeaders
}

func inArray(item string, array []string) bool {
	item = strings.ToLower(item)
	for _, entry := range array {
		if strings.Contains(item, entry) {
			return true
		}
	}
	return false
}
