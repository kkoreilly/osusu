package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector(colly.AllowedDomains("www.allrecipes.com"))

	c.OnXML("//sitemapindex/sitemap/loc", func(e *colly.XMLElement) {
		e.Request.Visit(e.Text)
	})

	total := 0
	c.OnXML("//urlset/url/loc", func(e *colly.XMLElement) {
		fmt.Println(e.Text)
		total++
		e.Request.Visit(e.Text)
	})

	c.Visit("https://www.allrecipes.com/sitemap.xml")
	fmt.Println(total)
}
