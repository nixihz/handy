package notionservice

import (
	"context"

	"github.com/kjk/notion"
	"github.com/tidwall/pretty"
)

// shows how to use https://developers.notion.com/reference/get-page API

// shows the info about:
//   regular page https://www.notion.so/Test-pages-for-notionapi-0367c2db381a4f8b9ce360f388a6b2e3
//   page in a database https://www.notion.so/A-row-that-is-not-empty-page-e56b74a6398a43848137cca2a0de20b2
// or the page given with -id argument (in which case also needs )

func getIndent(n int) string {
	s := ""
	for n > 0 {
		n -= 1
		s += "  "
	}
	return s
}

func ShowRichText(indent int, name string, richText []notion.RichText) {
	s := getIndent(indent)
	// TODO: better implementation
	if name != "" {
		logf("%s%s: %v\n", s, name, richText)
		return
	}
	logf("%s%v\n", s, richText)
}

func ppJSON(d []byte) {
	res := pretty.Pretty(d)
	logf("pretty printed JSON:\n%s\n", res)
}

func ShowPageInfo(page *notion.Page) {
	logf("ShowPageInfo:\n")
	logf("  ID: '%s'\n", page.ID)
	logf("  CreatedTime: '%s'\n", page.CreatedTime)
	logf("  LastEditedTime: '%s'\n", page.LastEditedTime)
	if page.Parent.PageID != nil {
		logf("  Parent: page with ID '%s'\n", *page.Parent.PageID)
	} else if page.Parent.DatabaseID != nil {
		logf("  Parent: database with ID '%s'\n", *page.Parent.DatabaseID)
	} else {
		logf("both page.Parent.PageID or page.Parent.DatabaseID are nil")
	}
	logf("  Archived: %v\n", page.Archived)
	switch prop := page.Properties.(type) {
	case notion.PageProperties:
		logf("  page properties:\n")
		ShowRichText(2, "Title", prop.Title.Title)
	case notion.DatabasePageProperties:
		logf("  database properties (NYI):\n")
	}
}

func GetPageInfo(apiKey string, pageID string) {
	logf("GetPageInfo: pageID='%s'\n", pageID)

	c := GetClient(apiKey)
	ctx := context.Background()
	page, err := c.GetPage(ctx, pageID)
	if err != nil {
		logf("GetPage() failed with: '%s'\n", err)
		logf("page.RawJSON: '%s'\n", page.RawJSON)
		ppJSON(page.RawJSON)
		return
	}
	ShowPageInfo(page)
}
