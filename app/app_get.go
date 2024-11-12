package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sub-ui/setup"
	"time"
)

func GetAppUrl() {
	MuTime.Lock()
	defer MuTime.Unlock()

	now := time.Now()

	if LastExecutionTime.IsZero() || now.Sub(LastExecutionTime) >= 24*time.Hour {

		var appUrl string
		var err error
		githuaProxy := setup.ConfigData.App.GitHubProxy

		var apps []Properties
		for i := range setup.ConfigData.App.AppList {

			if setup.ConfigData.App.AppList[i].Url != "" {
				apps = append(apps, Properties{
					OnlyCopy: setup.ConfigData.App.AppList[i].OnlyCopy,
					Label:    setup.ConfigData.App.AppList[i].Label,
					Url:      setup.ConfigData.App.AppList[i].Url,
				})
				continue
			}

			user := setup.ConfigData.App.AppList[i].User
			repo := setup.ConfigData.App.AppList[i].Repository
			regStr := setup.ConfigData.App.AppList[i].Regexp
			appUrl, err = GetLatestAppUrl(user, repo, regStr)

			if err == nil {
				apps = append(apps, Properties{
					OnlyCopy: setup.ConfigData.App.AppList[i].OnlyCopy,
					Label:    setup.ConfigData.App.AppList[i].Label,
					Url:      githuaProxy + appUrl,
				})
			}
		}

		AppsData = apps

		LastExecutionTime = now
	}
}

func GetLatestAppUrl(owner, repo, regStr string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API request failed with status code %d", resp.StatusCode)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	var matchBool bool
	var appUrl string

	for _, asset := range release.Assets {
		matchBool, _ = regexp.MatchString(regStr, asset.URL)
		if matchBool {
			appUrl = asset.URL
			break
		}
	}

	return appUrl, nil
}
