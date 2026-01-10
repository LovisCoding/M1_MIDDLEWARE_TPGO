package templater

import (
	"bytes"
	"embed"
	"strings"
	"text/template"

	"github.com/adrg/frontmatter"
)

//go:embed templates/*.html
var embeddedTemplates embed.FS

type Matter struct {
	Subject string `yaml:"subject"`
}

func GetStringFromEmbeddedTemplate(templateName string, body interface{}) (string, Matter, error) {
	var matter Matter

	// getting template file from embeddedTemplates
	// Note: We parse the raw file content first to extract frontmatter, then parse the template
	templatePath := "templates/" + templateName
	rawContent, err := embeddedTemplates.ReadFile(templatePath)
	if err != nil {
		return "", matter, err
	}

	// Separate frontmatter from content
	content, err := frontmatter.Parse(strings.NewReader(string(rawContent)), &matter)
	if err != nil {
		return "", matter, err
	}

	// Parse the content part as a Go template
	tmpl, err := template.New(templateName).Parse(string(content))
	if err != nil {
		return "", matter, err
	}

	var tpl bytes.Buffer
	if err = tmpl.Execute(&tpl, body); err != nil {
		return "", matter, err
	}

	// Re-process to interpolate subject variables if needed (optional, keeping simple for now)
	tmplSubject, err := template.New("subject").Parse(matter.Subject)
	if err == nil {
		var subjectBuffer bytes.Buffer
		if err := tmplSubject.Execute(&subjectBuffer, body); err == nil {
			matter.Subject = subjectBuffer.String()
		}
	}

	return tpl.String(), matter, nil
}
