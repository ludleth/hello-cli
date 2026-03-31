package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

const repo = "ludleth/hello-cli"

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update hello-cli to the latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		preview, _ := cmd.Flags().GetBool("preview")

		if Version == "dev" {
			fmt.Fprintln(cmd.ErrOrStderr(), "Running development version. Skipping update.")
			return nil
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Checking for latest version of %s...\n", repo)
		latest, found, err := selfupdate.DetectLatest(repo)
		if err != nil {
			return fmt.Errorf("error occurred while detecting version: %w", err)
		}
		if !found {
			return fmt.Errorf("no release found on GitHub")
		}

		if !preview && len(latest.Version.Pre) > 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "Latest version found (%s) is a pre-release. Skipping.\n", latest.Version)
			return nil
		}

		currentVersion := strings.TrimPrefix(Version, "v")
		if latest.Version.String() == currentVersion {
			fmt.Fprintln(cmd.OutOrStdout(), "Current version is the latest.")
			return nil
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Updating to %s...\n", latest.Version)
		exe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("could not locate executable path: %w", err)
		}

		if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
			return fmt.Errorf("error occurred while updating binary: %w", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Successfully updated to version %s\n", latest.Version)
		return nil
	},
}

func init() {
	updateCmd.Flags().Bool("preview", false, "include pre-release/beta versions")
	_ = updateCmd.Flags().MarkHidden("preview")
	rootCmd.AddCommand(updateCmd)
}
