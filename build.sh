#!/bin/bash

# Build script for terraform-state-visualizer
# This script builds the application for multiple platforms and creates releases

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="terraform-state-visualizer"
VERSION=${1:-"dev"}
BUILD_DIR="build"
DIST_DIR="dist"

# Print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Clean build directories
clean() {
    print_status "Cleaning build directories..."
    rm -rf ${BUILD_DIR} ${DIST_DIR}
    mkdir -p ${BUILD_DIR} ${DIST_DIR}
    print_success "Build directories cleaned"
}

# Build for a specific platform
build_platform() {
    local os=$1
    local arch=$2
    local ext=$3
    
    print_status "Building for ${os}/${arch}..."
    
    local output_name="${APP_NAME}"
    if [ "${os}" = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    GOOS=${os} GOARCH=${arch} go build \
        -ldflags "-X main.Version=${VERSION} -X main.BuildTime=$(date -u '+%Y-%m-%d_%H:%M:%S') -X main.GitCommit=$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')" \
        -o ${BUILD_DIR}/${os}_${arch}/${output_name} .
    
    print_success "Built ${os}/${arch}"
}

# Build all platforms
build_all() {
    print_status "Building for all platforms..."
    
    # Linux
    build_platform "linux" "amd64" ""
    build_platform "linux" "arm64" ""
    
    # macOS
    build_platform "darwin" "amd64" ""
    build_platform "darwin" "arm64" ""
    
    # Windows
    build_platform "windows" "amd64" ".exe"
    build_platform "windows" "arm64" ".exe"
    
    print_success "All platforms built"
}

# Create distribution packages
create_packages() {
    print_status "Creating distribution packages..."
    
    for platform_dir in ${BUILD_DIR}/*; do
        if [ -d "${platform_dir}" ]; then
            platform=$(basename ${platform_dir})
            print_status "Packaging ${platform}..."
            
            # Create tar.gz for Unix-like systems
            if [[ ${platform} == *"windows"* ]]; then
                # Create zip for Windows
                cd ${platform_dir}
                zip -r ../../${DIST_DIR}/${APP_NAME}-${VERSION}-${platform}.zip .
                cd - > /dev/null
            else
                # Create tar.gz for Unix-like systems
                tar -czf ${DIST_DIR}/${APP_NAME}-${VERSION}-${platform}.tar.gz -C ${platform_dir} .
            fi
            
            print_success "Created package for ${platform}"
        fi
    done
    
    print_success "All packages created"
}

# Build Docker image
build_docker() {
    print_status "Building Docker image..."
    
    docker build -t ${APP_NAME}:${VERSION} .
    docker tag ${APP_NAME}:${VERSION} ${APP_NAME}:latest
    
    print_success "Docker image built: ${APP_NAME}:${VERSION}"
}

# Run tests
run_tests() {
    print_status "Running tests..."
    
    go test -v ./...
    
    print_success "Tests passed"
}

# Show help
show_help() {
    echo "Usage: $0 [VERSION] [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  clean       Clean build directories"
    echo "  build       Build for all platforms"
    echo "  docker      Build Docker image"
    echo "  test        Run tests"
    echo "  package     Create distribution packages"
    echo "  all         Run clean, test, build, package, and docker (default)"
    echo "  help        Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                    # Build all with version 'dev'"
    echo "  $0 v1.0.0             # Build all with version 'v1.0.0'"
    echo "  $0 v1.0.0 build        # Build only for version 'v1.0.0'"
    echo "  $0 v1.0.0 docker       # Build Docker image only"
}

# Main execution
main() {
    local command=${2:-"all"}
    
    case ${command} in
        "clean")
            clean
            ;;
        "build")
            clean
            build_all
            ;;
        "docker")
            build_docker
            ;;
        "test")
            run_tests
            ;;
        "package")
            clean
            build_all
            create_packages
            ;;
        "all")
            clean
            run_tests
            build_all
            create_packages
            build_docker
            ;;
        "help")
            show_help
            ;;
        *)
            print_error "Unknown command: ${command}"
            show_help
            exit 1
            ;;
    esac
    
    print_success "Build completed successfully!"
    
    # Show summary
    if [ "${command}" = "all" ] || [ "${command}" = "package" ]; then
        echo ""
        print_status "Created packages:"
        ls -la ${DIST_DIR}/
    fi
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go 1.25.3 or later."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.25.3"
if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    print_error "Go version ${GO_VERSION} is too old. Please install Go ${REQUIRED_VERSION} or later."
    exit 1
fi

# Run main function
main "$@"
