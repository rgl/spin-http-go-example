package main

import (
	"fmt"
	"net/http"
	"strings"

	spinhttp "github.com/fermyon/spin/sdk/go/v2/http"
)

type indexTemplateData struct {
	SourceUrl  string
	Version    string
	Revision   string
	AuthorName string
	AuthorUrl  string
}

const indexTemplateText = `<!DOCTYPE html>
<html>
  <head>
    <title>Hello, World!</title>
    <link rel="icon" type="image/png" href="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAsklEQVR4nNXRSwrCQAwG4NFLeB9voIiINNl4BZPBVUUXVtyI4KbJrLyIl3Ljg4otFSlOcWbnvw1fkskY858pTAfFJmmadiMxH1G5AKVTDD5UmK/o5oMgD0rbaIxKGy9OcttvnSy0qzEIDxtF3pdFoeV3zJkXl9Mdz0DpUTXhrIHXP70ZxE5R6fbZBJVWQQcDJUDh+wuh8Ll1bV/Q8fi9ScxXlU1yHoHSJQrXmciiF41D8wTn54OWVxRsfQAAAABJRU5ErkJggg==" />
    <style>
      body {
        margin: 0;
        display: flex;
        align-items: center;
        justify-content: center;
        min-height: 100vh;
        text-align: center;
        font-size: 27px;
      }
      .version {
        color: #aaa;
        font-style: italic;
        font-size: 0.6em;
      }
      .version a {
        color: inherit;
        text-decoration: none;
      }
      .version a:hover {
        text-decoration: underline;
      }
    </style>
  </head>
  <body>
    <div>
      <p>Hello, World!</p>
      <p class="version">v{{.Version}}+{{.Revision}}<br/><a href="{{.SourceUrl}}">{{.SourceUrl}}</a><br>by <a href="{{.AuthorUrl}}">{{.AuthorName}}</a></p>
    </div>
  </body>
</html>
`

func init() {
	// NB we cannot use html/template with TinyGo. so we did a poor man's
	//    implementation.
	// see https://tinygo.org/docs/reference/lang-support/stdlib/#htmltemplate
	indexData := indexTemplateData{
		SourceUrl:  sourceUrl,
		Version:    version,
		Revision:   revision,
		AuthorName: authorName,
		AuthorUrl:  authorUrl,
	}
	indexHtml := indexTemplateText
	indexHtml = strings.ReplaceAll(indexHtml, "{{.SourceUrl}}", indexData.SourceUrl)
	indexHtml = strings.ReplaceAll(indexHtml, "{{.Version}}", indexData.Version)
	indexHtml = strings.ReplaceAll(indexHtml, "{{.Revision}}", indexData.Revision)
	indexHtml = strings.ReplaceAll(indexHtml, "{{.AuthorName}}", indexData.AuthorName)
	indexHtml = strings.ReplaceAll(indexHtml, "{{.AuthorUrl}}", indexData.AuthorUrl)

	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "Not Found.")
			return
		}
		if r.Method != "GET" && r.Method != "HEAD" {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintln(w, "Not Allowed.")
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, indexHtml)
	})
}

func main() {
	// this main definition is only required to compile; its not used at
	// runtime.
}
