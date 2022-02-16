// Package cmd /*
package cmd

import (
	"fmt"
	"github.com/NYTimes/logrotate"
	"github.com/apex/log"
	"github.com/apex/log/handlers/multi"
	"github.com/spf13/cobra"
	log2 "log"
	"math"
	"math/rand"
	"newthon/logger/cli"
	"os"
	"path/filepath"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "newthon",
	Short: "A brief description of your application",
	Run:   rootCmdRun,
	PreRun: func(cmd *cobra.Command, args []string) {
		initLogging()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().Int64P("max", "m", -1, "Set a maximum of trials")
	rootCmd.Flags().Float64P("diff", "d", 0.0000000000001, "Define how closely you want.")
}

func rootCmdRun(cmd *cobra.Command, _ []string) {
	log.Debug("running in debug mode")
	var (
		precision float64
		max       int64
	)
	precision, _ = cmd.Flags().GetFloat64("diff")
	log.Info("receiving precision flag")
	max, _ = cmd.Flags().GetInt64("max")
	log.Info("receiving max flag")

	if precision <= 0.00000000000000006969 {
		old := precision
		precision = 0.0000000000001
		log.WithFields(log.Fields{
			"old": old,
			"new": precision,
		}).Info("changed precision.")
	}

	fmt.Printf("%v\n\n", solve(10, precision, max))
}

func f(x float64) float64 {
	return math.Pow(x, 3) + math.Cos(x) + math.Sin(math.Pow(x, 2))
}

func solve(x float64, precision float64, max int64) []float64 {
	var xs []float64
	var tries uint64

	for true {
		var x float64

		if max != -1 {
			tries++
		}

		if len(xs) == 0 {
			x = math.Mod(rand.Float64(), 512)
		} else {
			x = xs[len(xs)-1]
		}

		x1 := x - f(x)/slope(x)
		if math.Abs(x1-x) <= precision || tries == uint64(max) {
			break
		}
		xs = append(xs, x1)
	}
	return xs
}

func slope(x float64) float64 {
	offset := 0.0001
	halfoffset := offset / 2
	x1 := x - halfoffset
	x2 := x + halfoffset
	return (f(x2) - f(x1)) / (x2 - x1)
}

func initLogging() {
	p := filepath.Join(".", "/atlas.log")
	w, err := logrotate.NewFile(p)
	if err != nil {
		log2.Fatalf("cmd/root: failed to create atlas log: %s", err)
	}
	log.SetLevel(log.InfoLevel)
	log.SetHandler(multi.New(cli.Default, cli.New(w.File, false)))
	log.WithField("path", p).Info("writing log files to disk")
}
