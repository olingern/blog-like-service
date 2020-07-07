package actions

import (
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

func ConvertName(name string) string {
	underscored := strings.Join(strings.Fields(name), "_")

	return strings.ToLower(underscored)
}

func (a *Actions) GetAllPosts() []*firestore.DocumentSnapshot {
	allPosts, err := a.store.GetAllPosts()

	if err != nil {
		log.Error("[GetAllPosts] Error retreiving all posts")
	}

	return allPosts
}

func (a *Actions) Scrape(w http.ResponseWriter, r *http.Request) {

	c := colly.NewCollector()

	convertedLinkNames := []string{}

	// Find and visit all links
	c.OnHTML("h2", func(e *colly.HTMLElement) {
		log.Info("Found h2: ", e.Text)
		convertedLinkNames = append(convertedLinkNames, ConvertName(e.Text))
	})

	c.OnRequest(func(r *colly.Request) {
		log.Info("Visiting: ", r.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		log.Info("Finished: ", r.Request.URL)

		a.store.ProcessPosts(convertedLinkNames)
	})

	log.Info("Visiting ", a.siteURL)
	c.Visit(a.siteURL)
}
