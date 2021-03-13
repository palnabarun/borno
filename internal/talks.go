package internal

import (
	"os"
	"sort"
	"time"

	"gopkg.in/yaml.v2"
)

// Presenter describes a presenter
type Presenter struct {
	Name string `yaml:"name"`
	Link string `yaml:"link"`
}

// Link describes a hyperlink
type Link struct {
	Name   string `yaml:"name"`
	Target string `yaml:"target"`
}

// Talk stores metadata for a talk
type Talk struct {
	Link         string      `yaml:"link"`
	Title        string      `yaml:"title"`
	Date         time.Time   `yaml:"date"`
	Location     string      `yaml:"location"`
	SlideURL     string      `yaml:"slides"`
	VideoURL     string      `yaml:"video"`
	CoPresenters []Presenter `yaml:"copresenters"`
}

// TalkGroup stores Talk objects for a year
type TalkGroup struct {
	Year  int
	Talks []Talk
}

// ConfigOpts describes options for the config parser
type ConfigOpts struct {
	ConfigFile string
}

// BornoConfig describes the configuration file format for borno
type BornoConfig struct {
	Author    string `yaml:"author"`
	PageTitle string `yaml:"page_title"`
	Links     []Link `yaml:"links"`
	Talks     []Talk `yaml:"talks"`
}

// ParseTalksFromConfig parses the borno config file and returns Talks
func ParseTalksFromConfig(opts *ConfigOpts) (BornoConfig, error) {
	var config BornoConfig

	b, err := os.ReadFile(opts.ConfigFile)
	if err != nil {
		return config, err
	}

	if err := yaml.Unmarshal(b, &config); err != nil {
		return config, err
	}

	return config, nil
}

func groupByYear(talks []Talk) []TalkGroup {
	yearMap := make(map[int]bool)

	for _, t := range talks {
		yearMap[t.Date.Year()] = true
	}

	groups := make([]TalkGroup, 0)

	getTalksForYear := func(talks []Talk, year int) []Talk {
		filteredTalks := make([]Talk, 0)
		for _, t := range talks {
			if t.Date.Year() == year {
				filteredTalks = append(filteredTalks, t)
			}
		}

		sort.SliceStable(filteredTalks, func(i, j int) bool {
			return filteredTalks[i].Date.After(filteredTalks[j].Date)
		})
		return filteredTalks
	}

	for y := range yearMap {
		groups = append(groups, TalkGroup{Year: y, Talks: getTalksForYear(talks, y)})
	}

	sort.SliceStable(groups, func(i, j int) bool {
		return groups[i].Year > groups[j].Year
	})

	return groups
}