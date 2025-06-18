# Integration Tests

This document describes how to run the integration tests for the ServiceLib library.

## Overview

ServiceLib includes integration tests that verify the interaction between different components and external systems. These tests are separate from unit tests and require additional setup to run.

Integration tests are marked with the `integration` build tag to prevent them from running during normal test execution.

## Test Structure

Integration tests are organized by package:

- `auth/integration/` - Tests for authentication and authorization components
  - `jwt_integration_test.go` - Tests JWT token validation flow
  - `oidc_integration_test.go` - Tests OIDC integration with a mock provider

- `telemetry/integration/` - Tests for telemetry components
  - `metrics_integration_test.go` - Tests metrics recording and Prometheus endpoint
  - `http_integration_test.go` - Tests HTTP instrumentation

## Running Integration Tests

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

## Troubleshooting

### Common Issues

- **Test Timeouts**: Integration tests may take longer to run than unit tests. Use the `-timeout` flag to increase the test timeout if needed.

  ```bash
  go test -tags=integration -timeout 5m ./...
  ```

- **External Dependencies**: If tests fail due to missing external dependencies, check that you've set up the required services and updated the test configuration.

- **Port Conflicts**: Integration tests may start servers on specific ports. If you encounter port conflicts, modify the test to use a different port or use `httptest.NewServer()` which assigns a random port.