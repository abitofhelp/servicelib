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
- **Dependency Injection Component** - Shows the structure of the dependency injection component
- **Health Check Component** - Shows the structure of the health check component
- **Error Handling Component** - Shows the structure of the error handling system, including the relationships between different error types
- **Retry Component** - Shows the structure of the retry package, including its configuration options and relationships with other components

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

- [Package Dependencies](http://www.plantuml.com/plantuml/proxy?src=https://raw.githubusercontent.com/abitofhelp/servicelib/main/docs/diagrams/source/package_dependencies.puml)
- [Improved Package Dependencies](http://www.plantuml.com/plantuml/proxy?src=https://raw.githubusercontent.com/abitofhelp/servicelib/main/docs/diagrams/source/improved_package_dependencies.puml)
- [Layered Architecture](http://www.plantuml.com/plantuml/proxy?src=https://raw.githubusercontent.com/abitofhelp/servicelib/main/docs/diagrams/source/layered_architecture.puml)
- [Authentication Component](http://www.plantuml.com/plantuml/proxy?src=https://raw.githubusercontent.com/abitofhelp/servicelib/main/docs/diagrams/source/auth_component.puml)
- [Dependency Injection Component](http://www.plantuml.com/plantuml/proxy?src=https://raw.githubusercontent.com/abitofhelp/servicelib/main/docs/diagrams/source/di_component.puml)
- [Health Check Component](http://www.plantuml.com/plantuml/proxy?src=https://raw.githubusercontent.com/abitofhelp/servicelib/main/docs/diagrams/source/health_component.puml)
- [Error Handling Component](http://www.plantuml.com/plantuml/proxy?src=https://raw.githubusercontent.com/abitofhelp/servicelib/main/docs/diagrams/source/errors_component.puml)
- [Retry Component](http://www.plantuml.com/plantuml/proxy?src=https://raw.githubusercontent.com/abitofhelp/servicelib/main/docs/diagrams/source/retry_component.puml)

### SVG Files

You can view the SVG files for each diagram directly:

- [Package Dependencies](svg/package_dependencies.svg)
- [Layered Architecture](svg/layered_architecture.svg)
- [Authentication Component](svg/auth_component.svg)
- [Dependency Injection Component](svg/di_component.svg)
- [Health Check Component](svg/health_component.svg)
- [Error Handling Component](svg/errors_component.svg)
- [Retry Component](svg/retry_component.svg)

### Source Files

You can view the source files for each diagram directly:

- [Package Dependencies](source/package_dependencies.puml)
- [Improved Package Dependencies](source/improved_package_dependencies.puml)
- [Layered Architecture](source/layered_architecture.puml)
- [Authentication Component](source/auth_component.puml)
- [Dependency Injection Component](source/di_component.puml)
- [Health Check Component](source/health_component.puml)
- [Error Handling Component](source/errors_component.puml)
- [Retry Component](source/retry_component.puml)

## Relationship to Code

The UML diagrams provide visual representations of the ServiceLib architecture and components. Here's how they relate to the codebase:

- **Package Dependencies**: Shows how the different packages in ServiceLib depend on each other. This helps developers understand the overall structure of the library and how components interact.

- **Layered Architecture**: Illustrates the layered architecture of ServiceLib, showing how components are organized into Core, Infrastructure, Service, and Application layers. This corresponds to the package organization in the codebase.

- **Authentication Component**: Shows the internal structure of the `auth` package, including the relationships between different authentication mechanisms (JWT, OIDC) and how they interact with middleware and other components.

- **Dependency Injection Component**: Shows the internal structure of the `di` package, illustrating how the dependency injection container works and how it manages service dependencies.

- **Health Check Component**: Shows the internal structure of the `health` package, including how health checks are registered, managed, and exposed via HTTP endpoints.

- **Error Handling Component**: Shows the internal structure of the `errors` package, including the hierarchy of error types (BaseError, DomainError, ApplicationError, InfrastructureError) and their relationships. This helps developers understand how to use the error handling system and how different error types relate to each other.

- **Retry Component**: Shows the internal structure of the `retry` package, including the configuration options, function types, and relationships with other packages like errors, logging, and telemetry. This helps developers understand how to use the retry functionality and how it integrates with other components in the system.

## Generating Diagrams

SVG versions of all diagrams are already available in the `svg` directory. However, if you need to regenerate them or create PNG versions, you can do so using PlantUML.

### Prerequisites

1. Install PlantUML: https://plantuml.com/starting
2. Navigate to the repository root directory

### Generating SVG Files

```bash
# Generate a specific diagram
plantuml -tsvg docs/diagrams/source/package_dependencies.puml

# Generate all diagrams
plantuml -tsvg docs/diagrams/source/*.puml
```

### Generating PNG Files

```bash
# Generate a specific diagram
plantuml -tpng docs/diagrams/source/package_dependencies.puml

# Generate all diagrams
plantuml -tpng docs/diagrams/source/*.puml
```

The generated files will be placed in the same directory as the source files. You may want to move them to a more appropriate location:

```bash
# Create directories if they don't exist
mkdir -p docs/diagrams/svg docs/diagrams/png

# Move SVG files
mv docs/diagrams/source/*.svg docs/diagrams/svg/

# Move PNG files
mv docs/diagrams/source/*.png docs/diagrams/png/
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
