package perl

import (
	//"bufio"
	//"bytes"
	"context"
	//"errors"
	//"fmt"
	//"go/format"
	//"strings"
	//"text/template"

	//"github.com/sqlc-dev/plugin-sdk-go/sdk"
	//"github.com/sqlc-dev/plugin-sdk-go/metadata"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
)

func Generate(ctx context.Context, req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {
	return generate(req)
}


func generate(req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {

	output := map[string]string{}
	output["Query.pm"] = "package Query";

	resp := plugin.GenerateResponse{}

	for filename, code := range output {
		resp.Files = append(resp.Files, &plugin.File{
			Name:     filename,
			Contents: []byte(code),
		})
	}

	return &resp, nil
}

