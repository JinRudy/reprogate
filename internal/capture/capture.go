package capture

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/JinRudy/reprogate/internal/redact"
	"github.com/JinRudy/reprogate/internal/report"
)

type Options struct {
	Command     []string
	WorkDir     string
	OutputPath  string
	MaxLogBytes int
	Stdout      io.Writer
	Stderr      io.Writer
}

type Result struct {
	ExitCode   int
	ReportPath string
}

func Run(ctx context.Context, opts Options) (Result, error) {
	if len(opts.Command) == 0 {
		return Result{}, fmt.Errorf("missing command")
	}
	if opts.WorkDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			return Result{}, err
		}
		opts.WorkDir = wd
	}
	if opts.OutputPath == "" {
		opts.OutputPath = filepath.Join(opts.WorkDir, ".reprogate", "repro.md")
	}
	if opts.MaxLogBytes <= 0 {
		opts.MaxLogBytes = 200000
	}

	start := time.Now()
	cmd := exec.CommandContext(ctx, opts.Command[0], opts.Command[1:]...)
	cmd.Dir = opts.WorkDir
	var logs limitedBuffer
	logs.limit = opts.MaxLogBytes
	cmd.Stdout = writerFor(&logs, opts.Stdout)
	cmd.Stderr = writerFor(&logs, opts.Stderr)
	err := cmd.Run()
	exitCode := 0
	if err != nil {
		exitCode = 1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}

	doc := report.Document{
		Command:     strings.Join(opts.Command, " "),
		ExitCode:    exitCode,
		Duration:    time.Since(start).Round(time.Millisecond).String(),
		Environment: probeEnvironment(),
		Dependency:  probeDependency(opts.WorkDir),
		Logs:        redact.Text(logs.String()),
	}

	if err := os.MkdirAll(filepath.Dir(opts.OutputPath), 0o755); err != nil {
		return Result{}, err
	}
	if err := os.WriteFile(opts.OutputPath, []byte(report.RenderMarkdown(doc)), 0o644); err != nil {
		return Result{}, err
	}
	return Result{ExitCode: exitCode, ReportPath: opts.OutputPath}, nil
}

func RunCLI(ctx context.Context, args []string, out io.Writer, errOut io.Writer) error {
	if len(args) == 0 || args[0] != "--" || len(args) == 1 {
		return fmt.Errorf("usage: reprogate capture -- <command> [args...]")
	}
	result, err := Run(ctx, Options{Command: args[1:], Stdout: out, Stderr: errOut})
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "wrote %s\n", result.ReportPath)
	if result.ExitCode != 0 {
		return fmt.Errorf("command exited with code %d", result.ExitCode)
	}
	return nil
}

type limitedBuffer struct {
	bytes.Buffer
	limit int
}

func (b *limitedBuffer) Write(p []byte) (int, error) {
	if b.Len() < b.limit {
		remaining := b.limit - b.Len()
		if len(p) > remaining {
			_, _ = b.Buffer.Write(p[:remaining])
		} else {
			_, _ = b.Buffer.Write(p)
		}
	}
	return len(p), nil
}

func writerFor(capture io.Writer, stream io.Writer) io.Writer {
	if stream == nil {
		return capture
	}
	return io.MultiWriter(capture, stream)
}

func probeEnvironment() map[string]string {
	values := map[string]string{
		"arch": runtime.GOARCH,
		"os":   runtime.GOOS,
	}
	if version, err := exec.Command("go", "version").Output(); err == nil {
		values["go"] = strings.TrimSpace(string(version))
	}
	if version, err := exec.Command("docker", "--version").Output(); err == nil {
		values["docker"] = strings.TrimSpace(string(version))
	}
	return values
}

func probeDependency(workDir string) map[string]string {
	values := map[string]string{}
	for _, name := range []string{"go.sum", "package-lock.json", "pnpm-lock.yaml", "yarn.lock"} {
		path := filepath.Join(workDir, name)
		if data, err := os.ReadFile(path); err == nil {
			values[name] = fmt.Sprintf("%d bytes", len(data))
		}
	}
	return values
}
