package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"runtime"

	"github.com/getcarina/carina/common"
	"github.com/spf13/cobra"
)

func newEnvCommand() *cobra.Command {
	var options struct {
		name  string
		shell string
		path  string
	}

	var cmd = &cobra.Command{
		Use:               "env <cluster-name>",
		Short:             "Show the command to connect docker/kubectl to a cluster",
		Long:              "Show the command to connect docker/kubectl to a cluster by setting environment variables in the current shell session",
		PersistentPreRunE: authenticatedPreRunE,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if options.shell == "" {
				shell := os.Getenv("SHELL")
				if shell != "" {
					options.shell = filepath.Base(shell)
					common.Log.WriteDebug("Shell: SHELL (%s)", options.shell)
				} else {
					options.shell = detectShell()
					if options.shell != "" {
						common.Log.WriteDebug("Shell: detected (%s)", options.shell)
					} else {
						return errors.New("Shell was not specified. Either use --shell or set SHELL")
					}
				}
			} else {
				common.Log.WriteDebug("Shell: --shell (%s)", options.shell)
			}

			return bindClusterNameArg(args, &options.name)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			sourceText, err := cxt.Client.GetSourceCommand(cxt.Account, options.shell, options.name, options.path)
			if err != nil {
				return err
			}

			fmt.Println(sourceText)
			return nil
		},
	}

	cmd.ValidArgs = []string{"cluster-name"}
	cmd.Flags().StringVar(&options.shell, "shell", "", "The parent shell type. Allowed values: bash, fish, powershell, cmd [SHELL]")
	cmd.Flags().StringVar(&options.path, "path", "", "Full path to the directory from which the credentials should be loaded")
	cmd.SetUsageTemplate(cmd.UsageTemplate())

	return cmd
}

func detectShell() string {
	if runtime.GOOS != "windows" {
		return ""
	}

	common.Log.WriteDebug("Detecting --shell")
	// Dirty hack: CMD seems to have PROMPT set, while PowerShell doesn't
	if os.Getenv("PROMPT") != "" {
		return "cmd"
	}

	return "powershell"
}
