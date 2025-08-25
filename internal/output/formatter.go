package output

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/tommyalmeida/envsync/internal/env"
)

type Formatter struct {
    useColor bool
    
    red    func(a ...interface{}) string
    green  func(a ...interface{}) string
    yellow func(a ...interface{}) string
    blue   func(a ...interface{}) string
    bold   func(a ...interface{}) string
}

func NewFormatter(useColor bool) *Formatter {
    f := &Formatter{useColor: useColor}
    
    if useColor {
        f.red = color.New(color.FgRed).SprintFunc()
        f.green = color.New(color.FgGreen).SprintFunc()
        f.yellow = color.New(color.FgYellow).SprintFunc()
        f.blue = color.New(color.FgBlue).SprintFunc()
        f.bold = color.New(color.Bold).SprintFunc()
    } else {
        f.red = fmt.Sprint
        f.green = fmt.Sprint
        f.yellow = fmt.Sprint
        f.blue = fmt.Sprint
        f.bold = fmt.Sprint
    }
    
    return f
}

func (f *Formatter) PrintValidationResult(result env.ValidationResult) error {
    if result.Valid {
        fmt.Println(f.green("✓ Validation passed"))
        return nil
    }
    
    fmt.Println(f.red("✗ Validation failed"))
    
    if len(result.Missing) > 0 {
        fmt.Printf("\n%s:\n", f.bold("Missing required variables"))
        for _, variable := range result.Missing {
            fmt.Printf("  - %s\n", f.red(variable))
        }
    }
    
    if len(result.Errors) > 0 {
        fmt.Printf("\n%s:\n", f.bold("Validation errors"))
        for _, err := range result.Errors {
            fmt.Printf("  - %s: %s\n", f.red(err.Variable), err.Message)
        }
    }
    
    if len(result.Extra) > 0 {
        fmt.Printf("\n%s:\n", f.bold("Extra variables (not in schema)"))
        for _, variable := range result.Extra {
            fmt.Printf("  - %s\n", f.yellow(variable))
        }
    }
    
    os.Exit(1)
    return nil
}

func (f *Formatter) PrintDiff(diff env.DiffResult, sourceFile, targetFile string) error {
    fmt.Printf("%s vs %s\n\n", f.bold(sourceFile), f.bold(targetFile))
    
    hasChanges := len(diff.Missing) > 0 || len(diff.Extra) > 0 || len(diff.Different) > 0
    
    if !hasChanges {
        fmt.Println(f.green("✓ Files are in sync"))
        return nil
    }
    
    if len(diff.Missing) > 0 {
        fmt.Printf("%s (%d):\n", f.bold("Missing in target"), len(diff.Missing))
        for _, key := range diff.Missing {
            fmt.Printf("  %s %s\n", f.red("-"), key)
        }
        fmt.Println()
    }
    
    if len(diff.Extra) > 0 {
        fmt.Printf("%s (%d):\n", f.bold("Extra in target"), len(diff.Extra))
        for _, key := range diff.Extra {
            fmt.Printf("  %s %s\n", f.green("+"), key)
        }
        fmt.Println()
    }
    
    if len(diff.Different) > 0 {
        fmt.Printf("%s (%d):\n", f.bold("Different values"), len(diff.Different))
        for key, values := range diff.Different {
            fmt.Printf("  %s %s\n", f.yellow("~"), key)
            fmt.Printf("    %s: %s\n", f.blue("source"), values.Source)
            fmt.Printf("    %s: %s\n", f.blue("target"), values.Target)
        }
    }
    
    return nil
}

func (f *Formatter) PrintSyncResult(result env.SyncResult, dryRun bool) error {
    action := "Synced"
    if dryRun {
        action = "Would sync"
    }
    
    if len(result.Added) == 0 {
        fmt.Println(f.green("✓ No variables need to be synced"))
        return nil
    }
    
    fmt.Printf("%s %d variables to %s:\n\n", action, len(result.Added), f.bold(result.FilePath))
    
    for _, key := range result.Added {
        fmt.Printf("  %s %s\n", f.green("+"), key)
    }
    
    if dryRun {
        fmt.Printf("\n%s\n", f.yellow("This was a dry run. Use --dry-run=false to apply changes."))
    } else {
        fmt.Printf("\n%s\n", f.green("✓ Sync completed successfully"))
    }
    
    return nil
}