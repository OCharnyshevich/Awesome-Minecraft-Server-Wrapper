package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

var templates = template.Must(template.ParseFiles("static/edit.html", "static/view.html", "static/index.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

type Page struct {
	Title string
	Body  []byte
}

type Web struct {
	Port string
	List map[string]*Page
}

func (web *Web) loadPage(title string) (*Page, error) {
	if page, ok := web.List[title]; !ok {
		return nil, fmt.Errorf("page '%s' not found", title)
	} else {
		return page, nil
	}

}

func (web *Web) save(page *Page) {
	web.List[page.Title] = page
	web.saveWeb()
}

func (web *Web) editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := web.loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func (web *Web) viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := web.loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func (web *Web) saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	web.save(p)
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func (web *Web) listHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, fmt.Sprintf("%s.html", "index"), web)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (web Web) saveWeb() {
	file, _ := json.MarshalIndent(web.List, "", " ")

	err := ioutil.WriteFile(fmt.Sprintf("%s/static/.web.json", web.currentDir()), file, 0777)
	if err != nil {
		log.Fatalf("Web can't be saved: %s", err)
	} else {
		log.Println("Web saved!")
	}
}

func (web Web) loadWeb() error {
	file, err := ioutil.ReadFile(fmt.Sprintf("%s/static/.web.json", web.currentDir()))

	if err != nil {
		if !os.IsExist(err) {
			web.saveWeb()
			return nil
		}

		return err
	}

	data := web.List

	err = json.Unmarshal(file, &data)

	return err
}

func (web Web) currentDir() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("can't get current dir '%v'", err)
	}

	return pwd
}

func NewWebServer() error {
	pages := &Web{
		Port: ":8080",
		List: map[string]*Page{},
	}
	err := pages.loadWeb()
	if err != nil {
		return err
	}

	http.HandleFunc("/", pages.listHandler)
	http.HandleFunc("/view/", makeHandler(pages.viewHandler))
	http.HandleFunc("/edit/", makeHandler(pages.editHandler))
	http.HandleFunc("/save/", makeHandler(pages.saveHandler))

	return http.ListenAndServe(pages.Port, nil)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return basicAuth(func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	})
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, fmt.Sprintf("%s.html", tmpl), p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()

		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte("myUsername"))
			expectedPasswordHash := sha256.Sum256([]byte("myPassword"))

			usernameMatch := subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1
			passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
