This project follows the following coding guidelines:
* You will explain decisions and design choices before implementing them.
* The architecture combines the best practices from Domain-Driven Design (DDD), Clean Architecture, and Hexagonal Architecture, creating a hybrid model that leverages the strengths of each.  All code must adhere to these designs.
* The codebase is written in Go, following best practices and conventions.
* Any function requiring a context parameter must have a context.Context parameter as the first argument and use ContextLogger to log messages.
* You must use Testify for assertions in tests.
* You must use GoMock for mocking in tests.
* You must not implement testing for the /examples and /tools folders.
* Implement unit tests for each package and when a new Go file is added.  The target is a minimum of 80% coverage of statements.  
* Unit tests should be fast, must exist in the same folder as the file being tested, and isolated from each other.
* Implement integration tests for each package as needed to adhere to best practices, with a minimum of 80% coverage of statements.
* Always use exhaustive error detection and handling, including retries and timeouts, where appropriate.
* Always update .svg files when updating .puml files.
* Always update the relevant.puml files whenever relevant content changes.
* Always update the documentation whenever there are relevant architecture changes or code changes.  This includes /docs, and any relevant .md files throughout the packages in the codebase.
* Do not embed examples in the documentation.  Rather, link to them in /examples/*.
* All new or existing README.md files must adhere to the *_README_TEMPLATE.md template related to its use. 
* Every package must have a README.md file, even if it is empty.
* All quickstart examples must be embedded in the relevant README.md file.  Verify that there are no issues with embedding golang in the readme.md file.
* Ensure that all readme.md files that use a template do not have any template filler data in them after generation or updating the readme.md.  Remove the offending data from the readme.md and replace it with actual data.  If no actual data exists, use "No Information Available".
* Ensure that the hierarchy of readme.md files from the root is consistent and targets the correct readme.md files in the packages at the same level as /readme.md.
* Ensure that all readme.md files that have links to examples are consistent and targets the correct example's readme.md file.
* Standardize naming conventions between .puml files and .svg files to use underscores and lowercase.
* All dependencies must point inward, never outward.
* All date and time strings must be in RFC3339 format.

  The project must follow these godoc guidelines, which are based on best practices and conventions:
1. Document Exported Identifiers:
   Rule: Every exported identifier (functions, types, variables, constants) should have a doc comment.
   Reason: godoc generates documentation from these comments, and well-documented exported elements are crucial for understanding and using your package.
2. Follow the Comment Format:
   Rule: Place doc comments directly above the declaration they describe.
   Reason: godoc uses this proximity to associate comments with the correct elements.
3. Start the Comment with the Identifier's Name:
   Rule: The first sentence of a doc comment should begin with the name of the documented element.
   Example: // Reader serves content from a ZIP archive.
   Reason: This convention makes the documentation easily readable and allows godoc to generate a synopsis.
4. Use Full Sentences:
   Rule: Write doc comments using complete, grammatically correct sentences.
   Reason: Full sentences improve readability and formatting in the generated documentation.
5. Explain the "Why," Not Just the "What":
   Rule: Focus on explaining the purpose, intent, and non-obvious aspects of your code.
   Reason: Redundant comments that simply rephrase the code's action are unhelpful.
6. Use Package-Level Documentation:
   Rule: Include a package-level comment (usually in doc.go) that provides an overview of the package's purpose and contents.
   Reason: This gives users a high-level understanding of what the package does.
7. Add Examples:
   Rule: Include runnable examples in your tests (Example... functions) to demonstrate how to use your code.
   Reason: These examples are automatically included in the godoc output and serve as practical demonstrations.
8. Use gofmt:
   Rule: Always run gofmt on your code.
   Reason: This ensures consistent formatting, which is crucial for godoc and overall code readability.
9. Be Concise and Clear:
   Rule: Write doc comments that are easy to understand and avoid unnecessary complexity.
   Reason: Concise and clear documentation is more likely to be read and used.
10. Include "BUG(who)" Comments (when appropriate):
    Rule: For known issues, use the // BUG(who): format.
    Reason: godoc specifically recognizes these comments and includes them in a separate "Bugs" section.
    By adhering to these principles, you'll create clear, informative documentation that is easily discoverable and usable via godoc. 
