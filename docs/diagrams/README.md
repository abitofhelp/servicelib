# ServiceLib UML Diagrams

This directory contains UML diagrams for the ServiceLib library. These diagrams provide visual representations of the library's architecture, component relationships, and internal structures.

## Available Diagrams

### Architecture Diagrams

- [Package Dependencies](package_dependencies.png) - Shows the dependencies between the main packages in ServiceLib
- [Layered Architecture](layered_architecture.png) - Illustrates the layered architecture of ServiceLib

### Component Diagrams

- [Authentication Component](auth_component.png) - Shows the structure of the authentication component
- [Dependency Injection Component](di_component.png) - Shows the structure of the dependency injection component
- [Health Check Component](health_component.png) - Shows the structure of the health check component

## Diagram Formats

All diagrams are available in the following formats:

- PNG - For viewing in browsers and documentation
- SVG - For scalable vector graphics
- PlantUML - Source files for generating and modifying the diagrams

## Generating Diagrams

The diagrams can be regenerated using PlantUML. To generate a diagram:

1. Install PlantUML: https://plantuml.com/starting
2. Run the following command:

```bash
plantuml -tpng diagrams/source/diagram_name.puml
```

## Contributing

If you'd like to contribute new diagrams or improve existing ones, please follow these guidelines:

1. Use PlantUML for creating diagrams
2. Follow the existing style and naming conventions
3. Include source files in the `source` directory
4. Generate PNG and SVG versions of the diagrams
5. Update this README.md file with information about new diagrams