package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// Parser parses proto files and extracts service definitions
type Parser struct {
	ImportPaths []string
}

// New creates a new parser with the given import paths
func New(importPaths []string) *Parser {
	return &Parser{
		ImportPaths: importPaths,
	}
}

// ServiceDefinition represents a parsed gRPC service
type ServiceDefinition struct {
	Name    string
	Package string
	Methods []MethodDefinition
	File    string
}

// MethodDefinition represents a gRPC method
type MethodDefinition struct {
	Name            string
	InputType       MessageDefinition
	OutputType      MessageDefinition
	ClientStreaming bool
	ServerStreaming bool
	Description     string
}

// MessageDefinition represents a protobuf message
type MessageDefinition struct {
	Name   string
	Fields []FieldDefinition
}

// FieldDefinition represents a message field
type FieldDefinition struct {
	Name     string
	Type     string
	Repeated bool
	Optional bool
	Map      bool
	KeyType  string
	JSONName string
}

// ParseDirectories parses all proto files in the given directories
func (p *Parser) ParseDirectories(dirs []string, excludePatterns []string) ([]ServiceDefinition, error) {
	// For now, use the simple parser to avoid protoc complications
	// TODO: Implement full protoc-based parsing later
	simple := NewSimple(p.ImportPaths)
	return simple.ParseDirectoriesSimple(dirs, excludePatterns)
}

// findProtoFiles recursively finds all .proto files in a directory
func (p *Parser) findProtoFiles(dir string, excludePatterns []string) ([]string, error) {
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

// parseProtoFiles parses the given proto files using protogen
func (p *Parser) parseProtoFiles(protoFiles []string) ([]ServiceDefinition, error) {
	// Create a minimal protoc request
	req := &pluginpb.CodeGeneratorRequest{
		FileToGenerate: protoFiles,
		Parameter:      proto.String(""),
	}

	// Read and parse each proto file
	for _, file := range protoFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", file, err)
		}

		// This is a simplified approach - in a real implementation,
		// we'd use protoc or a proper proto parser
		descriptor, err := p.parseProtoContent(file, content)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", file, err)
		}

		req.ProtoFile = append(req.ProtoFile, descriptor)
	}

	// Generate using protogen
	gen, err := protogen.Options{}.New(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create protogen: %w", err)
	}

	var services []ServiceDefinition

	// Extract services from generated files
	for _, file := range gen.Files {
		if !file.Generate {
			continue
		}

		for _, service := range file.Services {
			svcDef := ServiceDefinition{
				Name:    string(service.Desc.Name()),
				Package: string(file.Desc.Package()),
				File:    file.Desc.Path(),
			}

			// Extract methods
			for _, method := range service.Methods {
				methodDef := MethodDefinition{
					Name:            string(method.Desc.Name()),
					ClientStreaming: method.Desc.IsStreamingClient(),
					ServerStreaming: method.Desc.IsStreamingServer(),
					Description:     p.extractDescription(method.Comments),
				}

				// Extract input/output types
				methodDef.InputType = p.extractMessageDefinition(method.Input)
				methodDef.OutputType = p.extractMessageDefinition(method.Output)

				svcDef.Methods = append(svcDef.Methods, methodDef)
			}

			services = append(services, svcDef)
		}
	}

	return services, nil
}

// parseProtoContent parses proto file content into a descriptor
// This is a simplified implementation - real parser would be more robust
func (p *Parser) parseProtoContent(filename string, content []byte) (*descriptorpb.FileDescriptorProto, error) {
	// For now, return a minimal descriptor
	// In a real implementation, this would use a proper proto parser
	descriptor := &descriptorpb.FileDescriptorProto{
		Name:    proto.String(filename),
		Package: proto.String(""),
	}

	// TODO: Implement proper proto parsing
	// This is where we'd parse the actual proto syntax and extract:
	// - Package name
	// - Imports
	// - Messages
	// - Services
	// - Enums
	// etc.

	return descriptor, nil
}

// extractMessageDefinition extracts message definition from protogen message
func (p *Parser) extractMessageDefinition(msg *protogen.Message) MessageDefinition {
	msgDef := MessageDefinition{
		Name: string(msg.Desc.Name()),
	}

	// Extract fields
	for _, field := range msg.Fields {
		fieldDef := FieldDefinition{
			Name:     string(field.Desc.Name()),
			JSONName: field.Desc.JSONName(),
			Repeated: field.Desc.IsList(),
			Optional: field.Desc.HasOptionalKeyword(),
		}

		// Determine field type
		fieldDef.Type = p.protoTypeToString(field.Desc.Kind())

		// Handle map fields
		if field.Desc.IsMap() {
			fieldDef.Map = true
			mapEntry := field.Desc.MapValue()
			fieldDef.Type = p.protoTypeToString(mapEntry.Kind())
			fieldDef.KeyType = p.protoTypeToString(field.Desc.MapKey().Kind())
		}

		msgDef.Fields = append(msgDef.Fields, fieldDef)
	}

	return msgDef
}

// protoTypeToString converts protoreflect kind to string
func (p *Parser) protoTypeToString(kind protoreflect.Kind) string {
	switch kind {
	case protoreflect.BoolKind:
		return "bool"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return "int32"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "uint32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "int64"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "uint64"
	case protoreflect.FloatKind:
		return "float"
	case protoreflect.DoubleKind:
		return "double"
	case protoreflect.StringKind:
		return "string"
	case protoreflect.BytesKind:
		return "bytes"
	case protoreflect.EnumKind:
		return "enum"
	case protoreflect.MessageKind:
		return "message"
	default:
		return "unknown"
	}
}

// extractDescription extracts description from comments
func (p *Parser) extractDescription(comments protogen.CommentSet) string {
	if string(comments.Leading) != "" {
		return strings.TrimSpace(string(comments.Leading))
	}
	if string(comments.Trailing) != "" {
		return strings.TrimSpace(string(comments.Trailing))
	}
	return ""
}

// IsDangerousMethod checks if a method name matches dangerous patterns
func IsDangerousMethod(methodName string, patterns []string) bool {
	for _, pattern := range patterns {
		if matched, _ := filepath.Match(pattern, methodName); matched {
			return true
		}
	}
	return false
}

// IsExcludedMethod checks if a method should be excluded
func IsExcludedMethod(methodName string, patterns []string) bool {
	for _, pattern := range patterns {
		if matched, _ := filepath.Match(pattern, methodName); matched {
			return true
		}
	}
	return false
}