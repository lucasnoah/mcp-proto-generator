package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lucasnoah/mcp-proto-generator/pkg/config"
	"github.com/lucasnoah/mcp-proto-generator/pkg/generator"
	"github.com/lucasnoah/mcp-proto-generator/pkg/parser"
)

var (
	protoDir  string
	outputDir string
	dryRun    bool
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate MCP server from proto files",
	Long: `Generate a complete MCP server from your gRPC proto definitions.

This command will:
1. Parse all proto files in the specified directory
2. Extract services, methods, and message types
3. Generate a complete Go MCP server
4. Include API key authentication
5. Create Docker deployment files

Examples:
  # Generate from current directory
  mcp-proto-gen generate

  # Specify proto directory
  mcp-proto-gen generate --proto-dir ./api/proto

  # Custom output location
  mcp-proto-gen generate --output ./my-mcp-server

  # Dry run (show what would be generated)
  mcp-proto-gen generate --dry-run`,
	RunE: runGenerate,
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Flags
	generateCmd.Flags().StringVar(&protoDir, "proto-dir", ".", "directory containing proto files")
	generateCmd.Flags().StringVar(&outputDir, "output", "", "output directory (default from config)")
	generateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be generated without creating files")

	// Bind flags to viper
	viper.BindPFlag("proto.dirs", generateCmd.Flags().Lookup("proto-dir"))
	viper.BindPFlag("generate.output_dir", generateCmd.Flags().Lookup("output"))
}

func runGenerate(cmd *cobra.Command, args []string) error {
	verbose := viper.GetBool("verbose")
	
	if verbose {
		fmt.Println("🚀 Starting MCP server generation...")
	}

	// Load configuration
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Override config with command line flags
	if protoDir != "." {
		cfg.Proto.Dirs = []string{protoDir}
	}
	if outputDir != "" {
		cfg.Generate.OutputDir = outputDir
	}

	if verbose {
		fmt.Printf("📂 Proto directories: %v\n", cfg.Proto.Dirs)
		fmt.Printf("📁 Output directory: %s\n", cfg.Generate.OutputDir)
	}

	// Validate proto directories exist
	for _, dir := range cfg.Proto.Dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return fmt.Errorf("proto directory does not exist: %s", dir)
		}
	}

	// Parse proto files
	start := time.Now()
	if verbose {
		fmt.Println("🔍 Parsing proto files...")
	}

	p := parser.New(cfg.Proto.ImportPaths)
	services, err := p.ParseDirectories(cfg.Proto.Dirs, cfg.Proto.Exclude)
	if err != nil {
		return fmt.Errorf("failed to parse proto files: %w", err)
	}

	parseTime := time.Since(start)
	if verbose {
		fmt.Printf("✅ Parsed %d services in %v\n", len(services), parseTime)
		for _, svc := range services {
			fmt.Printf("   📋 %s (%d methods)\n", svc.Name, len(svc.Methods))
		}
	}

	// Generate MCP server
	if verbose {
		fmt.Println("⚙️  Generating MCP server...")
	}

	gen := generator.New(cfg, services)

	if dryRun {
		// Show what would be generated
		files, err := gen.Plan()
		if err != nil {
			return fmt.Errorf("failed to plan generation: %w", err)
		}

		fmt.Println("📝 Files that would be generated:")
		for _, file := range files {
			fmt.Printf("   📄 %s\n", file.Path)
		}
		return nil
	}

	// Actually generate
	files, err := gen.Generate()
	if err != nil {
		return fmt.Errorf("failed to generate MCP server: %w", err)
	}

	genTime := time.Since(start) - parseTime
	if verbose {
		fmt.Printf("✅ Generated %d files in %v\n", len(files), genTime)
	}

	// Write files
	if verbose {
		fmt.Println("💾 Writing files...")
	}

	for _, file := range files {
		path := filepath.Join(cfg.Generate.OutputDir, file.Path)
		
		// Create directory if needed
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

		// Write file
		if err := os.WriteFile(path, file.Content, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", path, err)
		}

		if verbose {
			fmt.Printf("   📄 %s\n", path)
		}
	}

	totalTime := time.Since(start)

	// Success message
	fmt.Printf("\n🎉 MCP server generated successfully!\n")
	fmt.Printf("📊 Stats:\n")
	fmt.Printf("   • Services: %d\n", len(services))
	
	totalMethods := 0
	for _, svc := range services {
		totalMethods += len(svc.Methods)
	}
	fmt.Printf("   • Methods: %d\n", totalMethods)
	fmt.Printf("   • Files: %d\n", len(files))
	fmt.Printf("   • Time: %v\n", totalTime)

	fmt.Printf("\n🏃 Next steps:\n")
	fmt.Printf("   cd %s\n", cfg.Generate.OutputDir)
	fmt.Printf("   go mod download\n")
	fmt.Printf("   go build .\n")
	fmt.Printf("   ./mcp-server\n")

	return nil
}