# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Build and Run
```
go build
go run main.go
```

### Tests
```
go test ./...                        # Run all tests
go test ./service/...                # Test specific package
go test -run TestForexRateServiceTestSuite/TestGetRateByCurrencyPair # Single test
go test ./dao/integrationtest/...    # Run integration tests
```

## Code Style Guidelines

1. **Package Structure**: Maintain clear separation between model, service, controller, dao, apiclient packages
2. **Naming**: Use CamelCase for functions/vars/types; interface names are verbs (e.g., ForexApiClient)
3. **Interfaces**: Define interfaces in separate `interface.go` files within packages
4. **Error Handling**: Explicit error checking with if statements; propagate errors up the call stack
5. **Testing**: Use testify suite for test organization; Given/When/Then pattern in comments
6. **Dependency Injection**: Use dig (go.uber.org/dig); constructor pattern with NewXXX functions
7. **Imports**: Group by standard lib, third-party libs, project imports with blank lines between
8. **Configuration**: Use YAML with environment-specific versions (application-local.yaml, etc.)