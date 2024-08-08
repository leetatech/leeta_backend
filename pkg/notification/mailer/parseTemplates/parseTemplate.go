package parseTemplates

import (
	"bytes"
	"fmt"
	"github.com/leetatech/leeta_backend/pkg/errs"
	"github.com/leetatech/leeta_backend/services/models"
	"html/template"
	"io"
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
	basePath := "pkg/messaging/mailer/templates"
	page := filepath.Join(basePath, fileName)

	ts, err := template.New(fileName).Funcs(functions).ParseFiles(page)
	if err != nil {
		return fmt.Errorf("failed to parse template file %s: %w", fileName, err)
	}

	layoutPattern := filepath.Join(basePath, "*.layout.gohtml")
	ts, err = ts.ParseGlob(layoutPattern)
	if err != nil {
		return fmt.Errorf("failed to parse layout templates: %w", err)
	}

	if err = ts.Execute(buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	return nil
}
