package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// Erratum are our individual chunks of errata
type Erratum struct {
	ID    int
	Date  time.Time
	Desc  string
	Patch string
	Link  string
	Sig   []string
}

// Fetch pulls down and parses the information for a given Erratum
func (e *Erratum) Fetch() error {
	resp, err := http.Get(e.Link)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	// First two lines are our signature
	// 3rd line is date
	// everything from the date down to ^Index is our description
	count := 0
	var descr []string
	var patch []string
	re := regexp.MustCompile(`^Index:`)
	matched := false
	for scanner.Scan() {
		s := scanner.Text()
		if count < 2 {
			e.Sig = append(e.Sig, s)
		}

		if count == 3 {
			parts := strings.Split(s, ",")
			if len(parts) == 3 {
				d := fmt.Sprintf("%s,%s",
					strings.Trim(parts[1], " "),
					strings.Replace(parts[2], ":", "", -1))
				e.Date, err = time.Parse("January 2, 2006", d)
				if err != nil {
					return err
				}
			}
		}

		if count > 3 {
			if re.MatchString(s) {
				matched = true
			}

			if !matched {
				descr = append(descr, s)
			} else {
				patch = append(patch, s)
			}

		}

		count = count + 1
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	e.Desc = strings.Join(descr, "\n")
	e.Patch = strings.Join(patch, "\n")

	return nil //fmt.Errorf("crap")
}

// Errata is a collection of Errata
type Errata struct {
	List   []Erratum
	Length int
}

// By is the type of a "less" function that defines the ordering of its Planet arguments.
type By func(p1, p2 *Erratum) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(errata []Erratum) {
	es := &errataSorter{
		errata: errata,
		by:     by,
	}
	sort.Sort(es)
}

// errataSorter joins a By function and a slice of Errata to be sorted.
type errataSorter struct {
	errata []Erratum
	by     func(p1, p2 *Erratum) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *errataSorter) Len() int {
	return len(s.errata)
}

// Swap is part of sort.Interface.
func (s *errataSorter) Swap(i, j int) {
	s.errata[i], s.errata[j] = s.errata[j], s.errata[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *errataSorter) Less(i, j int) bool {
	return s.by(&s.errata[i], &s.errata[j])
}

// ParseErrata does the actual parsing of html
func ParseErrata(body io.Reader, baseURL string) (*Errata, error) {
	var data []Erratum
	var errata = &Errata{}
	doc, err := html.Parse(body)
	if err != nil {
		return nil, err
	}
	var f func(*html.Node)
	f = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			if node.FirstChild != nil {
				var e Erratum
				var err error

				for _, a := range node.Attr {
					if a.Key == "href" {
						parts := strings.Split(a.Val, "_")

						if len(parts) >= 2 && a.Val != "" {
							e.Link = fmt.Sprintf("%s%s", baseURL, a.Val)
							e.ID, err = strconv.Atoi(parts[0])
							if err != nil {
								// Not a link we care about
								break
							}
							data = append(data, e)
						}
					}
				}
				return
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			f(child)
		}
	}
	f(doc)

	byID := func(p1, p2 *Erratum) bool {
		return p1.ID < p2.ID
	}

	By(byID).Sort(data)
	errata.List = data

	return errata, nil
}

// ParseRemoteErrata grabs all of the OpenBSD errata from an html page
func ParseRemoteErrata(s string) (*Errata, error) {
	resp, err := http.Get(s)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ParseErrata(resp.Body, s)
}

// PrintErrata pretty prints an errata
func PrintErrata(e *Erratum) string {
	return fmt.Sprintf("New OpenBSD Errata: %03d\n%s: %s\n%s",
		e.ID,
		e.Date.String(),
		e.Desc,
		e.Link,
	)
}

// PrintErrataMD pretty prints an errata in markdown
func PrintErrataMD(e *Erratum) string {
	return fmt.Sprintf("# OpenBSD Errata %03d (_%s_)\n<pre>%s</pre>\n[A source code patch exists which remedies this problem.](%s)",
		e.ID,
		e.Date.Format("January 2, 2006"),
		e.Desc,
		e.Link,
	)
}
