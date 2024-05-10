package book_crawler

import (
	"context"
	"encoding/json"
	
	"temporal-crawler/internal/entity"
)

// BuildBookDocusaurusIndexActivity builds book index file: _category_.json
//
// Its schema is like:
//
//	{
//	   "label": "1 Cô-rinh-tô",
//	   "position": 1,
//	   "link": { "type": "generated-index" }
//	 }
func BuildBookDocusaurusIndexActivity(ctx context.Context, bookInfo entity.BookInfo) ([]byte, error) {
	category := docusaurusCategory{
		Label:    bookInfo.BookName,
		Position: bookInfo.BookNumber,
		Link: docusaurusLink{
			Type: "generated-index",
		},
	}
	
	return json.Marshal(category)
}

type docusaurusCategory struct {
	Label    string         `json:"label"`
	Position uint           `json:"position"`
	Link     docusaurusLink `json:"link"`
}

type docusaurusLink struct {
	Type string `json:"type"`
}
