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
	// printer.queue <- config.Host
	idx.SetIfAbsent(config.Host, val)
	c := colly.NewCollector()
	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		url := e.Attr("href")

		if config.IsInHost(url) && !config.IsFile(url) && !idx.Has(url) {
			// fmt.Println("will visit ", url)
			idx.SetIfAbsent(url, val)
			e.Request.Visit(url)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting", r.URL)
		url := r.URL.String()
		printer.queue <- url
	})

	c.Visit(config.Host)

	// fmt.Println("------------------")
	// for i, r := range idx.Items() {
	// fmt.Println(i, "-", r)
	// }

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
			// fmt.Println("channel closed..")
			return
		}
		value = <-p.queue
		if value == "" {
			// fmt.Println("channel closed..")
			return
		}
		// file := strings.Replace(strings.TrimRight(strings.TrimLeft(value, p.config.Host), "/"), "/", "-", -1)
		// err := exec.Command(
		// 	"wkhtmltopdf",
		// 	value,
		// 	"--disable-external-links",
		// 	"--disable-internal-links",
		// 	"--orientation", "Portrait",
		// 	"--page-size", "A4",
		// 	p.config.Destination+"/"+file+".pdf",
		// ).Run()
		// if err != nil {
		// 	fmt.Println("err saving file: "+p.config.Destination+"/"+file+".pdf : ", err)
		// } else {
		// 	fmt.Println("pdf generated: " + p.config.Destination + "/" + file + ".pdf")
		// }
		links = append(links, value)
	}
	return
}

// // Crawler basic crawler
// type Crawler struct {
// 	fetch  *fetchbot.Fetcher
// 	queue  *fetchbot.Queue
// 	colly.
// 	config LinkConfig
// }

// // NewCrawler constructor
// func NewCrawler(config LinkConfig) *Crawler {
// 	crawler := &Crawler{config: config}
// 	f := fetchbot.New(fetchbot.HandlerFunc(crawler.Handle))
// 	crawler.fetch = f
// 	crawler.queue = f.Start()
// 	return crawler
// }

// // Handle receiving handle function
// func (c *Crawler) Handle(ctx *fetchbot.Context, res *http.Response, err error) {
// 	requestURL := ctx.Cmd.URL().RequestURI()
// 	if c.config.IsInHost(requestURL) && !c.config.IsFile(requestURL) {
// 		// sdad
// 		// sdasd
// 		// c.Add(requestURL)

// 	}
// }

// func (c *Crawler) Add(url string) {

// }

// // Execute run craw function
// func (c *Crawler) Execute() {

// }
