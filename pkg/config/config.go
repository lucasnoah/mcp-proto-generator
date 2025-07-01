package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the complete configuration for mcp-proto-gen
type Config struct {
	Proto    ProtoConfig    `yaml:"proto"`
	Generate GenerateConfig `yaml:"generate"`
	Auth     AuthConfig     `yaml:"auth"`
	Types    TypeConfig     `yaml:"types"`
	Server   ServerConfig   `yaml:"server"`
	Services map[string]ServiceConfig `yaml:"services"`
}

// ProtoConfig configures proto file discovery
type ProtoConfig struct {
	Dirs        []string `yaml:"dirs"`
	ImportPaths []string `yaml:"import_paths"`
	Exclude     []string `yaml:"exclude"`
}

// GenerateConfig configures code generation
type GenerateConfig struct {
	OutputDir      string             `yaml:"output_dir"`
	Module         string             `yaml:"module"`
	ServiceNaming  ServiceNamingConfig `yaml:"service_naming"`
	Methods        MethodConfig       `yaml:"methods"`
}

// ServiceNamingConfig configures service endpoint naming
type ServiceNamingConfig struct {
	Style         string `yaml:"style"`
	PortStart     int    `yaml:"port_start"`
	PortIncrement int    `yaml:"port_increment"`
}

// MethodConfig configures method filtering
type MethodConfig struct {
	DangerousPatterns []string `yaml:"dangerous_patterns"`
	ExcludePatterns   []string `yaml:"exclude_patterns"`
}

// AuthConfig configures authentication
type AuthConfig struct {
	Enabled bool   `yaml:"enabled"`
	EnvVar  string `yaml:"env_var"`
}

// TypeConfig configures type handling
type TypeConfig struct {
	BinarySizeLimit int    `yaml:"binary_size_limit"`
	AnyHandling     string `yaml:"any_handling"`
}

// ServerConfig configures the generated server
type ServerConfig struct {
	Port           int  `yaml:"port"`
	HealthEnabled  bool `yaml:"health_enabled"`
	MetricsEnabled bool `yaml:"metrics_enabled"`
}

// ServiceConfig configures individual services
type ServiceConfig struct {
	Description     string `yaml:"description"`
	EndpointEnv     string `yaml:"endpoint_env"`
	DefaultEndpoint string `yaml:"default_endpoint"`
}

// Load loads configuration from file or returns defaults
func Load(configFile string) (*Config, error) {
	// Start with defaults
	cfg := defaultConfig()

	// If no config file specified, look for default
	if configFile == "" {
		configFile = "mcp-gen.yaml"
	}

	// Check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// No config file, use defaults
		return cfg, nil
	}

	// Read config file
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configFile, err)
	}

	// Parse YAML
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configFile, err)
	}

	// Validate and set defaults
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// defaultConfig returns default configuration
func defaultConfig() *Config {
	return &Config{
		Proto: ProtoConfig{
			Dirs: []string{"./proto"},
			ImportPaths: []string{
				"./proto",
				"./vendor/proto",
				filepath.Join(os.Getenv("GOPATH"), "src"),
			},
			Exclude: []string{
				"**/test*.proto",
				"**/internal*.proto",
			},
		},
		Generate: GenerateConfig{
			OutputDir: "./mcp-server",
			Module:    "github.com/example/mcp-server",
			ServiceNaming: ServiceNamingConfig{
				Style:         "kebab-case",
				PortStart:     50051,
				PortIncrement: 1,
			},
			Methods: MethodConfig{
				DangerousPatterns: []string{
					"Delete*",
					"Drop*",
					"Truncate*",
					"Purge*",
					"Clear*",
					"Remove*",
				},
				ExcludePatterns: []string{
					"*Internal",
					"Debug*",
					"Health*",
				},
			},
		},
		Auth: AuthConfig{
			Enabled: true,
			EnvVar:  "MCP_API_KEYS",
		},
		Types: TypeConfig{
			BinarySizeLimit: 10 * 1024 * 1024, // 10MB
			AnyHandling:     "json_object",
		},
		Server: ServerConfig{
			Port:           3333,
			HealthEnabled:  true,
			MetricsEnabled: false,
		},
		Services: make(map[string]ServiceConfig),
	}
}

// validate validates the configuration
func (c *Config) validate() error {
	if len(c.Proto.Dirs) == 0 {
		return fmt.Errorf("proto.dirs cannot be empty")
	}

	if c.Generate.OutputDir == "" {
		return fmt.Errorf("generate.output_dir cannot be empty")
	}

	if c.Generate.Module == "" {
		return fmt.Errorf("generate.module cannot be empty")
	}

	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("server.port must be between 1 and 65535")
	}

	validStyles := []string{"kebab-case", "snake_case", "PascalCase"}
	validStyle := false
	for _, style := range validStyles {
		if c.Generate.ServiceNaming.Style == style {
			validStyle = true
			break
		}
	}
	if !validStyle {
		return fmt.Errorf("generate.service_naming.style must be one of: %v", validStyles)
	}

	return nil
}

// GetServiceEndpoint returns the endpoint for a service name
func (c *Config) GetServiceEndpoint(serviceName string) (envVar, defaultAddr string) {
	// Check if service has custom config
	if svcCfg, exists := c.Services[serviceName]; exists {
		return svcCfg.EndpointEnv, svcCfg.DefaultEndpoint
	}

	// Generate based on conventions
	envVar = c.serviceNameToEnvVar(serviceName)
	defaultAddr = c.serviceNameToDefaultAddr(serviceName)
	return envVar, defaultAddr
}

// serviceNameToEnvVar converts service name to environment variable
func (c *Config) serviceNameToEnvVar(serviceName string) string {
	// Convert to UPPER_SNAKE_CASE
	result := ""
	for i, r := range serviceName {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result += "_"
		}
		if r >= 'a' && r <= 'z' {
			result += string(r - 'a' + 'A')
		} else if r >= 'A' && r <= 'Z' {
			result += string(r)
		} else {
			result += string(r)
		}
	}
	return result + "_ADDR"
}

// serviceNameToDefaultAddr converts service name to default address
func (c *Config) serviceNameToDefaultAddr(serviceName string) string {
	var name string
	
	switch c.Generate.ServiceNaming.Style {
	case "kebab-case":
		name = toKebabCase(serviceName)
	case "snake_case":
		name = toSnakeCase(serviceName)
	case "PascalCase":
		name = serviceName
	default:
		name = toKebabCase(serviceName)
	}

	// Simple port assignment - should be smarter in real implementation
	port := c.Generate.ServiceNaming.PortStart
	return fmt.Sprintf("%s:%d", name, port)
}

// toKebabCase converts PascalCase to kebab-case
func toKebabCase(s string) string {
	result := ""
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result += "-"
		}
		if r >= 'A' && r <= 'Z' {
			result += string(r - 'A' + 'a')
		} else {
			result += string(r)
		}
	}
	return result
}

// toSnakeCase converts PascalCase to snake_case
func toSnakeCase(s string) string {
	result := ""
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result += "_"
		}
		if r >= 'A' && r <= 'Z' {
			result += string(r - 'A' + 'a')
		} else {
			result += string(r)
		}
	}
	return result
}