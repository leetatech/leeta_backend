package parseTemplates

import (
	"bytes"
	"fmt"
	"github.com/leetatech/leeta_backend/pkg/errs"
	"github.com/leetatech/leeta_backend/services/models"
	"html/template"
	"io"
	"os"
	"path/filepath"
)

type Template struct{}

func CreateSingleTemplate(templatePath string, data models.Message) (string, error) {
	var buf bytes.Buffer

	t := Template{}

	err := t.Create(&buf, templatePath, data)
	if err != nil {
		return "", errs.Body(errs.TemplateCreationError, err)
	}

	return buf.String(), nil
}

var functions = template.FuncMap{}

func (t *Template) Create(buf io.Writer, fileName string, data any) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}
	page := fmt.Sprintf("%s/%s", filepath.Join(dir, "pkg/messaging/mailer/templates"), fileName)

	ts, err := template.New(fileName).Funcs(functions).ParseFiles(page)

	if err != nil {
		return fmt.Errorf("failed to parse template file %s: %w", fileName, err)
	}

	ts, err = ts.ParseGlob(filepath.Join(dir, "pkg/messaging/mailer/templates/*.layout.gohtml"))

	if err != nil {
		return fmt.Errorf("failed to parse layout templates: %w", err)
	}

	if err = ts.Execute(buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	return nil
}
