package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/dongpd224/bookings/pkg/config"
	"github.com/dongpd224/bookings/pkg/models"
)

var function = template.FuncMap{}

var app *config.AppConfig

// Newtemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDerfaultData(td *models.TemplateData) *models.TemplateData {

	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	// get the template cache from the app config

	// get requeseted template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDerfaultData(td)

	_ = t.Execute(buf, td)

	// render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}
	testCase := map[string]string{}
	testCase["hello"] = "tuyet voi"
	log.Println(testCase)

	// get all of the files named *.page.tmpl from./templtes
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}
	//  range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		log.Println(ts, 51)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		log.Println(matches, 56)

		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			log.Println(ts, 63)

			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	log.Println(myCache)
	return myCache, nil
}
