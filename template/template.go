package template

import (
	"html/template"
	"path/filepath"
	"reflect"
	"strings"
)

// Template is main entity of package template
type Template map[string]*template.Template

// New to create new instance of template
func New() (Template, error) {
	var templates = make(map[string]*template.Template)

	layoutFiles, err := filepath.Glob("files/var/www/layout/" + "*.html")
	if err != nil {
		return nil, err
	}

	includeFiles, err := filepath.Glob("files/var/www/" + "*.html")
	if err != nil {
		return nil, err
	}

	for _, file := range includeFiles {
		fullFileName := filepath.Base(file)
		spltFullFileName := strings.Split(fullFileName, ".")
		fileName := spltFullFileName[0]
		files := append(layoutFiles, file)

		templates[fileName] = template.Must(template.New("").Funcs(template.FuncMap{
			"Field": func(v interface{}, name string) interface{} {
				var data interface{}
				ref := reflect.ValueOf(v)
				splt := strings.Split(name, "|")
				if len(splt) > 1 {
					buffer := ref.FieldByName(splt[0]).Interface()
					ref = reflect.ValueOf(buffer)
					data = ref.FieldByName(splt[1]).Interface()
				} else {
					data = ref.FieldByName(name).Interface()
				}

				return data
			},
		}).ParseFiles(files...))
	}

	return templates, nil
}
