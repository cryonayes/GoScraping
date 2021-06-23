package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Config struct {
	fileName      string
	fromURL       string
	jsPathQuery   string
	attrName      string
	attrValue     string
	readFromSTDIO bool
	verboseFlag   bool
}

func readFromFile(config *Config) {
	file, err := os.Open(config.fileName)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error while opening file!\n")
		return
	}
	readAll := bufio.NewReader(file)
	parseData(config, readAll)
}

func readFromStdIn(config *Config) {
	stdin := bufio.NewReader(os.Stdin)
	parseData(config, stdin)
}

func readFromURL(config *Config) {
	response, err := http.Get(config.fromURL)
	if err != nil || response == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Can't connect to host!")
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	reader := bufio.NewReader(response.Body)
	parseData(config, reader)
}

func parseData(config *Config, reader *bufio.Reader) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error while parsing HTML")
		return
	}

	doc.Find(config.jsPathQuery).Each(func(i int, selection *goquery.Selection) {
		if config.verboseFlag && !config.readFromSTDIO {
			_, _ = fmt.Fprint(os.Stdout, config.fileName+":")
		}
		if config.attrName == "" {
			_, _ = fmt.Fprintln(os.Stdout, selection.Text())
			return
		}

		if value, exist := selection.Attr(config.attrName); exist {
			if config.attrValue == "" {
				_, _ = fmt.Fprintln(os.Stdout, selection.Text())
				return
			}
			if strings.Contains(value, config.attrValue) {
				_, _ = fmt.Fprintln(os.Stdout, selection.Text())
			}
		}
	})
}

func main() {

	flag.Usage = func() {
		help := []string{
			"-u,  --url         Fetch HTML from URL",
			"-s,  --std         Read HTML from standard input",
			"-f,  --file        Read HTML file from disk",
			"-q,  --query       CSS selector to be used",
			"-a,  --attr        Only print nodes with specified attribute",
			"-av  --attr-val    Used with -a to get nodes with specified value of attribute",
			"-v,  --verbose     Verbose mode. (Prints file name in file mode)",
			"",
		}
		_, _ = fmt.Fprintf(os.Stdout, strings.Join(help, "\n"))
	}

	mConfig := Config{}
	flag.BoolVar(&mConfig.readFromSTDIO, "s", false, "")
	flag.BoolVar(&mConfig.readFromSTDIO, "std", false, "")

	flag.BoolVar(&mConfig.verboseFlag, "v", false, "")
	flag.BoolVar(&mConfig.verboseFlag, "verbose", false, "")

	flag.StringVar(&mConfig.fileName, "f", "", "")
	flag.StringVar(&mConfig.fileName, "file", "", "")

	flag.StringVar(&mConfig.fromURL, "u", "", "")
	flag.StringVar(&mConfig.fromURL, "url", "", "")

	flag.StringVar(&mConfig.jsPathQuery, "q", "", "")
	flag.StringVar(&mConfig.jsPathQuery, "query", "", "")

	flag.StringVar(&mConfig.attrName, "a", "", "")
	flag.StringVar(&mConfig.attrName, "attr", "", "")

	flag.StringVar(&mConfig.attrValue, "av", "", "")
	flag.StringVar(&mConfig.attrValue, "attr-val", "", "")

	flag.Parse()

	if !mConfig.readFromSTDIO && mConfig.fileName == "" && mConfig.fromURL == "" {
		_, _ = fmt.Fprintln(os.Stderr, "Please provide file or URL")
		return
	}

	if mConfig.readFromSTDIO {
		readFromStdIn(&mConfig)
	} else if mConfig.fromURL != "" {
		readFromURL(&mConfig)
	} else {
		readFromFile(&mConfig)
	}
}
