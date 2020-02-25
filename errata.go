package main

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// Errata stores all of the erratas grouped by date
type Errata struct {
	ID       int
	Type     string
	Date     time.Time
	Coverage string
	Desc     string
	Link     string
}

// Errati is a collection of Errata
type Errati struct {
	List   []Errata
	Length int
}

// By is the type of a "less" function that defines the ordering of its Planet arguments.
type By func(p1, p2 *Errata) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(erratas []Errata) {
	es := &errataSorter{
		erratas: erratas,
		by:      by,
	}
	sort.Sort(es)
}

// errataSorter joins a By function and a slice of Errati to be sorted.
type errataSorter struct {
	erratas []Errata
	by      func(p1, p2 *Errata) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *errataSorter) Len() int {
	return len(s.erratas)
}

// Swap is part of sort.Interface.
func (s *errataSorter) Swap(i, j int) {
	s.erratas[i], s.erratas[j] = s.erratas[j], s.erratas[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *errataSorter) Less(i, j int) bool {
	return s.by(&s.erratas[i], &s.erratas[j])
}

// ParseErrata does the actual parsing of html (poorly!)
func ParseErrata(body io.Reader) (*Errati, error) {
	var data []Errata
	var errati = &Errati{}
	doc, err := html.Parse(body)
	if err != nil {
		return nil, err
	}
	var f func(*html.Node)
	f = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "strong" {
			if node.FirstChild != nil {
				var e Errata
				var err error
				parts := strings.Split(node.FirstChild.Data, ": ")
				e.ID, err = strconv.Atoi(parts[0])
				if err != nil {
					return
				}
				e.Type = parts[1]
				e.Date, err = time.Parse("January 02, 2006", parts[2])
				if err != nil {
					return
				}

				// TODO: not this.

				if node.NextSibling != nil &&
					node.NextSibling.NextSibling != nil &&
					node.NextSibling.NextSibling.FirstChild.Data != "" {
					e.Coverage = node.NextSibling.NextSibling.FirstChild.Data
				}

				if node.NextSibling != nil &&
					node.NextSibling.NextSibling != nil &&
					node.NextSibling.NextSibling.NextSibling != nil &&
					node.NextSibling.NextSibling.NextSibling.NextSibling != nil &&
					node.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling != nil &&
					node.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.Data != "" {
					e.Desc = node.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.Data
					e.Desc = strings.TrimRight(e.Desc, "\n")
					e.Desc = strings.TrimLeft(e.Desc, "\n")
				}

				if node.NextSibling != nil &&
					node.NextSibling.NextSibling != nil &&
					node.NextSibling.NextSibling.NextSibling != nil &&
					node.NextSibling.NextSibling.NextSibling.NextSibling != nil &&
					node.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling != nil &&
					node.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling != nil &&
					node.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling != nil &&
					node.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling != nil {
					n := node.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling
					for _, a := range n.Attr {
						if a.Key == "href" {
							e.Link = a.Val
							break
						}
					}
				}

				data = append(data, e)
				errati.Length = errati.Length + 1
				return
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			f(child)
		}
	}
	f(doc)

	byDate := func(p1, p2 *Errata) bool {
		return p1.Date.Unix() < p2.Date.Unix()
	}

	By(byDate).Sort(data)
	errati.List = data

	return errati, nil
}

// ParseRemoteErrata grabs all of the OpenBSD errata from an html page
func ParseRemoteErrata(s string) (*Errati, error) {
	resp, err := http.Get(s)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ParseErrata(resp.Body)
}

// PrintErrata pretty prints an errata
func PrintErrata(e *Errata) string {
	return fmt.Sprintf("%03d: %s: %s %s\n%s\n%s",
		e.ID,
		e.Type,
		e.Date.String(),
		e.Coverage,
		e.Desc,
		e.Link,
	)
}

// PrintErrataMD pretty prints an errata in markdown
//func PrintErrataMD(e *Errata) string {
//	return fmt.Sprintf("#%03d: %s: %s _%s_\n%s\n[A source code patch exists which remedies this problem.](%s)",
//		e.ID,
//		e.Type,
//		e.Date.String(),
//		e.Coverage,
//		e.Desc,
//		e.Link,
//	)
//}
