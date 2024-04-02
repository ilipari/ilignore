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
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ilignore",
	Short: "CLI Tools to work with .gitignore files",
	Long: `CLI Tools to help prevent further commits in files that are already tracked in a git repo.
Ignored files specified in .gitignore file format.`,
	Version: "0.1",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var programLogLevel = new(slog.LevelVar)

func initLogs(level slog.Level) {
	programLogLevel.Set(level)
	// h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	// slog.SetDefault(slog.New(h))
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: programLogLevel}))
	slog.SetDefault(logger)
}

func init() {
	initLogs(slog.LevelWarn)
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ilignorerc.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output (same as --log info)")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.PersistentFlags().String("log", "", "log level: one of Debug, Info, Warn, Error. Case insensitive (default is WARN)")
	viper.BindPFlag("log", rootCmd.PersistentFlags().Lookup("log"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// look for config in the working directory
		viper.AddConfigPath(".")

		// optionally look for config in home directory
		viper.AddConfigPath(home)

		// with name ".ilignorerc" (with or without extension).
		viper.SetConfigType("yaml")
		viper.SetConfigName(".ilignorerc")
	}

	viper.AutomaticEnv() // read in environment variables that match

	err := viper.ReadInConfig()
	setProgramLogLevelFromViper()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error
			slog.Warn("Config file not found")
		} else {
			// Config file was found but another error was produced
			slog.Error("Error Using config file: ", err)
			os.Exit(1)
		}
	} else {
		// Config file found and successfully parsed
		slog.Info("Using config: ", "file", viper.ConfigFileUsed())
	}

}

func setProgramLogLevelFromViper() {
	verbose := viper.GetBool("verbose")
	if verbose {
		programLogLevel.Set(slog.LevelInfo)
	}
	logLevel, logErr := parseLogLevel(viper.GetString("log"))
	if logErr == nil {
		programLogLevel.Set(logLevel)
	}
}

func parseLogLevel(s string) (slog.Level, error) {
	var level slog.Level
	var err = level.UnmarshalText([]byte(s))
	return level, err
}
