/*
Copyright © 2024 Ignazio Lipari iglipari@gmail.com

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"log/slog"

	"ilipari/ilignore/service"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Checks which files conflicts with an ignore file",
	Long: `Analyzes a list of files looking for conflicts with an ignore file in .gitignore format. 
The list of files to check can be obtained through the execution of a shell command or via stdin.
A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("check command called")
		var s = service.NewService(viper.GetStringSlice(configKey(cmd, IGNORE_FILE_FLAG)), false)
		filesCh := service.NewFileSource(args, viper.GetString(configKey(cmd, LIST_FILES_FLAG)), true)
		// filesCh := service.NewCommandFileSource(viper.GetString(configKey(cmd, LIST_FILES_FLAG)), true)
		// filesCh := service.NewGitDiffFileSource(true, "")
		// filesCh := service.NewStdinFileSource()
		// filesCh := service.NewFixedFileSource([]string{"ciao.txt", "mondo.csv", ".vscode"})
		conflictsChannel := s.CheckFiles(filesCh)
		conflictsConsumerOutput := service.NewConsoleConflictConsumer(conflictsChannel, "")
		for err := range conflictsConsumerOutput {
			slog.Error("error!", "error", err)
		}
	},
}

const LIST_FILES_FLAG, IGNORE_FILE_FLAG = "files", "ignore"

func init() {
	rootCmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	checkCmd.Flags().StringP(LIST_FILES_FLAG, "f", "", "command to obtain list of files to check")
	viper.BindPFlag(configKey(checkCmd, LIST_FILES_FLAG), checkCmd.Flags().Lookup(LIST_FILES_FLAG))

	checkCmd.Flags().StringSliceP(IGNORE_FILE_FLAG, "i", []string{service.IGNORE_FILE}, "Ignore file")
	viper.BindPFlag(configKey(checkCmd, IGNORE_FILE_FLAG), checkCmd.Flags().Lookup(IGNORE_FILE_FLAG))

}

func configKey(cmd *cobra.Command, flag string) string {
	if cmd.HasParent() {
		return configKey(cmd.Parent(), "") + cmd.Name() + "." + flag
	}
	return ""
}
