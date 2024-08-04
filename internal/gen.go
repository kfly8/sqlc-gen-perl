package perl

import (
	//"bufio"
	"bufio"
	"bytes"
	"context"

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
	Q string
	SqlcVersion string
	SourceName string
}

func Generate(ctx context.Context, req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {
	options, err := opts.Parse(req)
	if err != nil {
		return nil, err
	}

	if err := opts.ValidateOpts(options); err != nil {
		return nil, err
	}

	return generate(req, options)
}


func generate(req *plugin.GenerateRequest, options *opts.Options) (*plugin.GenerateResponse, error) {

	tctx := tmplCtx{
		Q: "`",
		SqlcVersion: req.SqlcVersion,
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

