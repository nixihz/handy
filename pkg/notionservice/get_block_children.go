package notionservice

import (
	"context"
	"fmt"

	"github.com/kjk/notion"
)

func ShowBlock(b *notion.Block) {
	logf("  %s %s, has_children: %v\n", b.Type, b.ID, b.HasChildren)
	switch b.Type {
	case notion.BlockTypeParagraph:
		logf(" %v\n", b.Paragraph.Text)
	case notion.BlockTypeHeading1:
		logf(" %v\n", b.Heading1.Text)
	case notion.BlockTypeHeading2:
		logf(" %v\n", b.Heading2.Text)
	case notion.BlockTypeHeading3:
		logf(" %v\n", b.Heading3.Text)
	case notion.BlockTypeBulletedListItem:
		logf(" %v\n", b.BulletedListItem.Text)
	case notion.BlockTypeNumberedListItem:
		logf(" %v\n", b.NumberedListItem.Text)
	case notion.BlockTypeToDo:
		logf(" %v\n", b.ToDo.Text)
	case notion.BlockTypeToggle:
		logf(" %v\n", b.Toggle.Text)
	case notion.BlockTypeChildPage:
	case notion.BlockTypeUnsupported:
	}
}

func ShowBlockChildren(bcr *notion.BlockChildrenResponse) {
	logf("ShowBlockChildren:\n")
	logf("  hasMore: %v\n", bcr.HasMore)
	if bcr.NextCursor != "" {
		logf("  nextCursor: %v\n", bcr.NextCursor)
	}
	logf("  %d children:\n", len(bcr.Results))
	for _, b := range bcr.Results {
		ShowBlock(&b)
	}
}

func GetBlockChildren(apiKey string, blockID string) {
	fmt.Printf("GetBlockChildren: blockID='%s'\n", blockID)

	c := GetClient(apiKey)
	ctx := context.Background()
	rsp, err := c.GetBlockChildren(ctx, blockID, nil)
	if err != nil {
		logf("GetBlockChildren() failed with '%s'\n", err)
		logf("page.RawJSON: '%s'\n", rsp.RawJSON)
		ppJSON(rsp.RawJSON)
		return
	}
	ShowBlockChildren(rsp)
}
