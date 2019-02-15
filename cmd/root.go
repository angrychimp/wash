package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Unfortunately, cobra.Command.Execute() can only return error objects.
// Thus, the only way for us to let each command configure its own exit
// code is to wrap that value in an error object. This should be OK since
// we want the commands to handle their own errors.
type exitCode struct {
	value int
}

// Required to implement the error interface
func (e exitCode) Error() string {
	return ""
}

// This munging's necessary to ensure that all commandMain functions return
// an exit code while also letting them be used as RunE functions that can
// be passed into Cobra. Otherwise, Go's type-checker will complain even though
// exitCode is an error object.
type commandMain func(cmd *cobra.Command, args []string) exitCode
type runE func(cmd *cobra.Command, args []string) error

func toRunE(main commandMain) runE {
	return func(cmd *cobra.Command, args []string) error {
		return main(cmd, args)
	}
}

func rootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		// TODO: Set this to "" when we're ready to ship so that
		// when we alias our custom commands, someone typing in
		// e.g. `meta --help` will not see `wash meta` in the usage
		Use: "wash",
		// Need to set these so that Cobra will not output the usage +
		// error object when Execute() returns an error, which will always
		// happen in our case because the exitCode object is technically
		// an error.
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.AddCommand(serverCommand())
	rootCmd.AddCommand(metaCommand())
	rootCmd.AddCommand(lsCommand())

	return rootCmd
}

// Execute executes the root command, returning the exit code
func Execute() int {
	err := rootCommand().Execute()
	if err == nil {
		// This can happen if the user invokes `wash` without any
		// arguments, or if they invoke a help command.
		return 0
	}

	exitCode, ok := err.(exitCode)
	if !ok {
		// err is something Cobra-related, like e.g. a malformed
		// flag. Print the error, then return.
		fmt.Printf("Error: %v\n", err)
		return 1
	}

	return exitCode.value
}