package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// The following variables are injected at build time via ldflags
var (
	Version   = "dev"
	Commit    = "none"
	BuildTime = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version info of hello-cli",
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		if verbose {
			fmt.Fprintf(cmd.OutOrStdout(), "hello-cli version: %s\n", Version)
			fmt.Fprintf(cmd.OutOrStdout(), "commit: %s\n", Commit)
			fmt.Fprintf(cmd.OutOrStdout(), "build time: %s\n", BuildTime)
			fmt.Fprintf(cmd.OutOrStdout(), "platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "hello-cli version %s %s/%s\n", Version, runtime.GOOS, runtime.GOARCH)
		}
	},
}

func init() {
	versionCmd.Flags().Bool("verbose", false, "print detailed version info")
	_ = versionCmd.Flags().MarkHidden("verbose")
	rootCmd.AddCommand(versionCmd)
}
