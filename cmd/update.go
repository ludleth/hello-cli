package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

const repo = "ludleth/hello-cli"

var updateCmd = &cobra.Command{
	Use:   "update [version]",
	Short: "Update hello-cli to the latest or specified version",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		preview, _ := cmd.Flags().GetBool("preview")
		yes, _ := cmd.Flags().GetBool("yes")

		if Version == "dev" {
			fmt.Fprintln(cmd.ErrOrStderr(), "Running development version. Skipping update.")
			return nil
		}

		var targetVersion string
		if len(args) > 0 {
			targetVersion = args[0]
			if !strings.HasPrefix(targetVersion, "v") {
				targetVersion = "v" + targetVersion
			}
		}

		var latest *selfupdate.Release
		var found bool
		var err error

		if targetVersion != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "Looking for version %s of %s...\n", targetVersion, repo)
			latest, found, err = selfupdate.DetectVersion(repo, targetVersion)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "Checking for latest version of %s...\n", repo)
			latest, found, err = selfupdate.DetectLatest(repo)
		}
		if err != nil {
			return fmt.Errorf("error occurred while detecting version: %w", err)
		}
		if !found {
			return fmt.Errorf("version not found on GitHub")
		}

		if targetVersion == "" && !preview && len(latest.Version.Pre) > 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "Latest version found (%s) is a pre-release. Skipping.\n", latest.Version)
			return nil
		}

		currentVersion := strings.TrimPrefix(Version, "v")
		if latest.Version.String() == currentVersion {
			fmt.Fprintln(cmd.OutOrStdout(), "Current version is the latest.")
			return nil
		}

		if !yes {
			fmt.Fprintf(cmd.OutOrStdout(), "Update to %s? [y/N]: ", latest.Version)
			reader := bufio.NewReader(cmd.InOrStdin())
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(strings.ToLower(answer))
			if answer != "y" && answer != "yes" {
				fmt.Fprintln(cmd.OutOrStdout(), "Update cancelled.")
				return nil
			}
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
	updateCmd.Flags().BoolP("yes", "y", false, "skip confirmation prompt")
	rootCmd.AddCommand(updateCmd)
}
