# Action Based Access Control [ABAC]


## Tech Stack
- Golang
- Docker


## Entities

![entities](docs/diagrams/abac-entities.jpg)


## Routes

![entities](docs/diagrams/abac-routes.jpg)


## Project Structure

### `/pkg` - The Framework
    - `/pkg/utils/`, `/pkg/config/`, ...
    - No dependencies on `/cmd`
    - Can be imported by external programs
### `/cmd` - The Programs
    - `/cmd/abac/`, `/cmd/service/`, ...
    - Domain specific logic stays close to the main() func


## Todo:
- [X] setup migrations
- [X] user model implementation
- [X] CRUD for user model
- [X] password keeping mech
- [X] field level validation
- [X] password hiding
- [X] check for required fields on each scheme
- [x] basic JWT auth
- [X] jwt support for all routes
- [X] refactor and restructure the whole project
- [X] refactoring teardown
- [X] server graceful shutdown
- [X] basic logging
- [x] add auth middleware
- [x] add request body validators
- [ ] add decode interface to all request schemas, so decode and validation will be transferred there
- [ ] add Group and Action entities
- [ ] add migration schemas for Group and Action
- [ ] CRUD for Group entity
- [ ] CRUD for Action entity
- [ ] add extension for jwt token payload schema to handle needed BL
- [ ] extend existing AC to pass new payload schema
