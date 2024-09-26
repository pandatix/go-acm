package api

import (
	"net/http"
	"net/url"

	"github.com/antchfx/htmlquery"
	"github.com/gorilla/schema"
)

type (
	SearchParams struct {
		Request   string `schema:"AllField"`
		StartPage *int   `schema:"startPage,omitempty"`
		PageSize  *int   `schema:"pageSize,omitempty"`
	}

	SearchResponse struct {
		References []Reference
	}

	Reference struct {
		Category   string
		Title      string
		PubDate    string
		Conference string
	}
)

func (client *ACMClient) Search(params *SearchParams, opts ...Option) (*SearchResponse, error) {
	// Build request
	req, _ := http.NewRequest(http.MethodGet, "/action/doSearch", nil)

	// Encode parameters and append them to the slice, if any
	if params != nil {
		values := url.Values{}
		if err := schema.NewEncoder().Encode(params, values); err != nil {
			return nil, err
		}
		req.URL.RawQuery = values.Encode()
	}

	// Issue request
	res, err := client.call(req, opts...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Parse document content to extract actual data
	doc, err := htmlquery.Parse(res.Body)
	if err != nil {
		return nil, err
	}
	lis := htmlquery.Find(doc, "//li[contains(@class, 'search__item') and contains(@class, 'issue-item-container')]")
	refs := make([]Reference, 0, len(lis))
	for _, li := range lis {
		refs = append(refs, Reference{
			Category:   htmlquery.InnerText(li.FirstChild.NextSibling.FirstChild.FirstChild),
			Title:      htmlquery.InnerText(li.FirstChild.NextSibling.FirstChild.NextSibling.FirstChild.FirstChild),
			PubDate:    htmlquery.InnerText(li.FirstChild.NextSibling.FirstChild.FirstChild.NextSibling),
			Conference: htmlquery.InnerText(li.FirstChild.NextSibling.FirstChild.NextSibling.FirstChild.FirstChild.NextSibling.NextSibling.FirstChild),
		})
	}

	return &SearchResponse{
		References: refs,
	}, nil
}
