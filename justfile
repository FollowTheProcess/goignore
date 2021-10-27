PROJECT_NAME := "goignore"
PROJECT_PATH := "github.com/FollowTheProcess/goignore"
PROJECT_BIN := "./bin"
PROJECT_ENTRY_POINT := "."
COVERAGE_DATA := "coverage.out"
COVERAGE_HTML := "coverage.html"
GORELEASER_DIST := "dist"
COMMIT_SHA := `git rev-parse HEAD`
VERSION_LDFLAG := "main.version"
COMMIT_LDFLAG := "main.commit"

# By default print the list of recipes
_default:
    @just --list

# Tidy up dependencies in go.mod and go.sum
tidy:
    go mod tidy

# Compile the project binary
build: tidy fmt
    go build -ldflags="-s -w -X {{ VERSION_LDFLAG }}=dev -X {{ COMMIT_LDFLAG }}={{ COMMIT_SHA }}" -o {{ PROJECT_BIN }}/{{ PROJECT_NAME }} {{ PROJECT_ENTRY_POINT }}

# Run go fmt on all project files
fmt:
    gofumpt -extra -s -w .

# Run all project unit tests
test *flags: fmt
    go test -race ./... {{ flags }}

# Lint the project and auto-fix errors if possible
lint: fmt
    golangci-lint run --fix

# Calculate test coverage and render the html
cover:
    go test -race -cover -coverprofile={{ COVERAGE_DATA }} ./...
    go tool cover -html={{ COVERAGE_DATA }} -o {{ COVERAGE_HTML }}
    open {{ COVERAGE_HTML }}

# Remove build artifacts and other project clutter
clean:
    go clean ./...
    rm -rf {{ PROJECT_NAME }} {{ PROJECT_BIN }} {{ COVERAGE_DATA }} {{ COVERAGE_HTML }} {{ GORELEASER_DIST }}

# Run unit tests and linting in one go
check: test lint

# Run all recipes (other than clean) in a sensible order
all: build test lint cover

# Install the project on your machine
install:
    go install {{ PROJECT_ENTRY_POINT }}

# Uninstall the project from your machine
uninstall:
    rm -rf $GOBIN/{{ PROJECT_NAME }}
