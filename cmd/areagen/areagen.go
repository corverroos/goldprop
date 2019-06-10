package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

const start = "https://www.property24.com/for-sale/all-suburbs/cape-town/western-cape/432"

type area struct {
	area string
	code int
}

func main() {
	areas := make(map[int]area)
	
	c := colly.NewCollector(
		colly.AllowedDomains("www.property24.com"),
	)

	s := c.Clone()

	// Find all districts
	s.OnHTML(".p24_popular a[title]", func(e *colly.HTMLElement) {
		a := e.Text
		link := e.Attr("href")
		title := e.Attr("title")
		if !strings.Contains(link,"for-sale") {
			return
		}
		if !strings.Contains(title,a) {
			return
		}
		split := strings.Split(link, "/")
		code := split[len(split)-1]
		c, err := strconv.Atoi(code)
		if err != nil {
			return
		}
		areas[c] = area{
			area: a,
			code: c,
		}
	})

	// Find and visit all suburbs
	c.OnHTML(".row .checkbox", func(e *colly.HTMLElement) {
		code := e.ChildAttr("input","value")
		suburb := e.ChildText("a")
		c, err := strconv.Atoi(code)
		if err != nil {
			return
		}
		link := e.Request.AbsoluteURL(e.ChildAttr("a","href"))
		areas[c] = area{
			area: suburb,
			code: c,
		}
		err = s.Visit(link)
		if err != nil {
			log.Fatalf("error visiting suburb: %v", err)
		}
	})

	err := c.Visit(start)
	if err != nil {
		log.Fatal(err)
	}

	var al []area
	for _, a := range areas {
		al = append(al, a)
	}

	sort.Slice(al, func(i, j int) bool {
		return al[i].area < al[j].area
	})

	for _, a := range al {
		fmt.Printf("\"%s\":%d,\n",a.area,a.code)
	}
}
