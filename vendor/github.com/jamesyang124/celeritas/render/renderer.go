package render

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/CloudyKit/jet/v6"
)

type Render struct {
	Renderer   string
	RootPath   string
	Secure     bool
	Port       string
	ServerName string
	JetViews   *jet.Set
}

type TemplateData struct {
	IsAuthenticated bool
	IntMap          map[string]int
	StringMap       map[string]string
	FloatMap        map[string]float32
	DataMap         map[string]interface{}
	CSRFToken       string
	Port            string
	ServerName      string
	Secure          bool
}

func (c *Render) Page(w http.ResponseWriter, r *http.Request, view string, jetVariables, data interface{}) error {
	switch strings.ToLower(c.Renderer) {
	case "jet":
		return c.JetPage(w, r, view, jetVariables, data)
	case "go":
		return c.GoPage(w, r, view, data)
	}
	return nil
}

// render by default go template engine
func (c *Render) GoPage(w http.ResponseWriter, r *http.Request, view string, data interface{}) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/views/%s.page.tmpl", c.RootPath, view))
	if err != nil {
		return err
	}

	td := &TemplateData{}
	if data != nil {
		// we tend to not check assertion error in here,
		// later tmpl.Execute call would raise it if data type is invalid
		td = data.(*TemplateData)
	}

	err = tmpl.Execute(w, &td)
	if err != nil {
		return err
	}

	return nil
}

// render by jet template engine
func (c *Render) JetPage(w http.ResponseWriter, r *http.Request, view string, variables, data interface{}) error {
	var vars jet.VarMap
	if variables == nil {
		vars = make(jet.VarMap)
	} else {
		vars = variables.(jet.VarMap)
	}

	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}

	t, err := c.JetViews.GetTemplate(fmt.Sprintf("%s.jet", view))
	if err != nil {
		log.Println(err)
		return err
	}

	if err = t.Execute(w, vars, td); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
