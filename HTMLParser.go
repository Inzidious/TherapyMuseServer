package main

import (
	//"encoding/json"
	"fmt"

	"github.com/gocolly/colly"
)

func (s *HTMLParser) BuildBody(h *colly.HTMLElement) {
	dom := h.DOM

	secondDom := dom.FindNodes(dom.Children().Nodes[0])
	secondDom.Remove()

	thirdDom := dom.Find("div.d-xl-none")
	thirdDom.Remove()

	fourthDom := dom.Find("div.insertArea")
	fourthDom.Remove()

	fmt.Println("body url: ", h.Request.URL)

	s.pageNodes[h.Request.URL.String()].Body = dom.Text()
	//s.pageNodes[]

	/*string, error := dom.Html()

	if error == nil {
		println("dom Text after removed: ", string)
	} else {
		println("Error: ", error.Error())
	}*/
}

func (s *HTMLParser) BuildNode(h *colly.HTMLElement) {
	dom := h.DOM
	childNodes := dom.Children().Nodes

	secondDom := dom.FindNodes(childNodes[1])
	thirdDom := secondDom.FindNodes(secondDom.Children().Nodes[0])

	fourthDom := dom.FindNodes(childNodes[3])

	href, exists := thirdDom.Attr("href")

	if exists && s.hit < 2 {
		s.hit++
		realHref := "https://www.psychologytoday.com" + href
		//fmt.Println("node href: ", realHref)
		s.pageNodes[realHref] = &PageNode{Title: secondDom.Text(), Topic: fourthDom.Text()}

		c := colly.NewCollector()

		c.OnRequest(func(r *colly.Request) {
			r.Headers.Set("Accept-Language", "en")
			fmt.Printf("Visiting %s\n", r.URL)
		})

		c.OnHTML("div.field-name-body", s.BuildBody)
		c.Visit(realHref)
	}

	//fmt.Println(s.pageNodes[href].title)
	//fmt.Println("href:", href)
	//fmt.Println("node:", node.Data)

	//if exists {
	//fmt.Println("exists:", exists)
	//}

	//fmt.Println(h.Text)
}

func (s *HTMLParser) Collect() {

	s.c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "en")
		fmt.Printf("Visiting %s\n", r.URL)
	})

	s.c.OnHTML("div.teaser-lg__details", s.BuildNode)
	s.c.Visit(s.url)

	/*for key, element := range s.pageNodes {
		fmt.Println("Key:", key, "=>", "Element:", element.body)
	}*/

}

func (s *HTMLParser) InitCollector() {
	s.pageNodes = make(map[string]*PageNode)
	s.hit = 0
	s.c = colly.NewCollector()
}
