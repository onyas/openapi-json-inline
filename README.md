# OpenAPI Reference Inliner

A command-line tool written in Go that processes OpenAPI/Swagger JSON files by inlining all `$ref` references. It resolves both internal and external references, creating a single self-contained OpenAPI specification.

## Features

- Resolves internal references (e.g., `#/components/schemas/User`)
- Handles external file references (e.g., `./schemas/user.json`)
- Supports nested references (references within referenced files)
- Preserves the original OpenAPI structure
- Outputs beautifully formatted JSON with proper indentation
- Provides clear error messages for troubleshooting
- Can output to file or stdout

## Prerequisites

- Go 1.21 or higher

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/openapi-inline.git
cd openapi-inline
```

2. Build the executable:
```bash
go build -o openapi-inliner
```

## Usage

The tool provides a simple command-line interface with the following flags:

```bash
./openapi-inliner -input <input-file> [-output <output-file>]
```

### Flags
- `-input`: (Required) Path to the input OpenAPI JSON file
- `-output`: (Optional) Path where the inlined JSON should be saved. If not provided, output will be printed to stdout

### Examples

1. Process a file and print to console:
```bash
./openapi-inliner -input api.json
```

2. Process a file and save to a new file:
```bash
./openapi-inliner -input api.json -output inlined-api.json
```

3. Process a file with external references:
```bash
./openapi-inliner -input main-api.json -output combined-api.json
```

## Example Input/Output

### Input Structure
```
api/
├── main.json
├── schemas/
│   ├── user.json
│   └── product.json
```

main.json:
```json
{
  "components": {
    "schemas": {
      "User": {
        "$ref": "./schemas/user.json"
      },
      "Product": {
        "$ref": "./schemas/product.json"
      }
    }
  }
}
```

### Output
The tool will combine all referenced files into a single JSON file with resolved references.

## Error Handling

The tool provides clear error messages for common issues:

- Missing input file:
  ```
  Please provide an input file using -input flag
  ```

- Invalid JSON format:
  ```
  Failed to parse JSON: invalid character '}' looking for beginning of object key string
  ```

- Missing referenced files:
  ```
  Failed to read referenced file ./schemas/user.json: no such file or directory
  ```

- Invalid reference paths:
  ```
  Reference path not found: #/components/schemas/InvalidSchema
  ```

## Best Practices

1. Organize your OpenAPI files in a clear directory structure
2. Use relative paths for external references
3. Ensure all referenced files are valid JSON
4. Keep your references consistent and well-organized
5. Validate your OpenAPI specification before and after inlining

## Limitations

- Only supports JSON format (YAML not supported yet)
- External references must be relative to the input file's directory
- All referenced files must be accessible and valid JSON
- Circular references are not supported

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Troubleshooting

### Common Issues

1. **File Not Found**
   - Ensure the input file path is correct
   - Check if you have read permissions for the file
   - Verify that all referenced files exist

2. **Invalid JSON**
   - Validate your JSON files using a JSON validator
   - Check for missing commas or brackets
   - Ensure all referenced files are valid JSON

3. **Reference Resolution Failed**
   - Verify that reference paths are correct
   - Check if referenced schemas exist in the specified location
   - Ensure external file paths are relative to the input file

### Getting Help

If you encounter any issues:
1. Check the error message carefully
2. Verify your input files
3. Ensure all referenced files are accessible
4. Check file permissions
5. Create an issue in the GitHub repository

## Development

To contribute to the project:

1. Fork the repository
2. Create your feature branch
3. Run tests (when available)
4. Submit a pull request

## Future Enhancements

- YAML support
- Circular reference detection
- Schema validation
- Support for remote references (HTTP/HTTPS)
- Configuration file support
- Custom output formatting options
