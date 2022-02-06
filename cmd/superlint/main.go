package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"plugin"

	"github.com/ammario/superlint"
	"github.com/coder/flog"
	"github.com/spf13/cobra"
)

func main() {
	var verbose bool
	rootCmd := &cobra.Command{
		Use:     "superlint <ruleset plugin> [file regex]",
		Example: "superlint rules.so \".go$\"",
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			matcher := "(.*?)"
			if len(args) == 2 {
				matcher = args[1]
			}
			pl, err := plugin.Open(args[0])
			if err != nil {
				return fmt.Errorf("open %v: %w", args[0], err)
			}
			const rulesFuncName = "LoadRules"

			sym, err := pl.Lookup(rulesFuncName)
			if err != nil {
				return fmt.Errorf("lookup loader %s: %w", rulesFuncName, err)
			}

			loader, ok := sym.(*superlint.Loader)
			if !ok {
				return fmt.Errorf("loader not of type superlint.Loader")
			}

			log := &flog.Logger{
				W:          os.Stderr,
				TimeFormat: flog.ClockFormat + ".000",
			}

			rn := superlint.Runner{
				DebugLogger: flog.New(ioutil.Discard),
				Log:         log,
				Matcher:     matcher,
			}
			if verbose {
				rn.DebugLogger = log
			}
			rs := make(superlint.RuleSet, 0, 16)
			(*loader)(log, &rs)
			log.Info("loaded %v rules", len(rs))
			err = rn.Run(&rs)
			if err != nil {
				log.Fatal("%+v", err)
			}
			return nil
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode")

	err := rootCmd.Execute()
	if err != nil {
		flog.Fatal("%v", err)
	}
}
