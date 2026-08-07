package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	AppName    = "git-log-color"
	AppAuthor  = "Hadi Cahyadi"
	AppHome    = "github.com/cumulus13/git-log-color"
	AppVersion = "dev"
)

type CLIOptions struct {
	Path       string
	Limit      int
	All        bool
	OneLine    bool
	Format     string
	ConfigPath string
	NoColor    bool
	NoIcons    bool
	Full       bool
	GitBin     string
	ExtraArgs  []string
}

func main() {
	options, err := parseCLI(os.Args[1:])
	if err != nil {
		exitErr(err)
	}

	resolvedPath, err := resolveDirectory(options.Path)
	if err != nil {
		exitErr(err)
	}

	cfg, err := LoadConfig(options.ConfigPath, resolvedPath)
	if err != nil {
		exitErr(err)
	}

	if options.NoColor {
		cfg.UI.ColorEnabled = false
	}
	if options.NoIcons {
		cfg.UI.IconsEnabled = false
	}

	pager, out := StartPager(!options.Full, os.Stdout)

	app := NewApplication(cfg, options)
	runErr := app.Run(resolvedPath, out)
	pager.Wait()
	if runErr != nil {
		exitErr(runErr)
	}
}

func parseCLI(args []string) (CLIOptions, error) {
	opts := CLIOptions{}
	cliArgs, extraArgs, inferredPath := normalizeCLIArgs(args)
	fs := flag.NewFlagSet(AppName, flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	fs.StringVar(&opts.Path, "path", ".", "repository path (supports relative paths)")
	fs.StringVar(&opts.Path, "C", ".", "repository path (git-style shorthand)")
	fs.IntVar(&opts.Limit, "n", 10, "number of commits")
	fs.IntVar(&opts.Limit, "limit", 10, "number of commits")
	fs.BoolVar(&opts.All, "all", false, "include all branches")
	fs.BoolVar(&opts.OneLine, "oneline", false, "use compact one-line view for text output")
	fs.StringVar(&opts.Format, "format", "text", "output format: text,json,table,xml,yaml,toml")
	fs.StringVar(&opts.ConfigPath, "config", "", "config file path (JSON)")
	fs.BoolVar(&opts.NoColor, "no-color", false, "disable ANSI colors")
	fs.BoolVar(&opts.NoIcons, "no-icons", false, "disable emoji/icons")
	fs.BoolVar(&opts.Full, "f", false, "print full output directly, bypassing the pager (shorthand)")
	fs.BoolVar(&opts.Full, "full", false, "print full output directly, bypassing the pager")
	fs.StringVar(&opts.GitBin, "git-bin", "git", "git executable path")
	showVersion := fs.Bool("version", false, "print version information")
	showHelp := fs.Bool("help", false, "show help")

	if err := fs.Parse(cliArgs); err != nil {
		return opts, err
	}

	if *showHelp {
		printUsage(fs)
		os.Exit(0)
	}

	if *showVersion {
		fmt.Printf("%s %s\nAuthor: %s\nHomepage: %s\n", AppName, AppVersion, AppAuthor, AppHome)
		os.Exit(0)
	}

	if opts.Limit <= 0 {
		return opts, errors.New("limit must be greater than zero")
	}

	opts.Format = strings.ToLower(strings.TrimSpace(opts.Format))
	if !isSupportedFormat(opts.Format) {
		return opts, fmt.Errorf("unsupported format %q", opts.Format)
	}

	pathProvidedByFlag := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "path" || f.Name == "C" {
			pathProvidedByFlag = true
		}
	})

	if !pathProvidedByFlag && inferredPath != "" {
		opts.Path = inferredPath
	}

	opts.ExtraArgs = append(fs.Args(), extraArgs...)
	return opts, nil
}

func printUsage(fs *flag.FlagSet) {
	fmt.Fprintf(os.Stdout, "%s - colorful and structured git log output\n\n", AppName)
	fmt.Fprintf(os.Stdout, "Author: %s\nHomepage: %s\n\n", AppAuthor, AppHome)
	fmt.Fprintf(os.Stdout, "Usage:\n  %s [flags] [-- git log extra args]\n\n", AppName)
	fmt.Fprintf(os.Stdout, "Examples:\n")
	fmt.Fprintf(os.Stdout, "  %s --path . --format text\n", AppName)
	fmt.Fprintf(os.Stdout, "  %s -C .. --all --format table\n", AppName)
	fmt.Fprintf(os.Stdout, "  %s --format json -- --author=Hadi -- README.md\n\n", AppName)
	fmt.Fprintf(os.Stdout, "Output is paged through $GIT_PAGER/$PAGER (or less) by default when\n")
	fmt.Fprintf(os.Stdout, "attached to a terminal. Pass -f/--full to print everything directly.\n\n")
	fs.PrintDefaults()
}

func exitErr(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}

func resolveDirectory(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		path = "."
	}

	expanded := os.ExpandEnv(path)
	abs, err := filepath.Abs(expanded)
	if err != nil {
		return "", err
	}

	info, err := os.Stat(abs)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%s is not a directory", abs)
	}

	return abs, nil
}

func inferPositionalRepoPath(value string) (string, bool) {
	if strings.TrimSpace(value) == "" {
		return "", false
	}

	abs, err := filepath.Abs(os.ExpandEnv(value))
	if err != nil {
		return "", false
	}

	info, err := os.Stat(abs)
	if err != nil || !info.IsDir() {
		return "", false
	}

	return value, true
}

func normalizeCLIArgs(args []string) ([]string, []string, string) {
	flagNamesWithValue := map[string]bool{
		"-C":        true,
		"--path":    true,
		"-n":        true,
		"--limit":   true,
		"--format":  true,
		"--config":  true,
		"--git-bin": true,
	}

	var cliArgs []string
	var extraArgs []string
	var inferredPath string
	hasExplicitPath := false

	for index := 0; index < len(args); index++ {
		arg := args[index]

		if arg == "--" {
			extraArgs = append(extraArgs, args[index+1:]...)
			break
		}

		if strings.HasPrefix(arg, "-") {
			cliArgs = append(cliArgs, arg)
			if arg == "-C" || arg == "--path" {
				hasExplicitPath = true
			}

			if flagNamesWithValue[arg] && index+1 < len(args) {
				cliArgs = append(cliArgs, args[index+1])
				index++
			}
			continue
		}

		if !hasExplicitPath && inferredPath == "" {
			if candidate, ok := inferPositionalRepoPath(arg); ok {
				inferredPath = candidate
				continue
			}
		}

		extraArgs = append(extraArgs, arg)
	}

	return cliArgs, extraArgs, inferredPath
}
