package pkg

import (
	"strings"

	"github.com/gocolly/colly"
	cmap "github.com/orcaman/concurrent-map"
)

// LinkConfig config links
type LinkConfig struct {
	Host          string
	Destination   string
	ExternalLinks bool
	RespectTree   bool
	FilterIn      []string
	FilterOut     []string
}

// IsInHost returns true if is the same host
func (config LinkConfig) IsInHost(uri string) bool {
	// return
	return (strings.HasPrefix(uri, "/") && uri != "/") || (strings.Contains(uri, config.Host) && uri != config.Host) || config.ExternalLinks
}

// IsFile detects if the url is a file
func (config LinkConfig) IsFile(uri string) bool {
	split := strings.Split(uri, "/")
	last := split[len(split)-1]
	return strings.Contains(last, ".") || strings.Contains(last, "javascript:") || strings.Contains(last, "#")
}

// FetchLinks fetch links with crawler
func FetchLinks(config LinkConfig) []string {

	val := struct{}{}
	idx := cmap.New()
	printer := NewPrinter(config)
	idx.SetIfAbsent(config.Host, val)
	c := colly.NewCollector()
	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		url := e.Attr("href")

		if config.IsInHost(url) && !config.IsFile(url) && !idx.Has(url) {
			idx.SetIfAbsent(url, val)
			e.Request.Visit(url)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		url := r.URL.String()
		printer.queue <- url
	})

	c.Visit(config.Host)

	return printer.Start()
}

// Printer printer type for pdf renderer
type Printer struct {
	config LinkConfig
	queue  chan string
}

// NewPrinter creates a new Printer
func NewPrinter(config LinkConfig) *Printer {
	p := &Printer{
		config: config,
		queue:  make(chan string, 100000),
	}
	return p
}

// Start starts rendering pdf
func (p *Printer) Start() (links []string) {
	var value string
	links = make([]string, 0, len(p.queue))
	for {
		if len(p.queue) == 0 {
			break
		}
		value = <-p.queue
		if value == "" {
			break
		}
		links = append(links, value)
	}
	if p.config.RespectTree {
		tree := &LinkTree{}
		tree.AddLinks(links...)
		links = tree.GetLinks()
	}
	return
}

// LinkTree tree of links. It respects the link hierarchy
type LinkTree struct {
	Value  string
	Childs []*LinkTree
}

// AddLinks add links to a tree
func (t *LinkTree) AddLinks(links ...string) {
	if len(links) == 0 {
		return
	}
	// get the first one
	if t.Value == "" {
		t.Value = links[0]
		links = links[1:]
	}
	for _, l := range links {
		t.AddLink(l)
	}
}

// AddLink will add a link to itself as child or to its child if it is in the same root
// e.g. would add /a/b to /a but not /c to /a
func (t *LinkTree) AddLink(link string) {
	if t.Childs == nil {
		t.Childs = []*LinkTree{}
	}
	if !t.IsChild(link, false) {
		return
	}
	for _, c := range t.Childs {
		if c.IsChild(link, true) {
			c.AddLink(link)
			return
		}
	}
	t.Childs = append(t.Childs, &LinkTree{
		Value:  link,
		Childs: []*LinkTree{},
	})
}

// IsChild returns true if link is child of the node
func (t *LinkTree) IsChild(link string, checkChilds bool) (response bool) {
	if !strings.HasPrefix(link, t.Value) {
		return
	}
	response = true
	if !checkChilds {
		return
	}
	return
}

// GetLinks get links
func (t *LinkTree) GetLinks() (links []string) {
	links = make([]string, 0, 1+len(t.Childs))
	links = append(links, t.Value)
	for _, c := range t.Childs {
		links = append(links, c.GetLinks()...)
	}
	return
}
