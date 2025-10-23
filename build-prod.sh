#!/bin/bash
set -e

echo "Building frontend"
cd frontend
npm run build
cd ..

echo "Building backend"
cd backend
go build -o tm-server
cd ..

echo "Preparing build directory"
if [ -d "build" ]; then
  echo "Removing existing build directory"
  rm -rf build
fi

mkdir -p build

echo "Copying artifacts to build directory"
cp -r frontend/dist build/
cp backend/tm-server build/
mkdir build/data

echo ""
echo "Build complete!"
echo "Build artifacts are in the 'build/' directory:"
echo "  - build/dist/                    (frontend static files)"
echo "  - build/terraforming-mars-server (backend binary)"
echo "  - build/data/                    (database directory, empty)"
echo ""
echo "To run the production server:"
echo "  cd build && ./terraforming-mars-server"
echo ""
echo "The server will be available at http://localhost:8080"
