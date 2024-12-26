package file_operations

import (
	"fmt"
	"html/template"
	"net/url"
	"os"
	"path"

	. "github.com/emad-elsaid/xlog"
)

type PageRename struct {
	page Page
}

func (PageRename) Icon() string         { return "fa-solid fa-i-cursor" }
func (PageRename) Name() string         { return "Rename" }
func (PageRename) OnClick() template.JS { return "" }
func (PageRename) Link() string         { return "" }
func (f PageRename) Attrs() map[template.HTMLAttr]any {
	return map[template.HTMLAttr]any{
		"hx-get":    "/+/file/rename?page=" + url.QueryEscape(f.page.Name()),
		"hx-target": "body",
		"hx-swap":   "beforeend",
	}
}

func (f PageRename) Widget() template.HTML {
	return ""
}

func (f PageRename) Form(r Request) Output {
	name := r.FormValue("page")
	page := NewPage(name)

	return Render("rename-form", map[string]any{
		"page": page,
	})
}

func (f PageRename) Handler(r Request) Output {
	old := NewPage(r.FormValue("old"))
	if old == nil || !old.Exists() {
		return BadRequest("file doesn't exist")
	}

	ext := path.Ext(old.FileName())
	basename := r.FormValue("new")
	newName := basename + ext

	os.Rename(old.FileName(), newName)
	old.Write(Markdown(fmt.Sprintf("Renamed to: %s", basename)))

	return func(w Response, r Request) {
		w.Header().Add("HX-Redirect", "/"+basename)
	}
}
