# ServiceLib UML Diagrams

This directory contains UML diagrams for the ServiceLib library. These diagrams provide visual representations of the library's architecture, component relationships, and internal structures.

## Available Diagrams

The following diagrams can be generated from the source files in the `source` directory. See the [Generating Diagrams](#generating-diagrams) section below for instructions.

### Architecture Diagrams

- **Package Dependencies** - Shows the dependencies between the main packages in ServiceLib
- **Improved Package Dependencies** - Shows a proposed improved package structure with interface packages to reduce coupling
- **Layered Architecture** - Illustrates the layered architecture of ServiceLib

### Component Diagrams

- **Authentication Component** - Shows the structure of the authentication component
- **Database Component** - Shows the structure of the database package, including its interfaces and functions
- **Dependency Injection Component** - Shows the structure of the dependency injection component
- **Error Handling Component** - Shows the structure of the error handling system, including the relationships between different error types
- **GraphQL Component** - Shows the structure of the GraphQL package, including server configuration and directives
- **Health Check Component** - Shows the structure of the health check component
- **Logging Component** - Shows the structure of the logging package, including the logger interface and context logger
- **Middleware Component** - Shows the structure of the middleware package, including middleware functions and utilities
- **Retry Component** - Shows the structure of the retry package, including its configuration options and relationships with other components
- **Telemetry Component** - Shows the structure of the telemetry package, including metrics and tracing components

## Diagram Formats

The diagrams are available in the following formats:

- SVG - Scalable vector graphics for high-quality display at any resolution (available in the `svg` directory)
- PlantUML (.puml) - Source files for generating and modifying the diagrams (available in the `source` directory)

The SVG files have been pre-generated for convenience. If you need to regenerate them or create PNG versions, see the [Generating Diagrams](#generating-diagrams) section below for instructions.

## Viewing Diagrams Without Generating Them

If you don't want to install PlantUML locally, you can view the diagrams using online tools:

1. **PlantUML Online Server**: You can use the [PlantUML Online Server](http://www.plantuml.com/plantuml/uml/) to view and edit PlantUML diagrams. Simply copy the content of a .puml file and paste it into the editor.

2. **GitHub Integration**: If you're viewing the repository on GitHub, you can install a browser extension like [PlantUML Viewer](https://chrome.google.com/webstore/detail/plantuml-viewer/legbfeljfbjgfifnkmpoajgpgejojooj) for Chrome or [PlantUML Visualizer](https://addons.mozilla.org/en-US/firefox/addon/plantuml-visualizer/) for Firefox. These extensions will render PlantUML diagrams directly in GitHub.

3. **VS Code Extension**: If you're using Visual Studio Code, you can install the [PlantUML extension](https://marketplace.visualstudio.com/items?itemName=jebbs.plantuml) to preview diagrams directly in the editor.

### Quick View Links

For convenience, here are links to view each diagram directly in the PlantUML Online Server:

- [Package Dependencies](http://www.plantuml.com/plantuml/proxy?src=https://raw.githubusercontent.com/abitofhelp/servicelib/main/DOCS/diagrams/source/package_dependencies.puml)
- [Improved Package Dependencies](http://www.plantuml.com/plantuml/proxy?src=https://raw.githubusercontent.com/abitofhelp/servicelib/main/DOCS/diagrams/source/improved_package_dependencies.puml)
- [Layered Architecture](http://www.plantuml.com/plantuml/proxy?src=https://raw.githubusercontent.com/abitofhelp/servicelib/main/DOCS/diagrams/source/layered_architecture.puml)
- [Authentication Component](http://www.plantuml.com/plantuml/proxy?src=https://raw.githubusercontent.com/abitofhelp/servicelib/main/DOCS/diagrams/source/auth_component.puml)
- [Database Component](http://www.plantuml.com/plantuml/proxy?src=https://raw.githubusercontent.com/abitofhelp/servicelib/main/DOCS/diagrams/source/db_component.puml)
- [Dependency Injection Component](http://www.plantuml.com/plantuml/proxy?src=https://raw.githubusercontent.com/abitofhelp/servicelib/main/DOCS/diagrams/source/di_component.puml)
- [Error Handling Component](http://www.plantuml.com/plantuml/proxy?src=https://raw.githubusercontent.com/abitofhelp/servicelib/main/DOCS/diagrams/source/errors_component.puml)
- [GraphQL Component](http://www.plantuml.com/plantuml/proxy?src=https://raw.githubusercontent.com/abitofhelp/servicelib/main/DOCS/diagrams/source/graphql_component.puml)
- [Health Check Component](http://www.plantuml.com/plantuml/proxy?src=https://raw.githubusercontent.com/abitofhelp/servicelib/main/DOCS/diagrams/source/health_component.puml)
- [Logging Component](http://www.plantuml.com/plantuml/proxy?src=https://raw.githubusercontent.com/abitofhelp/servicelib/main/DOCS/diagrams/source/logging_component.puml)
- [Middleware Component](http://www.plantuml.com/plantuml/proxy?src=https://raw.githubusercontent.com/abitofhelp/servicelib/main/DOCS/diagrams/source/middleware_component.puml)
- [Retry Component](http://www.plantuml.com/plantuml/proxy?src=https://raw.githubusercontent.com/abitofhelp/servicelib/main/DOCS/diagrams/source/retry_component.puml)
- [Telemetry Component](http://www.plantuml.com/plantuml/proxy?src=https://raw.githubusercontent.com/abitofhelp/servicelib/main/DOCS/diagrams/source/telemetry_component.puml)

### SVG Files

You can view the SVG files for each diagram directly:

- [Package Dependencies](svg/package_dependencies.svg)
- [Improved Package Dependencies](svg/improved_package_dependencies.svg)
- [Layered Architecture](svg/layered_architecture.svg)
- [Authentication Component](svg/auth_component.svg)
- [Database Component](svg/db_component.svg)
- [Dependency Injection Component](svg/di_component.svg)
- [Error Handling Component](svg/errors_component.svg)
- [GraphQL Component](svg/graphql_component.svg)
- [Health Check Component](svg/health_component.svg)
- [Logging Component](svg/logging_component.svg)
- [Middleware Component](svg/middleware_component.svg)
- [Retry Component](svg/retry_component.svg)
- [Telemetry Component](svg/telemetry_component.svg)

### Source Files

You can view the source files for each diagram directly:

- [Package Dependencies](source/package_dependencies.puml)
- [Improved Package Dependencies](source/improved_package_dependencies.puml)
- [Layered Architecture](source/layered_architecture.puml)
- [Authentication Component](source/auth_component.puml)
- [Database Component](source/db_component.puml)
- [Dependency Injection Component](source/di_component.puml)
- [Error Handling Component](source/errors_component.puml)
- [GraphQL Component](source/graphql_component.puml)
- [Health Check Component](source/health_component.puml)
- [Logging Component](source/logging_component.puml)
- [Middleware Component](source/middleware_component.puml)
- [Retry Component](source/retry_component.puml)
- [Telemetry Component](source/telemetry_component.puml)

## Relationship to Code

The UML diagrams provide visual representations of the ServiceLib architecture and components. Here's how they relate to the codebase:

- **Package Dependencies**: Shows how the different packages in ServiceLib depend on each other. This helps developers understand the overall structure of the library and how components interact.

- **Layered Architecture**: Illustrates the layered architecture of ServiceLib, showing how components are organized into Core, Infrastructure, Service, and Application layers. This corresponds to the package organization in the codebase.

- **Authentication Component**: Shows the internal structure of the `auth` package, including the relationships between different authentication mechanisms (JWT, OIDC) and how they interact with middleware and other components.

- **Database Component**: Shows the internal structure of the `db` package, including the interfaces and functions for working with different types of databases (PostgreSQL, MongoDB, SQLite). This helps developers understand how to use the database functionality and how it integrates with other components.

- **Dependency Injection Component**: Shows the internal structure of the `di` package, illustrating how the dependency injection container works and how it manages service dependencies.

- **Error Handling Component**: Shows the internal structure of the `errors` package, including the hierarchy of error types (BaseError, DomainError, ApplicationError, InfrastructureError) and their relationships. This helps developers understand how to use the error handling system and how different error types relate to each other.

- **GraphQL Component**: Shows the internal structure of the `graphql` package, including server configuration, directives, and error handling. This helps developers understand how to use the GraphQL functionality and how it integrates with other components.

- **Health Check Component**: Shows the internal structure of the `health` package, including how health checks are registered, managed, and exposed via HTTP endpoints.

- **Logging Component**: Shows the internal structure of the `logging` package, including the logger interface, context logger, and their relationships with other components. This helps developers understand how to use the logging functionality and how it integrates with other components.

- **Middleware Component**: Shows the internal structure of the `middleware` package, including middleware functions, context utilities, and response writers. This helps developers understand how to use the middleware functionality and how it integrates with other components.

- **Retry Component**: Shows the internal structure of the `retry` package, including the configuration options, function types, and relationships with other packages like errors, logging, and telemetry. This helps developers understand how to use the retry functionality and how it integrates with other components in the system.

- **Telemetry Component**: Shows the internal structure of the `telemetry` package, including metrics, tracing, and their relationships with other components. This helps developers understand how to use the telemetry functionality and how it integrates with other components.

## Generating Diagrams

SVG versions of all diagrams are already available in the `svg` directory. However, if you need to regenerate them or create PNG versions, you can do so using PlantUML.

### Prerequisites

1. Install PlantUML: https://plantuml.com/starting
2. Navigate to the repository root directory

### Generating SVG Files

```bash
# Generate a specific diagram
plantuml -tsvg DOCS/diagrams/source/package_dependencies.puml

# Generate all diagrams
plantuml -tsvg DOCS/diagrams/source/*.puml
```

### Generating PNG Files

```bash
# Generate a specific diagram
plantuml -tpng DOCS/diagrams/source/package_dependencies.puml

# Generate all diagrams
plantuml -tpng DOCS/diagrams/source/*.puml
```

The generated files will be placed in the same directory as the source files. You may want to move them to a more appropriate location:

```bash
# Create directories if they don't exist
mkdir -p DOCS/diagrams/svg DOCS/diagrams/png

# Move SVG files
mv DOCS/diagrams/source/*.svg DOCS/diagrams/svg/

# Move PNG files
mv DOCS/diagrams/source/*.png DOCS/diagrams/png/
```

## Contributing

If you'd like to contribute new diagrams or improve existing ones, please follow these guidelines:

1. Use PlantUML for creating diagrams
2. Follow the existing style and naming conventions
3. Place source files in the `source` directory with a descriptive name (e.g., `component_name.puml`)
4. Generate SVG versions of the diagrams following the instructions in the [Generating Diagrams](#generating-diagrams) section
5. Place the generated SVG files in the `svg` directory with a consistent naming convention (e.g., `component_name.svg`)
6. Update this README.md file with information about new diagrams, including links to both the source and SVG files
7. See the project's [CONTRIBUTING.md](../../../CONTRIBUTING.md) file for general contribution guidelines
