// Package pastebin wraps the basic functions of the Pastebin API and exposes a
// Go API.
package pastebin

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	pastebinDevKey = "7b9c033d5a4e4b417fd2a22e6d598b01"

	paste_public = "0"
	paste_unlisted = "1"
	paste_private = "2"

	expire_never = "N"
	expire_10m = "10M"
	expire_1H = "1H"
	expire_1D = "1D"
	expire_1w = "1W"
	expire_2w = "2W"
	expire_month = "1M"
)

var (
	// ErrPutFailed is returned when a paste could not be uploaded to pastebin.
	ErrPutFailed = errors.New("pastebin put failed")
	// ErrGetFailed is returned when a paste could not be fetched from pastebin.
	ErrGetFailed = errors.New("pastebin get failed")
)

// Pastebin represents an instance of the pastebin service.
type Pastebin struct{}

// Put uploads text to Pastebin with optional title returning the ID or an error.
func (p Pastebin) Put(text, title, highlight string) (id string, err error) {
	data := url.Values{}
	highlight = detectHighlight(highlight)
	// Required values.
	data.Set("api_dev_key", pastebinDevKey)
	data.Set("api_option", "paste") // Create a paste.
	data.Set("api_paste_code", text)
	// Optional values.
	data.Set("api_paste_name", title)      // The paste should have title "title".
	data.Set("api_paste_private", paste_unlisted)
	data.Set("api_paste_expire_date", expire_1H)
	data.Set("api_paste_format", highlight)

	resp, err := http.PostForm("http://pastebin.com/api/api_post.php", data)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", ErrPutFailed
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return p.StripURL(string(respBody)), nil
}

// Get returns the text inside the paste identified by ID.
func (p Pastebin) Get(id string) (text string, err error) {
	resp, err := http.Get("http://pastebin.com/raw.php?i=" + id)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", ErrGetFailed
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(respBody), nil
}

// StripURL returns the paste ID from a pastebin URL.
func (p Pastebin) StripURL(url string) string {
	return strings.Replace(url, "http://pastebin.com/", "", -1)
}

// WrapID returns the pastebin URL from a paste ID.
func (p Pastebin) WrapID(id string) string {
	return "http://pastebin.com/" + id
}

func detectHighlight(extFileName string) string {
	extFileName = strings.ToLower(extFileName)
	mapHighlighting := map[string] string {
		"":			"text",
		".php": 	"php",
		".pas":		"delphi",
		".rb": 		"ruby",
		".as3": 	"actionscript3",
		".asm": 	"asm",
		".sh": 		"bash",
		".c": 		"c",
		".cpp": 	"cpp",
		".css": 	"css",
		".erl": 	"erlang",
		".go": 		"go",
		".lua": 	"lua",
		".html5": 	"html5",
		".html": 	"html5",
		".htm": 	"html5",
		".nginx": 	"nginx",
		".xml": 	"xml",
	}

	r, ok := mapHighlighting[extFileName];
	if !ok {
		r = "text"
	}

	return r
}
