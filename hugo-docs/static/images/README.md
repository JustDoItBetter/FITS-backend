# Documentation Images

This directory contains diagrams and images for the documentation.

## Creating Diagrams

### PlantUML

Install PlantUML:

```bash
# macOS
brew install plantuml

# Ubuntu
sudo apt install plantuml
```

Create diagrams:

```plantuml
@startuml
!theme plain

title Authentication Flow

actor Client
participant API
database Database

Client -> API: POST /auth/login
API -> Database: Validate credentials
Database --> API: User data
API -> API: Generate JWT
API -> Database: Create session
API --> Client: Return tokens

@enduml
```

Generate PNG:

```bash
plantuml diagram.puml
```

### Mermaid

Create diagrams in markdown:

```mermaid
graph TD
    A[Client Request] --> B{Authenticated?}
    B -->|Yes| C[Check Permissions]
    B -->|No| D[Return 401]
    C -->|Allowed| E[Process Request]
    C -->|Denied| F[Return 403]
```

Render using:
- GitHub (automatic)
- Mermaid Live Editor (https://mermaid.live)
- VS Code Mermaid extension

## Diagram Guidelines

1. Use consistent colors:
   - Blue: External systems
   - Green: Success paths
   - Red: Error paths
   - Gray: Internal components

2. Keep diagrams simple and focused

3. Include legends when necessary

4. Use SVG format when possible (better scaling)

5. Optimize PNG images:
   ```bash
   optipng -o7 diagram.png
   ```

## Recommended Tools

- **PlantUML**: Sequence diagrams, class diagrams
- **Mermaid**: Flow charts, sequence diagrams
- **Excalidraw**: Hand-drawn style diagrams
- **Draw.io**: General purpose diagrams
