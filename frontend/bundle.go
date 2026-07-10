package main

import (
	_ "embed"
	"fmt"
	"os"
	"slices"
	"strings"
	"text/template"

	"github.com/evanw/esbuild/pkg/api"
)

//go:embed index.tpl.html
var indexTpl string

func main() {
	result := api.Build(api.BuildOptions{
		LogLevel:          api.LogLevelDebug,
		EntryPoints:       []string{"./main.tsx"},
		Outdir:            ".",
		Bundle:            true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Write:             false,
		TreeShaking:       api.TreeShakingTrue,
		Platform:          api.PlatformBrowser,
		Target:            api.ESNext,
		Format:            api.FormatESModule,
		Define: map[string]string{
			"process.env.NODE_ENV": "\"production\"",
		},
		Tsconfig: "./tsconfig.json",
		Loader: map[string]api.Loader{
			".woff2": api.LoaderDataURL,
			".woff":  api.LoaderDataURL,
		},
	})

	if len(result.Errors) > 0 {
		os.Exit(1)
	}

	jsIndex := slices.IndexFunc(result.OutputFiles, func(s api.OutputFile) bool {
		return strings.HasSuffix(s.Path, "main.js")
	})
	cssIndex := slices.IndexFunc(result.OutputFiles, func(s api.OutputFile) bool {
		return strings.HasSuffix(s.Path, "main.css")
	})

	tmpl, err := template.New("foo").Parse(indexTpl)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(1)
	}

	f, err := os.Create("./.dist/index.html")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(1)
	}

	x := struct {
		Css string
		Js  string
	}{
		Css: string(result.OutputFiles[cssIndex].Contents),
		Js:  string(result.OutputFiles[jsIndex].Contents),
	}

	if err := tmpl.Execute(f, x); err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(1)
	}
}
