package opts

import (
	"encoding/json"
	"fmt"
	//"maps"
	///"path/filepath"

	"github.com/sqlc-dev/plugin-sdk-go/plugin"
)

type Options struct {
	Package                     string            `json:"package" yaml:"package"`
	Rename                      map[string]string `json:"rename,omitempty" yaml:"rename"`
	OutputModelsFileName        string            `json:"output_models_file_name,omitempty" yaml:"output_models_file_name"`
}

func Parse(req *plugin.GenerateRequest) (*Options, error) {
	options, err := parseOpts(req)
	if err != nil {
		return nil, err
	}

	return options, nil
}

func parseOpts(req *plugin.GenerateRequest) (*Options, error) {
	var options Options
	if len(req.PluginOptions) == 0 {
		return &options, nil
	}
	if err := json.Unmarshal(req.PluginOptions, &options); err != nil {
		return nil, fmt.Errorf("unmarshalling plugin options: %w", err)
	}

	return &options, nil
}

func ValidateOpts(opts *Options) error {
	// TODO: validate options

	return nil
}

