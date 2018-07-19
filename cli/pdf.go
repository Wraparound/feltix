package cli

import (
	"os"
	"time"

	"github.com/Wraparound/wrap/ast"
	"github.com/Wraparound/wrap/parser"
	"github.com/Wraparound/wrap/pdf"
	"github.com/spf13/cobra"
)

// pdfCmd represents the pdf command
var pdfCmd = &cobra.Command{
	Use:              "pdf [path to input file] [> output path]",
	Short:            "Export file as PDF",
	Args:             cobra.MaximumNArgs(1),
	TraverseChildren: true,
	Long:             longDescription,
	Run:              pdfRun,
}

var pdfNoscenenumbersFlag bool

func init() {
	pdfCmd.Flags().BoolVarP(&pdfNoscenenumbersFlag, "noscenenumbers", "s", false, "remove scenenumbers from output")

	WrapCmd.AddCommand(pdfCmd)
}

func pdfRun(cmd *cobra.Command, args []string) {
	startTime := time.Now()

	var (
		err    error
		output *os.File
		script *ast.Script
	)

	if len(args) == 0 {
		// Assume Wrap input
		parser.UseWrapExtensions = true

		script, err = getScriptFromStdin()
		handle(err)

		// Get the file to use during export.
		path, err := makeUnique("script", "pdf")
		handle(err)

		output = getOuput(path, "pdf")

	} else {
		pathToFile := args[0]

		if isWrapFile(pathToFile) {
			parser.UseWrapExtensions = true
		}

		script, err = parser.ParseFile(pathToFile)
		handle(err)

		// Get the file to use during export.
		output = getOuput(pathToFile, "pdf")
	}

	// Make sure to close the stream...
	defer output.Close()

	startExportTime := time.Now()

	if pdfNoscenenumbersFlag {
		pdf.AddSceneNumbers = false
	}

	err = pdf.WritePDF(script, output)
	handle(err)

	endTime := time.Now()

	if benchmarkFlag {
		printBenchmarks(startTime, startExportTime, endTime)
	}
}
