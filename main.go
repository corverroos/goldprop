package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/corverroos/goldprop/p24"
	"github.com/go-stack/stack"
	"github.com/gocolly/colly"
)

type tile struct {
	URL       string
	Price     string `selector:".p24_price"`
	Location  string `selector:".p24_location"`
	Bedrooms  string `selector:"[title='Bedrooms']"`
	Bathrooms string `selector:"[title='Bathrooms']"`
	Garages   string `selector:"[title='Garages']"`
	Size      string `selector:".p24_size"`
	Address   string `selector:".p24_address"`
}

type description struct {
	URL       string
	Description   string
}

type feature struct {
	URL       string
	Key   string
	Value   string
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.property24.com"),
		colly.CacheDir("/tmp/colly"),
	)
	maybeLog(c.Limit(&colly.LimitRule{
		RandomDelay: 5*time.Second,
		Parallelism: 2,
		DomainGlob: "www.property24.com.*",
	}))

	// TODO(corver): Process and results
	ch := make(chan interface{}, 10)

	// Visit all result tiles
	c.OnHTML(".js_resultTile .p24_content", func(e *colly.HTMLElement) {
		link, ok := e.DOM.Parent().Attr("href")
		if !ok {
			log.Printf("Cannot find parent url: %v", e)
		}
		if !strings.Contains(link, "for-sale") {
			return
		}
		var t tile
		maybeLog(e.Unmarshal(&t))
		t.URL = e.Request.AbsoluteURL(link)
		ch <- t
		maybeLog(c.Visit(t.URL))
	})

	// Scan detailed description
	c.OnHTML(".p24_listing .p24_dPl", func(e *colly.HTMLElement) {
		ch <- description{
			URL: e.Request.URL.String(),
			Description: e.Text,
		}
	})

	// Scan detailed features
	c.OnHTML(".p24_listing .p24_features .row", func(e *colly.HTMLElement) {
		f := feature{URL: e.Request.URL.String()}
		e.ForEach("div", func(i int, e *colly.HTMLElement) {
			if i == 0 {
				f.Key = e.Text
			} else if i == 2 {
				f.Value = e.Text
			}
		})
		ch <- f
	})

	// Visit next pages
	c.OnHTML(".p24_pager a.pull-right[data-pagenumber]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if !strings.Contains(link, "for-sale") {
			log.Printf("Unexpected next page: %v", e)
		}
		maybeLog(c.Visit(link))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	search, err := makeSearchUrl(0, "Marina Da Gama")
	maybeLog(err)
	maybeLog(c.Visit(search))
close(ch)

	c.Wait()
}

func maybeLog(err error) {
	if err == nil {
		return
	}
	c := stack.Caller(1)
	log.Printf("%v: Error: %v", c, err)
}

var (
	search = "https://www.property24.com/apartments-for-sale/advanced-search/results?"
)

func makeSearchUrl(minPrice int, arealist ...string) (string, error) {
	v := url.Values{}
	var codes []string
	for _, a := range arealist {
		c, ok := p24.Areas[a]
		if !ok {
			return "", fmt.Errorf("area not found: %s", a)
		}
		codes = append(codes, strconv.Itoa(c))
	}
	v.Add("sp", "s="+strings.Join(codes, ","))
	if minPrice != 0 {
		v.Add("pf", strconv.Itoa(minPrice))
	}
	return search + v.Encode(), nil
}

//
//func main() {
//	// Instantiate default collector
//	c := colly.NewCollector(colly.AllowURLRevisit())
//
//	// Rotate two socks5 proxies
//	rp, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:1337", "socks5://127.0.0.1:1338")
//	if err != nil {
//		log.Fatal(err)
//	}
//	c.SetProxyFunc(rp)
//
//	// Print the response
//	c.OnResponse(func(r *colly.Response) {
//		log.Printf("Proxy Address: %s\n", r.Request.ProxyURL)
//		log.Printf("%s\n", bytes.Replace(r.Body, []byte("\n"), nil, -1))
//	})
//
//	// Fetch httpbin.org/ip five times
//	for i := 0; i < 5; i++ {
//		err := c.Visit("https://httpbin.org/ip")
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
//}
