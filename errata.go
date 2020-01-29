package main

import (
	"net/http"

	"golang.org/x/net/html"
)

// ParseErrata grabs all of the OpenBSD errata from an html page
func ParseErrata(s string) ([]string, error) {
	var data []string
	resp, err := http.Get(s)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var f func(*html.Node)
	f = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "strong" {
			if node.FirstChild != nil {
				data = append(data, node.FirstChild.Data)
				return
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			f(child)
		}
	}
	f(doc)

	return data, nil
}
