package cmd

import (
	"github.com/spf13/cobra"
)

var ec2Cmd = &cobra.Command{
	Use: "ec2",
	//Short: "",
	//Long: ``,
}

func init() {
	RootCmd.AddCommand(ec2Cmd)
}
