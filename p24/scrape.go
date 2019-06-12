package p24

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/corverroos/goldprop/models"
	"github.com/go-stack/stack"
	"github.com/gocolly/colly"
	"github.com/jinzhu/gorm"
)

type description struct {
	URL         string
	Description string
}

type feature struct {
	URL   string
	Key   string
	Value string
}

func Scrape(db *gorm.DB, areas ...string) error {
	c := colly.NewCollector(
		colly.AllowedDomains("www.property24.com"),
		colly.CacheDir("/tmp/colly"),
	)
	maybeLog(c.Limit(&colly.LimitRule{
		RandomDelay: 5 * time.Second,
		Parallelism: 2,
		DomainGlob:  "www.property24.com.*",
	}))

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
	c.OnHTML(".p24_listing .p24_dPL", func(e *colly.HTMLElement) {
		ch <- description{
			URL:         e.Request.URL.String(),
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
			log.Printf("unexpectednext page: %v", e)
		}
		maybeLog(c.Visit(link))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	var wg sync.WaitGroup
	wg.Add(1)
	var consumeErr error
	go func() {
		consumeErr = consume(db, ch)
		wg.Done()
	}()

	search, err := makeSearchUrl(0, areas...)
	if err != nil {
		return err
	}

	err = c.Visit(search)
	if err != nil {
		return err
	}
	close(ch)
	wg.Wait()
	return consumeErr
}

func consume(db *gorm.DB, ch <-chan interface{}) error {
	for i := range ch {
		switch t := i.(type) {
		case tile:
			var l models.Listing
			db.Where("url=?", t.URL).First(&l)

			var err error
			l.Price, err = t.GetPrice()
			if err != nil {
				return err
			}
			l.Bathrooms, err = t.GetBathrooms()
			if err != nil {
				return err
			}
			l.Bedrooms, err = t.GetBedrooms()
			if err != nil {
				return err
			}
			l.Garages, err = t.GetGarages()
			if err != nil {
				return err
			}
			l.FloorSize, err = t.GetSize()
			if err != nil {
				return err
			}
			l.SiteID, err = getSiteID(t.URL)
			if err != nil {
				return err
			}
			l.Location = t.Location
			l.Address = t.Address
			l.URL = t.URL
			l.Site = models.SiteP24
			l.PropertyType = models.PropertyTypeApartment
			l.Listing = models.ListingTypeSale

			if l.ID == 0 {
				db.Create(&l)
			} else {
				db.Model(&l).Update(&l)
			}
		case feature:
			var l models.Listing
			if db.Where("url=?", t.URL).First(&l).RecordNotFound() {
				return fmt.Errorf("Cannot find listing for feature: %v", t.URL)
			}
			var f models.Features
			db.Where("listing_id = ? AND `key` = ?", l.ID, t.Key).First(&f)
			if f.ID != 0 && f.Value != t.Value {
				f.Value = t.Value
				db.Update(&f)
			} else if f.ID == 0 {
				f.ListingID = l.ID
				f.Key = t.Key
				f.Value = t.Value
				db.Create(&f)
			}
		case description:
			var l models.Listing
			if db.Where("url=?", t.URL).First(&l).RecordNotFound() {
				return fmt.Errorf("Cannot find listing: %v", t.URL)
			}
			l.Description = t.Description
			db.Model(&l).Update(&l)
		}
	}
	return nil
}

func getSiteID(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	return path.Base(u.Path), nil
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
	if len(arealist) == 0 {
		return "", errors.New("no areas specified")
	}
	v := url.Values{}
	var codes []string
	for _, a := range arealist {
		c, ok := Areas[strings.ToLower(a)]
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

func (t tile) GetPrice() (int, error) {
	if t.Price == "POA" {
		return 0, nil
	}
	p := strings.Replace(t.Price, "R", "", 1)
	p = strings.Replace(p, " ", "", -1)
	p = strings.TrimSpace(p)
	i, err := strconv.Atoi(p)
	if err != nil {
		return 0, fmt.Errorf("unexpected Price: %s, %s", t.Price, t.URL)
	}
	return i, nil
}

func (t tile) GetBedrooms() (int, error) {
	if t.Bedrooms == "" {
		return 0, nil
	}
	i, err := strconv.Atoi(t.Bedrooms)
	if err != nil {
		return 0, fmt.Errorf("unexpected Bedrooms: %s, %s", t.Bedrooms, t.URL)
	}
	return i, nil
}

func (t tile) GetGarages() (int, error) {
	if t.Garages == "" {
		return 0, nil
	}
	i, err := strconv.Atoi(t.Garages)
	if err != nil {
		return 0, fmt.Errorf("unexpected Garages: %s, %s", t.Garages, t.URL)
	}
	return i, nil
}

func (t tile) GetBathrooms() (int, error) {
	if t.Bathrooms == "" {
		return 0, nil
	}
	i, err := strconv.Atoi(t.Bathrooms)
	if err != nil {
		return 0, fmt.Errorf("unexpected Bathrooms: %s, %s", t.Bathrooms, t.URL)
	}
	return i, nil
}

func (t tile) GetSize() (int, error) {
	if t.Size == "" {
		return 0, nil
	}
	split := strings.Split(t.Size, " ")
	if len(split) != 2 {
		return 0, fmt.Errorf("unexpected Size: %s, %s", t.Size, t.URL)
	}
	i, err := strconv.Atoi(split[0])
	if err != nil {
		return 0, fmt.Errorf("unexpected Size: %s, %s", t.Size, t.URL)
	}
	return i, nil
}
