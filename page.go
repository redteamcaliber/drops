package drops

import (
	"html/template"
	"net/http"

	"context"
)

type MenuId uint8

// Menu configuration
type MenuCfg struct {
	MenuId  MenuId
	Ordinal int
}

type App struct {
	Menu          Menu
	Pages         map[string]Page
	Widgets       map[string]Widget
	TemplatePath  string
	Subdirectory  string
	Dev           bool // development mode, loads template from file
	TemplateFuncs template.FuncMap
}

func (t *App) Init() {
	loadIds(t)
	loadHandlers(t)
	loadTemplates(t, t.TemplatePath)
}

func (t *App) GetTemplate(id string) Template {
	tpl := t.Pages[id].Template
	if t.Dev {
		th := tpl.(*HtmlTemplate)
		tpl, s := loadTemplate(t.TemplatePath, th.File, id, t.TemplateFuncs)
		return NewHtmlTemplate(th.Name(), th.File, s, tpl)
	}
	return tpl
}

func (t *App) GetWidget(id string) Template {
	tpl := t.Widgets[id].Template
	if t.Dev {
		th := tpl.(*HtmlTemplate)
		tpl, s := loadTemplate(t.TemplatePath, th.File, id, t.TemplateFuncs)
		return NewHtmlTemplate(th.Name(), th.File, s, tpl)
	}
	return tpl
}

type Page struct {
	Id           string
	File         string
	Name         string
	Url          string
	Label        string
	Handler      http.Handler
	Menu         []MenuCfg
	Data         func(context.Context) (interface{}, error)
	Template     Template
	HtmlMenuItem string
	Parent       string
	Permission   int
	Submenu      MenuId
	Description  string
}

func (t *Page) HasMenu(mid MenuId) bool {
	for _, m := range t.Menu {
		if mid == m.MenuId {
			return true
		}
	}
	return false
}

func (t *Page) Ordinal(mid MenuId) int {
	for _, m := range t.Menu {
		if mid == m.MenuId {
			return m.Ordinal
		}
	}
	return 0
}

type Menu struct {
	Items []MenuItem
}

type MenuItem struct {
	Label   string
	Href    string
	Html    template.HTML
	Class   string
	Data    map[string]string
	Ordinal int
}

type MenuItemOrdered []MenuItem

func (a MenuItemOrdered) Len() int           { return len(a) }
func (a MenuItemOrdered) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a MenuItemOrdered) Less(i, j int) bool { return a[i].Ordinal < a[j].Ordinal }

type Footer struct {
	Template
}

type Widget struct {
	Id       string
	Name     string
	File     string
	Template Template
}

type SpfResponse struct {
	Title string                       `json:"title"`
	Url   string                       `json:"url"`
	Head  string                       `json:"head"`
	Body  map[string]string            `json:"body"`
	Attr  map[string]map[string]string `json:"attr"`
	Foot  string                       `json:"foot"`
}
