// Package render is the firewall between view models and Go source: view
// models go in, source fragments come out, and the only mechanism is
// text/template execution over the embedded .tmpl files. No naming or type
// decisions happen here.
package render

import (
	"embed"
	"strings"
	"text/template"

	"github.com/deploymenttheory/go-sdk-appleservices/device_management/internal/codegen/view"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

var templates = template.Must(template.New("dm").ParseFS(templateFS, "templates/*.tmpl"))

func execute(name string, data any) (string, error) {
	var b strings.Builder
	if err := templates.ExecuteTemplate(&b, name, data); err != nil {
		return "", err
	}
	return b.String(), nil
}

// StructDecl renders one struct's type declaration (fields only).
func StructDecl(s *view.Struct) (string, error) { return execute("structdecl", s) }

// StructFuncs renders one struct's methods: the wire-identifier method and
// Validate.
func StructFuncs(s *view.Struct) (string, error) { return execute("structfuncs", s) }

// EnumBlock renders one allowed-values const block.
func EnumBlock(e *view.EnumBlock) (string, error) { return execute("enumblock", e) }

// Registry renders a family's identifier->factory registry.
func Registry(r *view.Registry) (string, error) { return execute("registry", r) }
