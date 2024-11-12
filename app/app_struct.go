package app

import (
	"sync"
	"time"
)

var (
	LastExecutionTime time.Time
	MuTime            sync.Mutex
	AppsData          []Properties
)

type Properties struct {
	OnlyCopy bool
	Label    string
	Url      string
}

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name string `json:"name"`
	URL  string `json:"browser_download_url"`
}
