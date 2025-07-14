# OpenFHE Go - GitHub Workflow Setup

This document describes the GitHub workflow setup for building and testing the OpenFHE Go project.

## Workflow Overview

The GitHub workflow (`.github/workflows/build.yml`) is configured to:

1. **Use OpenFHE Buildbox Docker Image**: Uses `ghcr.io/matiasinsaurralde/openfhe-buildbox:latest` as the base container
2. **Build the Project**: Compiles the Go code with OpenFHE bindings using CGO
3. **Run Tests**: Executes all tests with proper OpenFHE environment
4. **Lint Code**: Runs golangci-lint to ensure code quality

## Key Features

### Docker Container
- **Image**: `ghcr.io/matiasinsaurralde/openfhe-buildbox:latest`
- **Platform**: Linux AMD64
- **Includes**: Pre-installed OpenFHE libraries and headers

### Build Configuration
- **Build Tags**: Uses `openfhe` build tag for conditional compilation
- **CGO**: Enabled with proper environment variables
- **Go Version**: 1.24.3
- **Libraries**: Links against OpenFHE core and PKE libraries

### Environment Variables
The workflow sets up the following environment variables:
```bash
CGO_ENABLED=1
LD_LIBRARY_PATH=/usr/local/lib
C_INCLUDE_PATH=/usr/local/include
CPLUS_INCLUDE_PATH=/usr/local/include
LIBRARY_PATH=/usr/local/lib
```

### Code Quality
- **golangci-lint**: Comprehensive linting with custom configuration
- **gofmt**: Automatic code formatting
- **Build Tags**: Proper conditional compilation for OpenFHE bindings

## Project Structure

```
.
├── .github/workflows/build.yml      # GitHub workflow
├── .golangci.yml                    # golangci-lint configuration
├── openfhe.go                       # Go bindings for OpenFHE
├── openfhe_test.go                  # Test suite
├── openfhe_c.h                      # C header for bindings
├── openfhe_c.cpp                    # C++ implementation
├── go.mod                           # Go module definition
├── test_workflow.sh                 # Local workflow test script
└── README.md                        # Project documentation
```

## Build Tags

The project uses build tags to conditionally compile OpenFHE-dependent code:

- **`openfhe`**: Required build tag for compilation
- **`cgo`**: Required for C interoperability

### Example Build Commands

```bash
# Build with OpenFHE support
go build -tags="openfhe" -v ./...

# Run tests with OpenFHE support
go test -tags="openfhe" -v ./...

# Run linter
golangci-lint run --timeout=10m
```

## Local Testing

Use the provided test script to simulate the workflow locally:

```bash
./test_workflow.sh
```

**Note**: This script will fail at the compilation step without OpenFHE libraries installed, which is expected behavior.

## OpenFHE Integration

The project provides Go bindings for OpenFHE with the following features:

### Supported Operations
- **Context Management**: BGV-RNS crypto context creation
- **Key Generation**: Public/private key pair generation
- **Encryption/Decryption**: Encrypt/decrypt uint64 values
- **Homomorphic Operations**: Addition of encrypted values
- **Serialization**: Binary serialization of keys and ciphertexts

### Key Types
- `Context`: BGV-RNS crypto context wrapper
- `PublicKey`: Public key wrapper
- `SecretKey`: Secret key wrapper  
- `Ciphertext`: Encrypted data wrapper

### Example Usage

```go
// Create context
ctx := NewBGVRNS(2, 65537)
defer ctx.Free()

// Generate keys
pk, sk, err := ctx.KeyGenPtr()
if err != nil {
    log.Fatal(err)
}
defer pk.Free()
defer sk.Free()

// Encrypt value
ct, err := ctx.EncryptU64ToPtr(pk, 42)
if err != nil {
    log.Fatal(err)
}
defer ct.Free()

// Decrypt value
plain, err := ctx.DecryptU64FromPtr(sk, ct)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Decrypted value: %d\n", plain)
```

## Continuous Integration

The workflow runs automatically on:
- **Push** to `main` or `master` branches
- **Pull requests** to `main` or `master` branches

### Workflow Steps
1. **Checkout**: Get the latest code
2. **Setup Go**: Install Go 1.24.3
3. **Verify Environment**: Check OpenFHE installation
4. **Install golangci-lint**: Install the latest linter
5. **Download Dependencies**: Get Go module dependencies
6. **Build**: Compile with OpenFHE tags
7. **Test**: Run test suite
8. **Lint**: Run code quality checks

## Troubleshooting

### Common Issues

1. **OpenFHE Headers Not Found**
   - Ensure the buildbox image contains OpenFHE in `/usr/local/include/openfhe/`
   - Check that environment variables are set correctly

2. **Library Linking Errors**
   - Verify OpenFHE libraries are in `/usr/local/lib/`
   - Check `LD_LIBRARY_PATH` includes `/usr/local/lib`

3. **Build Tag Issues**
   - Always use `-tags="openfhe"` when building
   - Code won't compile without proper build tags

4. **CGO Compilation Errors**
   - Ensure `CGO_ENABLED=1`
   - Check that C++ compiler is available
   - Verify OpenFHE headers are accessible

### Monitoring Workflow

To monitor the workflow:

1. Go to the repository's Actions tab
2. Check the latest workflow run
3. View logs for each step
4. Look for compilation or test failures

### Iteration Process

When the workflow fails:

1. Check the workflow logs for specific errors
2. Fix the identified issues
3. Commit and push changes
4. Monitor the new workflow run
5. Repeat until all tests pass

## Development Notes

- The project requires OpenFHE libraries for compilation
- Build tags ensure code only compiles when OpenFHE is available
- Local development requires OpenFHE installation or Docker container
- Tests verify encryption, decryption, and homomorphic operations
- Serialization support allows persistence of cryptographic objects

## Future Improvements

- Add support for additional OpenFHE schemes (CKKS, BFV)
- Implement more homomorphic operations (multiplication, rotation)
- Add benchmarking tests
- Support for different parameter sets
- Documentation generation
- Release automation