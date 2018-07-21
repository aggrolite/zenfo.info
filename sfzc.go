package zenfo

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

// Sfzc crawls sfzc.org, satisfies Worker interface
type Sfzc struct {
	venueMap map[string]*Venue
	client   *Client
}

type eventJSON struct {
	Name      string    `json:"Name"`
	Blurb     string    `json:"Introduction"`
	Desc      string    `json:"Description"`
	Start     time.Time `json:"DateStart"`
	End       time.Time `json:"DateEnd"`
	URL       string    `json:"Link"`
	ContentID int       `json:"ContentId"`
	Location  string    `json:"Location"`
	Type      string    `json:"EventType"`
}

func html2text(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", errors.Wrap(err, "goquery.NewDocument()")
	}
	out := doc.Text()
	return strings.TrimSpace(out), nil
}

// Init sets HTTP client and defines internal venue map
func (s *Sfzc) Init(client *Client) error {

	s.client = client
	s.venueMap = make(map[string]*Venue)

	s.venueMap["City Center"] = &Venue{
		Name:    "San Francisco Zen Center",
		Addr:    "300 Page St, San Francisco, CA 94102",
		Phone:   "+1 (415) 863-3136",
		Email:   "ccoffice@sfzc.org",
		Lat:     37.773789,
		Lng:     -122.426153,
		Website: "https://sfzc.org",
	}
	s.venueMap["Tassajara"] = &Venue{
		Name:    "Tassajra Zen Mountain Center",
		Addr:    "39171 Tassajara Road, Carmel Valley, CA 93924",
		Phone:   "+1 (831) 659-2229",
		Email:   "rezoffice@sfzc.org",
		Lat:     36.234131,
		Lng:     -121.550031,
		Website: "http://sfzc.org/tassajara",
	}
	s.venueMap["Green Gulch"] = &Venue{
		Name:    "Green Gulch Farm Zen Center",
		Addr:    "1601 Shoreline Highway, Muir Beach, CA 94965",
		Phone:   "+1 (415) 383-3134",
		Email:   "ggfoffice@sfzc.org",
		Lat:     37.865967,
		Lng:     -122.563911,
		Website: "http://sfzc.org/green-gulch",
	}

	return nil
}

// Desc returns description for website crawled
func (s *Sfzc) Desc() string {
	return "San Francisco Zen Center (sfzc.org)"
}

// Events hits sfcz JSON API and returns slice of Event types
func (s *Sfzc) Events() ([]*Event, error) {
	resp, err := s.client.Get("http://sfzc.org/api/eventsapi/allevents")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var events []eventJSON
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, err
	}

	var final []*Event
	for _, e := range events {

		blurb, err := html2text(e.Blurb)
		if err != nil {
			return nil, err
		}
		desc, err := html2text(e.Desc)
		if err != nil {
			return nil, err
		}
		e.Blurb = blurb
		e.Desc = desc
		if e.End.Before(e.Start) {
			e.End = e.Start
		}
		var link string
		if !strings.HasPrefix(e.URL, "http") {
			link = "http://sfzc.org"
		}
		u := fmt.Sprintf("%s%s", link, e.URL)

		venue, ok := s.venueMap[e.Location]
		if !ok {
			return nil, fmt.Errorf("Failed to match venue for '%s'", e.Location)
		}

		finalEvent := &Event{
			Name:  e.Name,
			Blurb: e.Blurb,
			Desc:  e.Desc,
			Start: e.Start,
			End:   e.End,
			URL:   u,
			Venue: venue,
		}

		final = append(final, finalEvent)
	}

	return final, nil
}
