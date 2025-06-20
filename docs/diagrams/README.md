# ServiceLib UML Diagrams

This directory contains UML diagrams for the ServiceLib library. These diagrams provide visual representations of the library's architecture, component relationships, and internal structures.

## Available Diagrams

The following diagrams can be generated from the source files in the `source` directory. See the [Generating Diagrams](#generating-diagrams) section below for instructions.

### Architecture Diagrams

- **Package Dependencies** - Shows the dependencies between the main packages in ServiceLib
- **Layered Architecture** - Illustrates the layered architecture of ServiceLib

### Component Diagrams

- **Authentication Component** - Shows the structure of the authentication component
- **Dependency Injection Component** - Shows the structure of the dependency injection component
- **Health Check Component** - Shows the structure of the health check component

## Diagram Formats

The diagrams can be generated in the following formats:

- PNG - For viewing in browsers and documentation
- SVG - For scalable vector graphics

The source files are provided in:

- PlantUML (.puml) - Source files for generating and modifying the diagrams

Note: The PNG and SVG files need to be generated from the PlantUML source files. See the [Generating Diagrams](#generating-diagrams) section below for instructions.

## Generating Diagrams

The diagrams can be generated using PlantUML. Follow these steps:

1. Install PlantUML: https://plantuml.com/starting
2. Navigate to the repository root directory
3. Generate PNG files:

```bash
# Generate a specific diagram
plantuml -tpng docs/diagrams/source/package_dependencies.puml

# Generate all diagrams
plantuml -tpng docs/diagrams/source/*.puml
```

4. Generate SVG files:

```bash
# Generate a specific diagram
plantuml -tsvg docs/diagrams/source/package_dependencies.puml

# Generate all diagrams
plantuml -tsvg docs/diagrams/source/*.puml
```

The generated files will be placed in the same directory as the source files. You may want to move them to a more appropriate location, such as the `docs/diagrams` directory.

## Contributing

If you'd like to contribute new diagrams or improve existing ones, please follow these guidelines:

1. Use PlantUML for creating diagrams
2. Follow the existing style and naming conventions
3. Place source files in the `source` directory with a descriptive name (e.g., `component_name.puml`)
4. Generate PNG and SVG versions of the diagrams following the instructions in the [Generating Diagrams](#generating-diagrams) section
5. Place the generated files in the appropriate location
6. Update this README.md file with information about new diagrams
7. See the project's [CONTRIBUTING.md](../../../CONTRIBUTING.md) file for general contribution guidelines
