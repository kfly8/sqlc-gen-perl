package perl

import (
	//"bufio"
	"bufio"
	"bytes"
	"context"
	"sort"
	"unicode"

	//"errors"
	//"fmt"
	//"go/format"
	"strings"
	"text/template"

	"github.com/kfly8/sqlc-gen-perl/internal/opts"
	"github.com/sqlc-dev/plugin-sdk-go/sdk"

	//"github.com/sqlc-dev/plugin-sdk-go/metadata"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
)


type tmplCtx struct {
	ModelPackage string
	Structs []Struct
	SqlcVersion string
	SourceName string
}

type Struct struct {
	Table *plugin.Identifier
	Name string
	Fields []Field
	Comment string
}

type Field struct {
	Name string // TODO: CamelCase ? snaked-case?
	DBName string // Name as used in the DB
	Type string
	Comment string
}

func Generate(ctx context.Context, req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {
	options, err := opts.Parse(req)
	if err != nil {
		return nil, err
	}

	if err := opts.ValidateOpts(options); err != nil {
		return nil, err
	}

	structs := buildStructs(req, options)

	return generate(req, options, structs)
}


func generate(req *plugin.GenerateRequest, options *opts.Options, structs []Struct) (*plugin.GenerateResponse, error) {

	modelPackage := "Models"
	if options.ModelPackage != "" {
		modelPackage = options.ModelPackage
	}

	tctx := tmplCtx{
		ModelPackage: modelPackage,
		SqlcVersion: req.SqlcVersion,
		Structs: structs,
	}

	funcMap := template.FuncMap{
		"lowerTitle": sdk.LowerTitle,
		"comment":    sdk.DoubleSlashComment,
		"escape":     sdk.EscapeBacktick,
		"hasPrefix":  strings.HasPrefix,
	}

	tmpl := template.Must(
		template.New("table").
			Funcs(funcMap).
			ParseFS(
				templates,
				"templates/*.tmpl",
			),
	)

	output := map[string]string{}

	execute := func(name, templateName string) error {
		var b bytes.Buffer
		w := bufio.NewWriter(&b)
		tctx.SourceName = name
		err := tmpl.ExecuteTemplate(w, templateName, &tctx)
		w.Flush()
		if err != nil {
			return err
		}

		if !strings.HasSuffix(name, ".pm") {
			name += ".pm"
		}

		output[name] = b.String()
		return nil
	}

	modelsFileName := "Models.pm"
	if options.OutputModelsFileName != "" {
		modelsFileName = options.OutputModelsFileName
	}
	if err := execute(modelsFileName, "modelsFile"); err != nil {
		return nil, err
	}

	resp := plugin.GenerateResponse{}

	for filename, code := range output {
		resp.Files = append(resp.Files, &plugin.File{
			Name:     filename,
			Contents: []byte(code),
		})
	}

	return &resp, nil
}

func buildStructs(req *plugin.GenerateRequest, options *opts.Options) []Struct {
	var structs []Struct
	for _, schema := range req.Catalog.Schemas {
		if schema.Name == "pg_catalog" || schema.Name == "information_schema" {
			continue
		}
		for _, table := range schema.Tables {
			var tableName string
			if schema.Name == req.Catalog.DefaultSchema {
				tableName = table.Rel.Name
			} else {
				tableName = schema.Name + "_" + table.Rel.Name
			}

			s := Struct{
				Table:   &plugin.Identifier{Schema: schema.Name, Name: table.Rel.Name},
				Name:    structName(tableName, options),
				Comment: table.Comment,
			}
			for _, column := range table.Columns {
				s.Fields = append(s.Fields, Field{
					Name:    fieldName(column.Name, options),
					Type:    perlType(req, options, column),
					Comment: column.Comment,
				})
			}
			structs = append(structs, s)
		}
	}
	if len(structs) > 0 {
		sort.Slice(structs, func(i, j int) bool { return structs[i].Name < structs[j].Name })
	}
	return structs
}

func structName(name string, options *opts.Options) string {
	if rename := options.Rename[name]; rename != "" {
		return rename
	}
	out := ""
	name = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) {
			return r
		}
		if unicode.IsDigit(r) {
			return r
		}
		return rune('_')
	}, name)

	for _, p := range strings.Split(name, "_") {
		if p == "id" {
			out += "ID"
		} else {
			out += strings.Title(p)
		}
	}

	return out;
}

func fieldName(name string, options *opts.Options) string {
	if rename := options.Rename[name]; rename != "" {
		return rename
	}
	return name;
}

func perlType(req *plugin.GenerateRequest, options *opts.Options, column *plugin.Column) string {
	// TODO: implement
	return "Str"
}
