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


