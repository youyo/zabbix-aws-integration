package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	latest "github.com/tcnksm/go-latest"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	//Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := execVersion(os.Stdout, Version, CommitHash, BuildTime, GoVersion); err != nil {
			log.Fatal(err)
		}
	},
}

func execVersion(w io.Writer, version, commitHash, buildTime, goVersion string) (err error) {
	outdate, currentVersion, err := versionCheck(version)
	if err != nil {
		return
	}
	fmt.Fprintf(w, "Version: %s\nCommitHash: %s\nBuildTime: %s\nGoVersion: %s\n", version, commitHash, buildTime, goVersion)
	if outdate {
		fmt.Fprintf(w, "\n%s is not latest, you should upgrade to %s\n", version, currentVersion)
	}
	return
}

func versionCheck(version string) (outdate bool, currentVersion string, err error) {
	githubTag := &latest.GithubTag{
		Owner:      "youyo",
		Repository: Name,
	}
	res, err := latest.Check(githubTag, version)
	if err == nil {
		if res.Outdated {
			outdate = true
			currentVersion = res.Current
			return
		}
	} else {
		return
	}
	return
}

func init() {
	RootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}