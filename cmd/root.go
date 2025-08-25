package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tommyalmeida/envsync/internal/config"
	"github.com/tommyalmeida/envsync/internal/env"
	"github.com/tommyalmeida/envsync/internal/output"
)

var (
	cfgFile    string
	jsonOutput bool
	verbose    bool
)

var rootCmd = &cobra.Command{
	Use:   "envsync",
	Short: "Keep environment variable files consistent across environments",
}

var validateCmd = &cobra.Command{
	Use:   "validate [env-file]",
	Short: "Validate an environment file against a schema",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		return runValidate(args[0])
	},
}

var diffCmd = &cobra.Command{
	Use:   "diff [source-env] [target-env]",
	Short: "Compare two environment files and show differences",
	Args:  cobra.ExactArgs(2),
	RunE: func(_ *cobra.Command, args []string) error {
		return runDiff(args[0], args[1])
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync [source-env] [target-env]",
	Short: "Synchronize missing variables from source to target",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		return runSync(args[0], args[1], dryRun)
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .envsync.yaml)")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "output in JSON format")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	syncCmd.Flags().Bool("dry-run", false, "show what would be synced without making changes")

	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(syncCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".envsync")
		viper.SetConfigType("yaml")

		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		if verbose {
			log.Printf("Warning: No config file found: %v\n", err)
		}
	}
}

func runDiff(sourceFile, targetFile string) error {
	sourceVars, err := env.ParseFile(sourceFile)

	if err != nil {
		return fmt.Errorf("failed to parse source file: %w", err)
	}

	targetVars, err := env.ParseFile(targetFile)

	if err != nil {
		return fmt.Errorf("failed to parse target file: %w", err)
	}

	diff := env.CompareEnvs(sourceVars, targetVars)

	if jsonOutput {
		return outputJSON(diff)
	}

	formatter := output.NewFormatter(!jsonOutput)
	return formatter.PrintDiff(diff, sourceFile, targetFile)
}

func runSync(sourceFile, targetFile string, dryRun bool) error {
	sourceVars, err := env.ParseFile(sourceFile)

	if err != nil {
		return fmt.Errorf("failed to parse source file: %w", err)
	}

	targetVars, err := env.ParseFile(targetFile)

	if err != nil {
		return fmt.Errorf("failed to parse target file: %w", err)
	}

	cfg, err := config.Load()

	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	syncer := env.NewSyncer(cfg)
	result, err := syncer.Sync(sourceVars, targetVars, targetFile, dryRun)

	if err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}

	if jsonOutput {
		return outputJSON(result)
	}

	formatter := output.NewFormatter(!jsonOutput)
	return formatter.PrintSyncResult(result, dryRun)
}

func runValidate(envFile string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	envVars, err := env.ParseFile(envFile)

	if err != nil {
		return fmt.Errorf("failed to parse env file: %w", err)
	}

	validator := env.NewValidator(cfg.Schema)
	result := validator.Validate(envVars)

	if jsonOutput {
		return outputJSON(result)
	}

	formatter := output.NewFormatter(!jsonOutput)
	return formatter.PrintValidationResult(result)
}

func outputJSON(v any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	return encoder.Encode(v)
}

func Execute() error {
	return rootCmd.Execute()
}
