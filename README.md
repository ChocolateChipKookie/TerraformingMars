
# Terraforming Mars Tracker

A web application for tracking Terraforming Mars game statistics.

## Tech Stack

- **Frontend**: React 19 with Vite
- **Backend**: Go
- **Routing**: React Router v7

## Development Setup

### Prerequisites

- Node.js (for frontend development)
- Go (for backend development)
- npm or similar package manager

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/ChocolateChipKookie/TerraformingMars.git
   cd TerraformingMars
   ```

2. Install frontend dependencies:
   ```bash
   npm install
   ```

3. Install backend dependencies:
   ```bash
   cd backend
   go mod download
   cd ..
   ```

### Running the Application

**Frontend (Development)**:
```bash
npm run dev
```

**Backend**:
```bash
cd backend
go run main.go
```

**Build for Production**:
```bash
npm run build
```

## Contributing

### Development Workflow

1. Create a new branch for your feature or bug fix
2. Make your changes following the project's coding standards
3. Test your changes thoroughly
4. Commit your changes with clear, descriptive commit messages
5. Push to your branch and create a pull request

### Code Style

- Follow standard Go conventions for backend code
- Use React best practices for frontend components
- Ensure code is properly formatted before committing

### Pull Request Guidelines

- Provide a clear description of the changes
- Reference any related issues
- Ensure all tests pass
- Keep PRs focused on a single feature or fix

## TODO

### High Priority

- [ ] Make it production ready
  - Go backend should serve the webpage
- [ ] Add migration framework for database schema changes

### Features

- [ ] **Player Ratings Page**
  - Calculate all player ratings every time the page is opened (prevents issues with deleting/updating games)
  - Add caching mechanism for performance optimization


# DOCKER

```




```


