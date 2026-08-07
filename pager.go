package main

import (
	"io"
	"os"
	"os/exec"
	"strings"
)

// Pager wraps an external pager process (e.g. "less") that output can be
// streamed into, mirroring how `git log` pages its own output by default.
type Pager struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	active bool
}

// StartPager launches a pager and returns both the Pager handle (for
// Wait) and the io.Writer the caller should write output to.
//
// If enabled is false, stdout isn't a terminal (e.g. output is piped or
// redirected), or no pager program can be found, it falls back to
// writing directly to fallback and returns a no-op Pager.
func StartPager(enabled bool, fallback *os.File) (*Pager, io.Writer) {
	if !enabled || !isTerminal(fallback) {
		return &Pager{}, fallback
	}

	name, args := pagerCommand()
	if name == "" {
		return &Pager{}, fallback
	}

	cmd := exec.Command(name, args...)
	cmd.Stdout = fallback
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return &Pager{}, fallback
	}

	if err := cmd.Start(); err != nil {
		return &Pager{}, fallback
	}

	return &Pager{cmd: cmd, stdin: stdin, active: true}, stdin
}

// Wait closes the pager's input stream and blocks until the user exits
// the pager (e.g. presses "q"). Safe to call more than once, and safe to
// call on a no-op Pager returned when paging was skipped.
func (p *Pager) Wait() {
	if p == nil || !p.active {
		return
	}
	p.active = false
	_ = p.stdin.Close()
	_ = p.cmd.Wait()
}

// pagerCommand resolves which pager to run, following the same
// precedence git itself uses: $GIT_PAGER, then $PAGER, then a built-in
// default ("less", with sane flags).
func pagerCommand() (string, []string) {
	if value := strings.TrimSpace(os.Getenv("GIT_PAGER")); value != "" {
		return splitPagerCommand(value)
	}
	if value := strings.TrimSpace(os.Getenv("PAGER")); value != "" {
		return splitPagerCommand(value)
	}

	if path, err := exec.LookPath("less"); err == nil {
		// -R: render raw ANSI color codes instead of escaping them
		// -F: exit immediately if the content fits on one screen
		// -X: don't clear the screen on exit
		return path, []string{"-R", "-F", "-X"}
	}

	if path, err := exec.LookPath("more"); err == nil {
		return path, nil
	}

	return "", nil
}

func splitPagerCommand(value string) (string, []string) {
	fields := strings.Fields(value)
	if len(fields) == 0 {
		return "", nil
	}
	return fields[0], fields[1:]
}

// isTerminal reports whether f is attached to a terminal. Used to skip
// paging automatically when output is redirected or piped, same as git.
func isTerminal(f *os.File) bool {
	if f == nil {
		return false
	}
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}
