# LensLocked - Photo Gallery Application

A modern photo gallery application built with Go, allowing users to create and manage photo galleries. Users can sign up, create galleries, upload photos, and share their galleries with others.

## Features

- 🔐 User authentication (signup, signin, signout)
- 📸 Create and manage photo galleries
- 🖼️ Upload and manage images
- 🌓 Dark/Light theme support
- 🔒 Protected routes and resources
- 🗃️ PostgreSQL database integration
- 🖥️ Docker support for development
- 🚀 Heroku deployment ready

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- PostgreSQL client (for database access)
- Make (for running commands)

## Quick Start

1. **Clone the repository**
```bash
git clone https://github.com/enlistedmango/lenslocked.git
cd lenslocked
```

2. **Set up environment variables**
```bash
cp .env.example .env
# Edit .env with your configuration:
# - Database settings
# - CSRF key
# - Session key
# - FiveManage API key
```

3. **Start the application with Docker**
```bash
make docker-build  # Build the Docker images
make docker-run   # Start the application
```

The application will be available at `http://localhost:3000`

## Development

### Local Development
```bash
# Start only the database
make local-db

# Run the application locally
make run
```

### Testing
```bash
# Run all tests
make test-all

# Run specific tests
make test-local      # Test local environment
make test-prod       # Test production environment
make test-integration  # Run integration tests
```

## Deployment

### Heroku Deployment

1. **Create and configure Heroku app**
```bash
# Create new Heroku app
make heroku-create

# Configure environment variables
make heroku-config

# Set your FiveManage API key
heroku config:set FIVEMANAGE_API_KEY=your-key-here
```

2. **Deploy**
```bash
make deploy-prod
```

## Docker Commands

```bash
# Build images
make docker-build

# Start services
make docker-run

# Stop services
make docker-down

# View logs
docker-compose logs
```

## Project Structure

```
.
├── controllers/     # Request handlers
├── middleware/     # HTTP middleware
├── models/         # Database models
├── services/       # Business logic
├── templates/      # HTML templates
├── views/          # View rendering logic
├── migrations/     # Database migrations
├── scripts/        # Utility scripts
└── static/         # Static assets
```

## Testing Strategy

The project includes comprehensive testing:

- **Local Environment Tests**: Verify Docker services, database connection, and web service functionality
- **Integration Tests**: Test all endpoints and their expected behaviors
- **Production Tests**: Verify Heroku deployment and production environment
- **GitHub Actions**: Automated testing on push and pull requests

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Based on [usegolang.com](https://www.usegolang.com/) course by Jon Calhoun
- Extended with additional features:
  - Docker support
  - Dark/Light theme
  - Automated testing
  - CI/CD pipeline

