package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

// Reference represents a JSON reference
type Reference struct {
	Ref string `json:"$ref"`
}

// isReference checks if a map contains a $ref key
func isReference(m map[string]interface{}) bool {
	_, ok := m["$ref"]
	return ok
}

// resolveReference resolves a reference path and returns the referenced content
func resolveReference(refPath string, rootDoc map[string]interface{}, baseDir string) (interface{}, error) {
	// Handle internal references (starting with #)
	if strings.HasPrefix(refPath, "#/") {
		parts := strings.Split(strings.TrimPrefix(refPath, "#/"), "/")
		current := rootDoc
		
		for _, part := range parts {
			if v, ok := current[part]; ok {
				if m, ok := v.(map[string]interface{}); ok {
					current = m
				} else {
					return v, nil
				}
			} else {
				return nil, fmt.Errorf("reference path not found: %s", refPath)
			}
		}
		return current, nil
	}

	// Handle external references (files)
	filePath := filepath.Join(baseDir, refPath)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read referenced file %s: %v", filePath, err)
	}

	var referencedDoc interface{}
	if err := json.Unmarshal(content, &referencedDoc); err != nil {
		return nil, fmt.Errorf("failed to parse referenced file %s: %v", filePath, err)
	}

	return referencedDoc, nil
}

// inlineReferences recursively inlines all references in the document
func inlineReferences(doc interface{}, rootDoc map[string]interface{}, baseDir string) (interface{}, error) {
	switch v := doc.(type) {
	case map[string]interface{}:
		if isReference(v) {
			refPath := v["$ref"].(string)
			resolved, err := resolveReference(refPath, rootDoc, baseDir)
			if err != nil {
				return nil, err
			}
			return inlineReferences(resolved, rootDoc, baseDir)
		}

		result := make(map[string]interface{})
		for key, value := range v {
			inlined, err := inlineReferences(value, rootDoc, baseDir)
			if err != nil {
				return nil, err
			}
			result[key] = inlined
		}
		return result, nil

	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			inlined, err := inlineReferences(item, rootDoc, baseDir)
			if err != nil {
				return nil, err
			}
			result[i] = inlined
		}
		return result, nil

	default:
		return v, nil
	}
}

func main() {
	// Parse command line arguments
	inputFile := flag.String("input", "", "Path to the OpenAPI JSON file")
	outputFile := flag.String("output", "", "Path to the output JSON file (optional)")
	flag.Parse()

	if *inputFile == "" {
		log.Fatal("Please provide an input file using -input flag")
	}

	// Read the input file
	content, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Failed to read input file: %v", err)
	}

	// Parse the JSON content
	var doc map[string]interface{}
	if err := json.Unmarshal(content, &doc); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Get the base directory for resolving external references
	baseDir := filepath.Dir(*inputFile)

	// Inline all references
	inlined, err := inlineReferences(doc, doc, baseDir)
	if err != nil {
		log.Fatalf("Failed to inline references: %v", err)
	}

	// Convert the result to JSON
	result, err := json.MarshalIndent(inlined, "", "  ")
	if err != nil {
		log.Fatalf("Failed to generate JSON output: %v", err)
	}

	// Write the output
	if *outputFile == "" {
		// Print to stdout if no output file is specified
		fmt.Println(string(result))
	} else {
		// Write to the specified output file
		if err := ioutil.WriteFile(*outputFile, result, 0644); err != nil {
			log.Fatalf("Failed to write output file: %v", err)
		}
		fmt.Printf("Successfully wrote inlined OpenAPI spec to %s\n", *outputFile)
	}
} 