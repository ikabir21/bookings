package render

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/ikabir21/bookings/internal/config"
	"github.com/ikabir21/bookings/internal/models"
	"github.com/justinas/nosurf"
)

func RenderTemplateTest(w http.ResponseWriter, tmpl string) {
	parsedTemplate, _ := template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.tmpl")
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println("error parsing template: ", err)
		return
	}
}

// METHOD - 1 : Caching template

// var cache = make(map[string]*template.Template)

// func RenderTemplate(w http.ResponseWriter, s string) {
// 	var tmpl *template.Template
// 	var err error
// 	_, isExists := cache[s]
// 	if !isExists {
// 		// create the template
// 		fmt.Println("creating template and adding to cache")
// 		err = createTemplateCache(s)
// 		if err != nil {
// 			fmt.Println("using cached template")
// 		}
// 	} else {
// 		// template is already present in cache
// 		fmt.Println("template is already present in cache")
// 	}
// 	tmpl = cache[s]
// 	err = tmpl.Execute(w, nil)
// 	if err != nil {
// 		fmt.Println("using cached template")
// 	}
// }

// func createTemplateCache(s string) error {
// 	templates := [] string {
// 		fmt.Sprintf("./templates/%s", s),
// 		"./templates/base.layout.tmpl",
// 	}
// 	// parse the templates
// 	tmpl, err := template.ParseFiles(templates...)

// 	if err != nil {
// 		return err
// 	}
// 	cache[s] = tmpl
// 	return nil
// }

var app *config.AppConfig

// NewTemplates renders templates using html/template
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	return td
}

// METHOD - 2 : Caching templates
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	var tempCache map[string]*template.Template
	if app.UseCache {
		// Get the template cache from the app config
		tempCache = app.TemplateCache

	} else {
		tempCache, _ = CreateTemplateCache()
	}

	// get requested templates from cache
	t, ok := tempCache[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}
	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)
	_ = t.Execute(buf, td)

	// render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := make(map[string]*template.Template)
	// get all of the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")

	if err != nil {
		return myCache, err
	}
	// rnage through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, nil
		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, nil
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, nil
			}
		}

		myCache[name] = ts
	}
	return myCache, nil
}
