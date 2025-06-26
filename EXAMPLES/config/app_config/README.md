# config app_config Example

## Overview

This example demonstrates how to implement and use the AppConfigProvider interface from the ServiceLib config package. It shows how to create a custom configuration struct that provides both standard application configuration and additional application-specific settings.

## Features

- **Custom App Configuration**: Create a custom configuration struct that implements the AppConfigProvider interface
- **Extended Functionality**: Add application-specific methods beyond the standard interface
- **Feature Flags**: Implement and use feature flags to control application behavior
- **Configuration Access**: Access configuration through both the standard interface and custom methods

## Running the Example

To run this example, navigate to this directory and execute:

```bash
go run main.go
```

## Code Walkthrough

### Implementing AppConfigProvider Interface

The example defines a custom configuration struct that implements the required AppConfigProvider interface:

```
// AppSettings is a custom configuration struct that implements the AppConfigProvider interface
type AppSettings struct {
    AppVersion     string
    AppName        string
    AppEnvironment string
    LogLevel       string
    Features       map[string]bool
}

// GetAppVersion implements the AppConfigProvider interface
func (s *AppSettings) GetAppVersion() string {
    return s.AppVersion
}

// GetAppName implements the AppConfigProvider interface
func (s *AppSettings) GetAppName() string {
    return s.AppName
}

// GetAppEnvironment implements the AppConfigProvider interface
func (s *AppSettings) GetAppEnvironment() string {
    return s.AppEnvironment
}
```

### Adding Custom Methods

The example adds application-specific methods beyond the standard interface:

```
// Additional methods specific to AppSettings
func (s *AppSettings) GetLogLevel() string {
    return s.LogLevel
}

func (s *AppSettings) IsFeatureEnabled(featureName string) bool {
    if enabled, ok := s.Features[featureName]; ok {
        return enabled
    }
    return false
}
```

### Accessing Configuration Values

The example demonstrates how to access configuration through both the standard interface and custom methods:

```
// Create a config adapter
adapter := config.NewGenericConfigAdapter(settings)

// Get the app configuration through the adapter
appConfig := adapter.GetApp()

// Use the standard AppConfig interface methods
fmt.Println("=== Application Configuration ===")
fmt.Println("Version:", appConfig.GetVersion())
fmt.Println("Name:", appConfig.GetName())
fmt.Println("Environment:", appConfig.GetEnvironment())

// Use the original settings object for additional functionality
fmt.Println("\n=== Additional Settings ===")
fmt.Println("Log Level:", settings.GetLogLevel())
fmt.Println("Dark Mode Enabled:", settings.IsFeatureEnabled("darkMode"))
```

## Expected Output

```
=== Application Configuration ===
Version: 2.1.0
Name: FeatureApp
Environment: staging

=== Additional Settings ===
Log Level: debug
Dark Mode Enabled: true
Beta Features Enabled: true
Notifications Enabled: false
```

## Related Examples

- [basic_usage](../basic_usage/README.md) - Demonstrates the basic usage of the configuration package with a custom configuration structure
- [custom_adapter](../custom_adapter/README.md) - Shows how to create a custom configuration adapter that reads from environment variables
- [database_config](../database_config/README.md) - Illustrates how to implement database-specific configuration

## Related Components

- [config Package](../../../config/README.md) - The config package documentation.
- [Errors Package](../../../errors/README.md) - The errors package used for error handling.
- [Context Package](../../../context/README.md) - The context package used for context handling.
- [Logging Package](../../../logging/README.md) - The logging package used for logging.

## License

This project is licensed under the MIT License - see the [LICENSE](../../../LICENSE) file for details.
