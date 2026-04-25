package main

import (
	"fmt"
	"os"
)

type Application struct {
	config  Config
	options CLIOptions
}

func NewApplication(config Config, options CLIOptions) *Application {
	return &Application{
		config:  config,
		options: options,
	}
}

func (a *Application) Run(repoPath string) error {
	client := GitClient{
		Bin:      a.options.GitBin,
		RepoPath: repoPath,
	}

	if err := client.ValidateRepo(); err != nil {
		return err
	}

	if a.options.Format == "text" {
		lines, err := client.RunGraphLog(a.options.Limit, a.options.All, a.options.OneLine, a.options.ExtraArgs)
		if err != nil {
			return err
		}

		renderer := NewTextRenderer(a.config)
		for _, line := range lines {
			fmt.Fprintln(os.Stdout, renderer.RenderLine(line, a.options.OneLine))
		}
		return nil
	}

	commits, err := client.RunStructuredLog(a.options.Limit, a.options.All, a.options.ExtraArgs)
	if err != nil {
		return err
	}

	out, err := FormatCommits(a.options.Format, commits, a.config)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, out)
	return nil
}
