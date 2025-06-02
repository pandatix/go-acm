package api

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type (
	SearchParams struct {
		Request    string `schema:"AllField"`
		StartPage  *int   `schema:"startPage,omitempty"`
		PageSize   *int   `schema:"pageSize,omitempty"`
		BeforeYear *int   `schema:"BeforeYear,omitempty"`
		AfterYear  *int   `schema:"AfterYear,omitempty"`
	}

	SearchResponse struct {
		Results    int
		References []*Reference
	}

	Reference struct {
		Category   string
		Title      string
		Abstract   string
		DOI        *string
		PubDate    string
		Conference *string
		OpenAccess bool
		FreeAccess bool
		Metrics    ReferenceMetrics
	}

	ReferenceMetrics struct {
		Citations      int
		TotalDownloads *int
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
	defer func() {
		_ = res.Body.Close()
	}()

	// Parse document content to extract actual data
	doc, err := htmlquery.Parse(res.Body)
	if err != nil {
		return nil, err
	}
	// => Total results
	totals := htmlquery.Find(doc, "//span[contains(@class, 'hitsLength')]")
	if len(totals) != 1 {
		return nil, errors.New("could not find total results")
	}
	total, err := acmAtoi(htmlquery.InnerText(totals[0]))
	if err != nil {
		return nil, err
	}
	// => References
	lis := htmlquery.Find(doc, "//li[contains(@class, 'search__item') and contains(@class, 'issue-item-container')]")
	refs := make([]*Reference, 0, len(lis))
	var merr error
	for lid, li := range lis {
		ref := &Reference{}
		var rmerr error

		// Category
		cats := htmlquery.Find(li, "//div[contains(@class, 'issue-heading')]")
		switch len(cats) {
		case 0:
			rmerr = multierr.Append(rmerr, errors.New("category not found"))
		case 1:
			ref.Category = htmlquery.InnerText(cats[0])
		default:
			rmerr = multierr.Append(rmerr, errors.New("too many categories match"))
		}

		// PubDate
		pubDates := htmlquery.Find(li, "//div[contains(@class, 'bookPubDate')]")
		switch len(pubDates) {
		case 0:
			rmerr = multierr.Append(rmerr, errors.New("publication date not found"))
		case 1:
			ref.PubDate = htmlquery.InnerText(pubDates[0])
		default:
			rmerr = multierr.Append(rmerr, errors.New("too many publication dates match"))
		}

		// Title
		titles := htmlquery.Find(li, "//h5[contains(@class, 'issue-item__title')]")
		switch len(titles) {
		case 0:
			rmerr = multierr.Append(rmerr, errors.New("title not found"))
		case 1:
			ref.Title = htmlquery.InnerText(titles[0])
		default:
			rmerr = multierr.Append(rmerr, errors.New("too many titles match"))
		}

		// Conference
		conferences := htmlquery.Find(li, "//span[contains(@class, 'epub-section__title')]")
		switch len(conferences) {
		case 0:
			// Is okay (e.g., doctoral thesis are not published at a conference)
		case 1:
			conf := htmlquery.InnerText(conferences[0])
			ref.Conference = &conf
		default:
			rmerr = multierr.Append(rmerr, errors.New("too many conferences match"))
		}

		// OpenAccess
		openAccesses := htmlquery.Find(li, "//div[contains(@class, 'open-access')]")
		switch len(openAccesses) {
		case 0:
			ref.OpenAccess = false
		case 1:
			ref.OpenAccess = true
		default:
			rmerr = multierr.Append(rmerr, errors.New("too many open accesses match"))
		}

		// FreeAccess
		freeAccesses := htmlquery.Find(li, "//div[contains(@class, 'free-access')]")
		switch len(freeAccesses) {
		case 0:
			ref.FreeAccess = false
		case 1:
			ref.FreeAccess = true
		default:
			rmerr = multierr.Append(rmerr, errors.New("too many free accesses match"))
		}

		// TODO authors
		// -> Truncated if too many, not in the list by default -> require UI interaction ? Where does the data comes from ??

		// DOI
		dois := htmlquery.Find(li, "//a[contains(@class, 'issue-item__doi')]")
		switch len(dois) {
		case 0:
			// Not that bad, is not mandatory
		case 1:
			doi := strings.TrimPrefix("https://doi.org/", htmlquery.InnerText(dois[0]))
			ref.DOI = &doi
		default:
			rmerr = multierr.Append(rmerr, errors.New("too many dois match"))
		}

		// Abstract
		abstracts := htmlquery.Find(li, "//div[contains(@class, 'issue-item__abstract')]")
		switch len(abstracts) {
		case 0:
			rmerr = multierr.Append(rmerr, errors.New("abstract not found"))
		case 1:
			ref.Abstract = htmlquery.InnerText(abstracts[0])
		default:
			rmerr = multierr.Append(rmerr, errors.New("too many abstract match"))
		}

		// Metrics: citations
		citations := htmlquery.Find(li, "//span[contains(@class, 'citation')]")
		switch len(citations) {
		case 0:
			rmerr = multierr.Append(rmerr, errors.New("citation not found"))
		case 1:
			ref.Metrics.Citations, err = acmAtoi(htmlquery.InnerText(citations[0]))
			if err != nil {
				rmerr = multierr.Append(rmerr, err)
			}
		default:
			rmerr = multierr.Append(rmerr, errors.New("too many citations match"))
		}
		// Metrics: downloads
		downloads := htmlquery.Find(li, "//span[contains(@class, 'metric')]")
		switch len(downloads) {
		case 0:
			// Not that bad, is not mandatory
		case 1:
			downloads, err := acmAtoi(htmlquery.InnerText(downloads[0]))
			if err != nil {
				rmerr = multierr.Append(rmerr, err)
			}
			ref.Metrics.TotalDownloads = &downloads
		default:
			rmerr = multierr.Append(rmerr, errors.New("too many downloads match"))
		}

		if rmerr != nil {
			merr = multierr.Append(merr, errors.Wrapf(rmerr, "reference %d", lid+1))
		} else {
			refs = append(refs, ref)
		}
	}
	if merr != nil {
		return nil, merr
	}

	return &SearchResponse{
		Results:    total,
		References: refs,
	}, nil
}

func acmAtoi(str string) (int, error) {
	str = strings.ReplaceAll(str, ",", "")
	str = strings.TrimSpace(str)
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0, errors.Wrapf(err, "cannot extract total results from %s", str)
	}
	return i, nil
}
