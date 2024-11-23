package serve

import (
	"html/template"
	"net/http"
	"sub-ui/app"
	"sub-ui/setup"
	"sub-ui/users"
	"time"
)

func (s Server) login(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie(setup.CookieName)
	if err == nil &&
		cookie.Value == setup.CookieValue {
		http.Redirect(w, r, setup.ConfigData.Server.Home.Url+"/sub-ui", http.StatusFound)
		return
	}

	error := ""
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == setup.ConfigData.Server.Home.User && password == setup.ConfigData.Server.Home.Password {

			http.SetCookie(w, &http.Cookie{
				Name:    setup.CookieName,
				Value:   setup.CookieValue,
				Path:    setup.ConfigData.Server.Home.Url,
				Expires: time.Now().Add(time.Duration(setup.CookieDay) * 24 * time.Hour),
			})
			http.Redirect(w, r, setup.ConfigData.Server.Home.Url+"/sub-ui", http.StatusFound)
			return
		} else {
			error = "用户或密码不正确!"
		}
	}

	w.Header().Set("Content-Type", "text/html")

	tmpl, err := template.ParseFiles("html/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	variables := struct {
		Login string
		Error string
	}{
		Login: setup.ConfigData.Server.Home.Url + "/login",
		Error: error,
	}

	err = tmpl.Execute(w, variables)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (s Server) logout(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		http.SetCookie(w, &http.Cookie{
			Name:   setup.CookieName,
			Value:  "",
			Path:   setup.ConfigData.Server.Home.Url,
			MaxAge: -1,
		})
		http.Redirect(w, r, setup.ConfigData.Server.Home.Url+"/login", http.StatusFound)
	} else {
		http.Redirect(w, r, setup.ConfigData.Server.Home.Url+"/login", http.StatusSeeOther)
	}
}

func (s Server) home(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie(setup.CookieName)
	if err != nil ||
		cookie.Value != setup.CookieValue {
		http.Redirect(w, r, setup.ConfigData.Server.Home.Url+"/login", http.StatusSeeOther)
		return
	}

	var titleStr string

	if setup.ConfigData.Server.Home.Title != "" {
		titleStr = `<h1>` + setup.ConfigData.Server.Home.Title + `</h1>`
	}

	var subAddr string

	if setup.ConfigData.Users.Domain == "" {

		subAddr = r.Host

	} else {

		if setup.ConfigData.Users.Port == "" {
			subAddr = "https://" + setup.ConfigData.Users.Domain
		} else {
			subAddr = "https://" + setup.ConfigData.Users.Domain + ":" + setup.ConfigData.Users.Port
		}

	}

	var setTagStr, usersLiSrt string
	var staticStr string
	var backupStr string
	var appUrl string

	users.UsersListHttp(subAddr, &setTagStr, &usersLiSrt)

	if setup.ConfigData.Static.Enabled && len(setup.ConfigData.Static.ConstList) == 0 {
		staticStr = `
	    </br></br>
	    <button onclick="renewUsersData('static')">设置已选中的用户为静态链接</button>
        `
	}

	if setup.ConfigData.Backup.Enabled {
		backupStr = `
		</br></br>
	    <button onclick="renewBackupSetup('exclude')">设置已选中的用户不使用备用链接</button>
		</br></br>
		<button onclick="renewBackupSetup('renew')">立即更新备用连接</button>
		`
	}

	app.GetAppUrl()
	appLen := len(app.AppsData) - 1
	for i := range app.AppsData {
		appUrl += `
	    <button onclick="copyContent('` + app.AppsData[i].Url + `')">
        ` + app.AppsData[i].Label + `复制</button>`

		if !app.AppsData[i].OnlyCopy {
			appUrl += `
			<button
			onclick="showQRCode('','` + app.AppsData[i].Url + `','` + app.AppsData[i].Label + `')">` + app.AppsData[i].Label + `二维码</button>`
		}

		if appLen != i {
			appUrl += `</br></br>`
		}

	}

	variables := struct {
		Logout        string
		TitleStr      template.HTML
		SetTagStr     template.HTML
		UsersLiSrt    template.HTML
		StaticStr     template.HTML
		BackupStr     template.HTML
		AppUrl        template.HTML
		ToggleContent string
		SetPostUrl    string
		RenewPostUrl  string
		BackupPostUrl string
	}{
		Logout:        setup.ConfigData.Server.Home.Url + "/logout",
		TitleStr:      template.HTML(titleStr),
		SetTagStr:     template.HTML(setTagStr),
		UsersLiSrt:    template.HTML(usersLiSrt),
		StaticStr:     template.HTML(staticStr),
		BackupStr:     template.HTML(backupStr),
		AppUrl:        template.HTML(appUrl),
		ToggleContent: ToggleContent,
		SetPostUrl:    setup.ConfigData.Server.Post.Set,
		RenewPostUrl:  setup.ConfigData.Server.Post.Renew,
		BackupPostUrl: setup.ConfigData.Server.Post.Backup,
	}

	ToggleContent = ""
	w.Header().Set("Content-Type", "text/html")

	tmpl, err := template.ParseFiles("html/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, variables)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
