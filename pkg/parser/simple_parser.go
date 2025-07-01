package parser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// SimpleParser is a basic proto parser that extracts services using regex
// This is a simplified approach for the MVP - full proto parsing would use protoc
type SimpleParser struct {
	ImportPaths []string
}

// NewSimple creates a new simple parser
func NewSimple(importPaths []string) *SimpleParser {
	return &SimpleParser{
		ImportPaths: importPaths,
	}
}

// ParseDirectoriesSimple parses proto files using simple regex matching
func (p *SimpleParser) ParseDirectoriesSimple(dirs []string, excludePatterns []string) ([]ServiceDefinition, error) {
	var protoFiles []string

	// Find all proto files
	for _, dir := range dirs {
		files, err := findProtoFiles(dir, excludePatterns)
		if err != nil {
			return nil, fmt.Errorf("failed to find proto files in %s: %w", dir, err)
		}
		protoFiles = append(protoFiles, files...)
	}

	if len(protoFiles) == 0 {
		return nil, fmt.Errorf("no proto files found in directories: %v", dirs)
	}

	// Parse each file
	var services []ServiceDefinition
	for _, file := range protoFiles {
		fileSvcs, err := p.parseProtoFileSimple(file)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", file, err)
		}
		services = append(services, fileSvcs...)
	}

	return services, nil
}

// parseProtoFileSimple parses a single proto file using regex
func (p *SimpleParser) parseProtoFileSimple(filename string) ([]ServiceDefinition, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var services []ServiceDefinition
	var currentService *ServiceDefinition
	var packageName string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") || line == "" {
			continue
		}

		// Extract package name
		if strings.HasPrefix(line, "package ") {
			packageRegex := regexp.MustCompile(`package\s+([^;]+);`)
			if matches := packageRegex.FindStringSubmatch(line); len(matches) > 1 {
				packageName = strings.TrimSpace(matches[1])
			}
		}

		// Start of service
		if strings.HasPrefix(line, "service ") {
			serviceRegex := regexp.MustCompile(`service\s+(\w+)\s*{`)
			if matches := serviceRegex.FindStringSubmatch(line); len(matches) > 1 {
				if currentService != nil {
					services = append(services, *currentService)
				}
				currentService = &ServiceDefinition{
					Name:    matches[1],
					Package: packageName,
					File:    filename,
					Methods: []MethodDefinition{},
				}
			}
		}

		// RPC method
		if currentService != nil && strings.Contains(line, "rpc ") {
			rpcRegex := regexp.MustCompile(`rpc\s+(\w+)\s*\(\s*(\w+)\s*\)\s*returns\s*\(\s*(\w+)\s*\)`)
			if matches := rpcRegex.FindStringSubmatch(line); len(matches) > 3 {
				method := MethodDefinition{
					Name: matches[1],
					InputType: MessageDefinition{
						Name: matches[2],
					},
					OutputType: MessageDefinition{
						Name: matches[3],
					},
					ClientStreaming: false, // Simplified - would need stream detection
					ServerStreaming: false,
					Description:     "", // Would extract from comments
				}
				currentService.Methods = append(currentService.Methods, method)
			}
		}

		// End of service
		if currentService != nil && line == "}" {
			services = append(services, *currentService)
			currentService = nil
		}
	}

	// Handle case where file ends without closing brace
	if currentService != nil {
		services = append(services, *currentService)
	}

	return services, scanner.Err()
}

// findProtoFiles recursively finds all .proto files (moved from original parser)
func findProtoFiles(dir string, excludePatterns []string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".proto") {
			return nil
		}

		// Check exclude patterns
		for _, pattern := range excludePatterns {
			if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
				return nil
			}
		}

		files = append(files, path)
		return nil
	})

	return files, err
}