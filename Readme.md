# GoScraping

## Web scraping tool made for bug hunters
### This tool helps you to get data from HTML using CSS selectors.

## <b>Usage: </b>
```
➜  GoScraping -h
-u,  --url         Fetch HTML from URL
-s,  --std         Read HTML from standard input
-f,  --file        Read HTML file from disk
-q,  --query       CSS selector to be used
-a,  --attr        Only print nodes with specified attribute
-av  --attr-val    Used with -a to get nodes with specified value of attribute
-v,  --verbose     Verbose mode. (Prints file name in file mode)
```

## <b>Installing</b>
```bash
➜  go get github.com/cryonayes/GoScraping
```

## <b>Examples</b>
```bash
➜  GoScraping -u "https://www.w3schools.com/howto/howto_css_example_website.asp" -q "#main > div.w3-content > div > div:nth-child(2) > a"
```

Combining with fff output.
```bash
➜  find . -type f -name "*\.body" | xargs -n1 -I{} GoScraping -f {} -q "#main > div.content > div > a"
```

Filtering by attributes
```bash
➜  GoScraping -u "https://www.w3schools.com/howto/tryhow_css_example_website.htm" -q "body > div.row > div.main > div" -a "class" -av "fakeimg"
```
