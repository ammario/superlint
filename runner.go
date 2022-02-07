package superlint

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sync/atomic"

	"github.com/coder/flog"
)

// Runner runs the RuleSet.
type Runner struct {
	Matcher     string
	DebugLogger *flog.Logger
	Log         *flog.Logger
	failed      int64
}

// FileInfo is an os.FileInfo combined with a fully qualified path.
type FileInfo struct {
	os.FileInfo
	Path string
}

func (rn *Runner) runRule(files []FileInfo, r Rule) error {
	log := rn.DebugLogger.WithPrefix("%v: ", r.Name)

	if r.Linter == nil {
		return fmt.Errorf("no validator configured")
	}

	matchedFiles := make(map[string]FileInfo)
	for _, finfo := range files {
		matchedFiles[finfo.Name()] = finfo
	}
	return r.Linter(matchedFiles, func(ref FileReference, message string) {
		atomic.AddInt64(&rn.failed, 1)

		fmt.Printf("%v: %v: %v\n", r.Name, ref.Name, message)
		file, err := ioutil.ReadFile(ref.Name)
		if err != nil {
			log.Error("read %v: %v", ref.Name, err)
			return
		}
		prettyPrintReference(os.Stdout, file, ref)
	})
}

func (rn *Runner) Run(rs *RuleSet) error {
	var matches []FileInfo
	matchRegex, err := regexp.Compile(rn.Matcher)
	if err != nil {
		return fmt.Errorf("compile matcher: %w", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get wd: %w", err)
	}

	err = filepath.Walk(wd, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		path, err = filepath.Rel(wd, path)
		if err != nil {
			return err
		}

		if matchRegex.MatchString(path) {
			finfo, err := os.Stat(path)
			if err != nil {
				return err
			}
			matches = append(matches, FileInfo{
				FileInfo: finfo,
				Path:     path,
			})
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk: %w", err)
	}

	rn.DebugLogger.Info("%s matched %v files", rn.Matcher, len(matches))
	for _, r := range *rs {
		err := rn.runRule(matches, r)
		if err != nil {
			rn.Log.Error("%v: %v", r.Name, err)
		}
	}
	if rn.failed > 0 {
		return fmt.Errorf("%v violations found", rn.failed)
	}
	return nil
}
