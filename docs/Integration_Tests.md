# Integration Tests

This document describes how to run the integration tests for the ServiceLib library.

## Overview

ServiceLib includes integration tests that verify the interaction between different components and external systems. These tests are separate from unit tests and require additional setup to run.

Integration tests are marked with the `integration` build tag to prevent them from running during normal test execution.

## Test Structure

Integration tests are organized by package:

- `auth/integration/` - Tests for authentication and authorization components
  - [`jwt_integration_test.go`](../auth/integration/jwt_integration_test.go) - Tests JWT token validation flow
  - [`oidc_integration_test.go`](../auth/integration/oidc_integration_test.go) - Tests OIDC integration with a mock provider

- `db/integration/` - Tests for database components
  - [`db_integration_test.go`](../db/integration/db_integration_test.go) - Tests direct database connections and operations

- `telemetry/integration/` - Tests for telemetry components
  - [`metrics_integration_test.go`](../telemetry/integration/metrics_integration_test.go) - Tests metrics recording and Prometheus endpoint
  - [`http_integration_test.go`](../telemetry/integration/http_integration_test.go) - Tests HTTP instrumentation

## Running Integration Tests

> **Important**: In order to ensure that developers can build and work with this package within different IDEs and environments, please use the Makefile to build, test, etc.

Using the Makefile (recommended):

```bash
# Run all integration tests
make test-integration

# Run integration tests with coverage
make coverage-integration
```

Alternatively, you can use Go commands directly:

To run all integration tests:

```bash
go test -tags=integration ./...
```

To run integration tests for a specific package:

```bash
go test -tags=integration ./auth/integration
go test -tags=integration ./telemetry/integration
```

To run a specific integration test:

```bash
go test -tags=integration -run TestJWTAuthenticationFlow ./auth/integration
go test -tags=integration -run TestPrometheusMetricsEndpoint ./telemetry/integration
```

## Test Environment Requirements

### Auth Integration Tests

The JWT integration tests run without any external dependencies.

The OIDC integration tests are skipped by default because they require an external OIDC provider. To run these tests:

1. Set up an OIDC provider (e.g., Keycloak, Auth0)
2. Update the test configuration with your provider details
3. Remove the `t.Skip()` line from the tests

### DB Integration Tests

The database integration tests are skipped by default because they require real database connections:
- PostgreSQL health check and transaction tests
- MongoDB health check tests
- SQLite health check and transaction tests

To run these tests:

1. Set up the required database (PostgreSQL, MongoDB, or SQLite)
2. Update the connection strings in the tests
3. Remove the `t.Skip()` line from the tests you want to run

### Telemetry Integration Tests

Some telemetry integration tests run without external dependencies:
- Prometheus metrics endpoint tests
- HTTP instrumentation tests
- Span operation tests

The metrics recording tests are skipped by default because they require an OpenTelemetry collector. To run these tests:

1. Set up an OpenTelemetry collector
2. Update the test configuration with your collector details
3. Remove the `t.Skip()` line from the tests

## Mock Implementations

The integration tests use mock implementations where appropriate:

- `auth/integration/oidc_integration_test.go` includes a mock OIDC provider
- `telemetry/integration/metrics_integration_test.go` includes fallback behavior when no OpenTelemetry collector is available

## Best Practices

When writing new integration tests:

1. Use the `integration` build tag
2. Create test servers and clients for testing HTTP components
3. Use `testify/assert` and `testify/require` for assertions
4. Implement proper test setup and teardown
5. Skip tests that require external dependencies by default
6. Document the requirements for running skipped tests
7. Ensure tests are isolated and can run independently
8. Add comprehensive assertions to verify behavior

## Common Test Scenarios

Here are some common test scenarios and how to run them:

### Testing JWT Authentication

The JWT authentication tests verify that tokens can be generated, validated, and used for authorization:

```bash
go test -tags=integration -run TestJWTAuthenticationFlow ./auth/integration
```

See [jwt_integration_test.go](../auth/integration/jwt_integration_test.go) for implementation details.

### Testing Metrics Collection

The metrics collection tests verify that metrics are properly recorded and exposed via Prometheus:

```bash
go test -tags=integration -run TestPrometheusMetricsEndpoint ./telemetry/integration
```

See [metrics_integration_test.go](../telemetry/integration/metrics_integration_test.go) for implementation details.

### Testing HTTP Instrumentation

The HTTP instrumentation tests verify that HTTP requests are properly instrumented with tracing and metrics:

```bash
go test -tags=integration -run TestHTTPInstrumentation ./telemetry/integration
```

See [http_integration_test.go](../telemetry/integration/http_integration_test.go) for implementation details.

### Testing Database Connections

The database connection tests verify that the database health check and transaction functions work correctly:

```bash
go test -tags=integration -run TestCheckPostgresHealthDirect ./db/integration
go test -tags=integration -run TestExecutePostgresTransactionDirect ./db/integration
```

See [db_integration_test.go](../db/integration/db_integration_test.go) for implementation details.

## Troubleshooting

### Common Issues

- **Test Timeouts**: Integration tests may take longer to run than unit tests. Use the `-timeout` flag to increase the test timeout if needed.

  ```bash
  go test -tags=integration -timeout 5m ./...
  ```

- **External Dependencies**: If tests fail due to missing external dependencies, check that you've set up the required services and updated the test configuration.

- **Port Conflicts**: Integration tests may start servers on specific ports. If you encounter port conflicts, modify the test to use a different port or use `httptest.NewServer()` which assigns a random port.

- **Authentication Failures**: OIDC integration tests may fail due to authentication issues. Check that your OIDC provider is properly configured and that the client credentials are correct.

- **Network Connectivity**: Tests that communicate with external services may fail due to network issues. Check your network connectivity and firewall settings.

- **Environment Variables**: Some tests may require specific environment variables to be set. Check the test file for any required environment variables.

### Common Test Failures

- **JWT Token Validation Failures**: These often occur due to incorrect secret keys or token expiration. Check that the JWT configuration in the test matches the configuration in the code being tested.

- **OIDC Provider Unavailable**: If the OIDC integration tests fail with connection errors, check that the OIDC provider is running and accessible.

- **Metrics Recording Failures**: These may occur if the OpenTelemetry collector is not properly configured or not running. Check the collector configuration and status.

- **HTTP Instrumentation Failures**: These may occur if the HTTP server or client is not properly configured. Check the HTTP configuration in the test.

- **Database Connection Failures**: These may occur if the database is not running or if the connection string is incorrect. Check that the database is running and accessible, and that the connection string is correct.

### Debugging Test Failures

When a test fails, follow these steps to debug the issue:

1. **Run the test with verbose output**: Use the `-v` flag to see detailed output from the test.

   ```bash
   go test -tags=integration -v -run TestSpecificTest ./path/to/package
   ```

2. **Check the test logs**: Look for error messages in the test output that might indicate the cause of the failure.

3. **Inspect the test code**: Look at the test file to understand what the test is doing and what might be causing it to fail.

4. **Use the Go debugger**: If necessary, use the Go debugger to step through the test code and identify the issue.

   ```bash
   dlv test --build-flags="-tags=integration" ./path/to/package -- -test.run TestSpecificTest
   ```

5. **Check external dependencies**: If the test interacts with external services, check that those services are running and accessible.

### Common Test Failures

- **JWT Token Validation Failures**: These often occur due to incorrect secret keys or token expiration. Check that the JWT configuration in the test matches the configuration in the code being tested.

- **OIDC Provider Unavailable**: If the OIDC integration tests fail with connection errors, check that the OIDC provider is running and accessible.

- **Metrics Recording Failures**: These may occur if the OpenTelemetry collector is not properly configured or not running. Check the collector configuration and status.

- **HTTP Instrumentation Failures**: These may occur if the HTTP server or client is not properly configured. Check the HTTP configuration in the test.

- **Database Connection Failures**: These may occur if the database is not running or if the connection string is incorrect. Check that the database is running and accessible, and that the connection string is correct.
