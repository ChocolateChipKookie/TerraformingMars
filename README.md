# Terraforming Mars Tracker

A web application for tracking Terraforming Mars game statistics.


## Backstory

It is a complete rewrite of the prvious terraform mars tracking website which had no backend, but stored the data in github directly. This version of the website actually has a backend, but has the same style as the previous iteration.

The previous version broke at some point when Github made some changes how API keys work, and I never had the actual will to fix it. It was very buggy, and required too much work for me to actually get it working correctly. But some time this year we filled up the small notebook we used for tracking games, so the time has come to resurrect the project (also I had a little helper on my side that made rewriting the front-end a bit less painful (it's not the peak of engineering, actually it is quite the slop, but it's working)).

The front-end of the website is the full domain of Claude, I am not really sure why some things work there, but it works and gets the job done.
Claude also had fingers in the backend, but I had some idea of what I wanted to do so it is in a much cleaner state (there are still some issues that happened, but I would guess much less then in the front-end).

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
   cd frontend
   npm install
   cd ..
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
cd frontend
npm run dev
```

**Backend**:
```bash
cd backend
go run main.go
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

- [ ] Add migration framework for database schema changes
- [ ] **Player Ratings Page**
  - Calculate all player ratings every time the page is opened (prevents issues with deleting/updating games)
  - Add caching mechanism for performance optimization

## Docker

### Running with Docker

The application can be run as a Docker container that includes both the frontend and backend services in a single container.

#### Quick Start

1. Use the provided script to build and run the container:
   ```bash
   ./run-docker.sh [port]
   ```

   Examples:
   ```bash
   ./run-docker.sh        # Runs on default port 8080
   ./run-docker.sh 3000   # Runs on port 3000
   ```

   This script will:
   - Build the Docker image
   - Stop any existing container with the same name
   - Create a data directory (if it doesn't exist)
   - Start the container with appropriate settings

2. Access the application at `http://localhost:[port]` (default: `http://localhost:8080`)

#### Manual Docker Commands

Build the image:
```bash
docker build -t terraforming-mars:latest .
```

Run the container:
```bash
docker run -d \
  --name tm-server \
  --restart unless-stopped \
  -p 8080:8080 \
  -v ./data:/data \
  terraforming-mars:latest
```

#### Docker Configuration

- **Port**: The application runs on port 8080 inside the container
- **Data Persistence**: Game data is persisted in the `./data` directory mounted as a volume
- **Resource Limits**: By default limited to 256MB of memory
