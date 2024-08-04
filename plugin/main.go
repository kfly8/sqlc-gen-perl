package main

import (
	"github.com/sqlc-dev/plugin-sdk-go/codegen"

	perl "github.com/kfly8/sqlc-gen-perl/internal"
)

func main() {
	codegen.Run(perl.Generate)
}
