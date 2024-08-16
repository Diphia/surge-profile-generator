package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/template"
	"gopkg.in/yaml.v3"
)

// interfaceSlice converts an interface{} to a []interface{} if possible
func interfaceSlice(slice interface{}) ([]interface{}, bool) {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil, false
	}
	ret := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}
	return ret, true
}

// joinFunc is a template function that joins slice elements
func joinFunc(sep string, slice interface{}) string {
	if strSlice, ok := slice.([]string); ok {
		return strings.Join(strSlice, sep)
	}
	if interfaceSlice, ok := interfaceSlice(slice); ok {
		strSlice := make([]string, len(interfaceSlice))
		for i, v := range interfaceSlice {
			strSlice[i] = fmt.Sprint(v)
		}
		return strings.Join(strSlice, sep)
	}
	return fmt.Sprint(slice)
}

// mergeValues merges multiple maps, with later maps taking precedence
func mergeValues(valuesList ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, values := range valuesList {
		for k, v := range values {
			if existing, ok := result[k]; ok {
				if existingMap, ok := existing.(map[string]interface{}); ok {
					if newMap, ok := v.(map[string]interface{}); ok {
						// Recursively merge nested maps
						result[k] = mergeValues(existingMap, newMap)
						continue
					}
				}
				if existingSlice, ok := existing.([]interface{}); ok {
					if newSlice, ok := v.([]interface{}); ok {
						// Merge slices by prepending new elements
						result[k] = append(newSlice, existingSlice...)
						continue
					}
				}
			}
			// For non-map and non-slice types or if the key doesn't exist, simply override
			result[k] = v
		}
	}
	return result
}

func main() {
	// Define command-line flags
	templateFile := flag.String("template", "", "Path to the template file")
	outputFile := flag.String("output", "", "Path to the output file (optional, defaults to stdout)")
	
	// Parse command-line flags
	flag.Parse()

	// Check if required flags are provided
	if *templateFile == "" || flag.NArg() == 0 {
		fmt.Println("Usage: go run main.go -template <template_file> [-output <output_file>] <values_file1> [<values_file2> ...]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Get the list of value files from remaining arguments
	valueFiles := flag.Args()

	// Read the template file
	templateContent, err := os.ReadFile(*templateFile)
	if err != nil {
		fmt.Printf("Error reading template file: %v\n", err)
		os.Exit(1)
	}

	// Read and merge all value files
	var allValues []map[string]interface{}
	for _, file := range valueFiles {
		valuesContent, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading values file %s: %v\n", file, err)
			os.Exit(1)
		}
		var values map[string]interface{}
		err = yaml.Unmarshal(valuesContent, &values)
		if err != nil {
			fmt.Printf("Error parsing YAML from file %s: %v\n", file, err)
			os.Exit(1)
		}
		allValues = append(allValues, values)
	}

	// Merge all values, with later files taking precedence
	mergedValues := mergeValues(allValues...)

	// Create a new template and parse the content
	tmpl, err := template.New("config").Funcs(template.FuncMap{
		"join": joinFunc,
	}).Parse(string(templateContent))
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		os.Exit(1)
	}

	// Determine the output destination
	var output *os.File
	if *outputFile != "" {
		output, err = os.Create(*outputFile)
		if err != nil {
			fmt.Printf("Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer output.Close()
	} else {
		output = os.Stdout
	}

	// Execute the template with the merged values
	err = tmpl.Execute(output, mergedValues)
	if err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		os.Exit(1)
	}

	if *outputFile != "" {
		fmt.Printf("Configuration written to %s\n", *outputFile)
	}
}
