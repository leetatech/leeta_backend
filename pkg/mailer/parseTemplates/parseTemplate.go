package parseTemplates

import (
	"bytes"
	"fmt"
	"github.com/leetatech/leeta_backend/services/models"
	"html/template"
	"io"
	"os"
	"path/filepath"
)

type Template struct{}

func CreateSingleTemplate(templatePath string, data models.Message) (string, error) {
	templateStore := map[string]string{}

	t := Template{}
	buf := new(bytes.Buffer)

	err := t.Create(buf, templatePath, data)
	if err != nil {
		return "", err
	}

	templateStore[templatePath] = buf.String()

	return templateStore[templatePath], nil
}

var functions = template.FuncMap{}

func (t *Template) Create(buf io.Writer, fileName string, data any) error {
	dir, err := os.Getwd()
	page := fmt.Sprintf("%s/%s", filepath.Join(dir, "pkg/mailer/templates"), fileName)

	ts, err := template.New(fileName).Funcs(functions).ParseFiles(page)

	if err != nil {
		return err
	}

	ts, err = ts.ParseGlob("./pkg/mailer/templates/*.layout.gohtml")

	if err != nil {
		return err
	}

	if err = ts.Execute(buf, data); err != nil {
		return err
	}
	return nil
}
