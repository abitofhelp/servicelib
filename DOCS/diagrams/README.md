# UML Diagrams for ServiceLib

This directory contains UML diagrams that provide a visual representation of the ServiceLib architecture and components.

## Directory Structure

- `source/`: Contains the PlantUML (.puml) source files for the diagrams
- `svg/`: Contains the generated SVG files

## Available Diagrams

1. **Architecture Overview** - High-level overview of the ServiceLib architecture based on Clean Architecture, DDD, and Hexagonal Architecture
2. **Package Diagram** - Diagram showing the packages in ServiceLib and their relationships
3. **HTTP Request Sequence** - Sequence diagram illustrating the HTTP request processing flow
4. **Errors Package Class Diagram** - Class diagram for the errors package

## Updating Diagrams

To update a diagram:

1. Edit the corresponding .puml file in the `source/` directory
2. Run the `tools/generate_svg.sh` script to generate the SVG files:
   ```bash
   ./tools/generate_svg.sh
   ```

## Adding New Diagrams

To add a new diagram:

1. Create a new .puml file in the `source/` directory
2. Run the `tools/generate_svg.sh` script to generate the SVG files
3. Update this README.md file to include the new diagram
4. Consider adding a reference to the new diagram in the main README.md file

## PlantUML Resources

- [PlantUML Official Website](https://plantuml.com/)
- [PlantUML Language Reference Guide](https://plantuml.com/guide)
- [PlantUML Class Diagram Syntax](https://plantuml.com/class-diagram)
- [PlantUML Sequence Diagram Syntax](https://plantuml.com/sequence-diagram)
- [PlantUML Component Diagram Syntax](https://plantuml.com/component-diagram)