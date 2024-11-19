package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sub-ui/app"
	"sub-ui/read"
	"sub-ui/setup"
	"time"
)

func downloadFile(url string, filepath string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func DownloadRules() {

	for i := range setup.ConfigData.Download.RuleList {
		url := setup.ConfigData.Download.RuleList[i].Url
		filepath := setup.ConfigData.Download.Folder + "/rules/" + setup.ConfigData.Download.RuleList[i].Name
		downloadFile(url, filepath)
	}

}

func DownloadApp() {
	isUpdate := false

	for i := range setup.ConfigData.Download.GithubList {
		if getGithubUrl(&setup.ConfigData.Download.GithubList[i]) {
			isUpdate = true
		}
	}

	if isUpdate {
		setup.SavedConfig()
	}
}

func getGithubUrl(github *setup.Github) bool {
	user := github.User
	repo := github.Repository
	regStr := github.Regexp
	newUrl, err := app.GetLatestAppUrl(user, repo, regStr)

	if err != nil || newUrl == "" {
		return false
	}

	url := github.Url

	filepath := setup.ConfigData.Download.Folder + "/app/" + github.Name

	if url == newUrl && read.CheckExistence(filepath) != "nil" {
		return false
	}

	github.Url = newUrl

	fmt.Println(newUrl)
	fmt.Println(filepath)

	//downloadFile(newUrl, filepath)
	//downloadFile(url, filepath)
	return true

}

func downloadRuleTicker(day int, hour int, url string, filepath string) {

	WaitStart(hour)

	ticker := time.NewTicker(time.Duration(day) * 24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			downloadFile(url, filepath)

			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location())

			nextEnd := next.Add(20 * time.Minute)
			if !now.After(next) || !now.Before(nextEnd) {
				go downloadRuleTicker(day, hour, url, filepath)
				return
			}
		}
	}
}

func RulesScheduledTasks() {

	hour := setup.ConfigData.Download.StartTime

	if hour < 0 || hour > 23 {
		hour = 1
	}

	for i := range setup.ConfigData.Download.RuleList {
		day := 6

		if setup.ConfigData.Download.RuleList[i].UpdateInterval > 0 &&
			setup.ConfigData.Download.RuleList[i].UpdateInterval < 31 {
			day = setup.ConfigData.Download.RuleList[i].UpdateInterval - 1
		}
		url := setup.ConfigData.Download.RuleList[i].Url
		filepath := setup.ConfigData.Download.Folder + "/rules/" + setup.ConfigData.Download.RuleList[i].Name

		go downloadRuleTicker(day, hour, url, filepath)

	}

}

func AppScheduledTasks() {

	hour := setup.ConfigData.Download.StartTime

	if hour < 0 || hour > 23 {
		hour = 1
	}

	WaitStart(hour)

	day := setup.ConfigData.Download.AppUpdateInterval

	if day < 7 || day > 90 {
		day = 14
	} else {
		day -= 1
	}

	ticker := time.NewTicker(time.Duration(day) * 24 * time.Hour)
	defer ticker.Stop()

	for {

		isUpdate := false
		select {
		case <-ticker.C:
			for i := range setup.ConfigData.Download.GithubList {
				if getGithubUrl(&setup.ConfigData.Download.GithubList[i]) {
					isUpdate = true
				}
			}
			if isUpdate {
				setup.SavedConfig()
			}
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location())

			nextEnd := next.Add(20 * time.Minute)
			if !now.After(next) || !now.Before(nextEnd) {
				go AppScheduledTasks()
				return
			}
		}

	}

}

func WaitStart(hour int) {

	now := time.Now()

	next := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location())

	if now.After(next) {
		next = next.Add(24 * time.Hour)
	}

	duration := time.Until(next)

	//fmt.Printf("等待时间: %v\n", duration)

	time.Sleep(duration)

}
