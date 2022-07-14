package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/gobwas/glob"
	"gopkg.in/yaml.v3"
)

type Job struct {
	BeforeScipt []string `yaml:"before_script"`
	Script      []string `yaml:"script"`
	AfterScipt  []string `yaml:"after_script"`
	Image       string   `yaml:"image"`
	Stage       string   `yaml:"stage"`
}

func (j Job) getScript() string {
	builder := strings.Builder{}
	for _, l := range j.BeforeScipt {
		builder.Write([]byte(l))
		builder.Write([]byte("\n"))
	}
	for _, l := range j.Script {
		builder.Write([]byte(l))
		builder.Write([]byte("\n"))
	}
	for _, l := range j.AfterScipt {
		builder.Write([]byte(l))
		builder.Write([]byte("\n"))
	}
	return builder.String()
}

func (j Job) inferShell() string {
	switch img := strings.ToLower(j.Image); {
	case strings.HasPrefix(img, "docker"):
		return "sh"
	case strings.HasPrefix(img, "python"):
		return "bash"
	default:
		fmt.Printf("Cannot infer shell type for '%s'. Defaulting to 'bash'\n", j.Image)
		return "bash"
	}
}

func main() {

	var filterArg string

	flag.StringVar(&filterArg, "filter", "*", "Filter by job name.")

	flag.Parse()
	posArgs := flag.Args()

	filter := glob.MustCompile(filterArg)

	if len(posArgs) == 0 {
		fmt.Print("No path to .gitlab-ci.yml file\n")
		os.Exit(0)
	}

	content, err := os.ReadFile(posArgs[0])
	if err != nil {
		fmt.Printf("Error reading the file: %v\n", err)
		os.Exit(1)
	}

	jobs := make(map[string]Job, 1)

	err = yaml.Unmarshal(content, &jobs)
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	dir, err := os.MkdirTemp(os.TempDir(), "sh-check")
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	defer os.RemoveAll(dir)

	for name, job := range jobs {
		if !filter.Match(name) {
			continue
		}
		tmpFile := path.Join(dir, name)
		fp, err := os.Create(tmpFile)
		if err != nil {
			fmt.Printf("Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer fp.Close()

		fp.Write([]byte(job.getScript()))

		cmd := exec.Command("shellcheck", fmt.Sprintf("--shell=%s", job.inferShell()), tmpFile)
		out, err := cmd.Output()
		// Shellcheck output is != nil when any lint errors arise
		if err != nil {
			fmt.Print("===============================================\n")
			fmt.Printf("%s\n", out)
		}
	}
}
