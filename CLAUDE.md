# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

bc-cli is a Go-based CLI application for Butler Coffee that enables users to discover coffee knowledge, manage subscriptions, and purchase coffee directly from the terminal. It uses the Cobra framework for command structure and interacts with the Butler Coffee API at butler.coffee.

## Development Commands

### Build
```bash
make compile         # Builds the binary as bc-cli
go build -o bc-cli . # Direct build command
```

### Development Setup
```bash
make install  # Installs dependencies, sets up pre-commit hooks
```

### Dependency Management
```bash
make upgrade  # Updates Go dependencies and runs pip-upgrade
go mod tidy   # Cleans up go.mod and go.sum
```

### API Schema Reference
The backend API schema is available in OpenAPI 3.0 format:
```bash
# View the full API schema (requires backend running)
curl http://localhost:8000/api/schema/

# Set BASE_HOSTNAME for local development
export BASE_HOSTNAME=http://localhost:8000

# The CLI automatically uses BASE_HOSTNAME if set, otherwise defaults to https://api.butler.coffee
```

**Important**: Always check the API schema when implementing or modifying API endpoints to ensure correct paths, methods, and request/response formats.

## Architecture

### Command Structure
- **Entry point**: `main.go` → `cmd.Execute()`
- **Commands**: Cobra commands defined in `cmd/` directory
  - `cmd/root.go` - Root command and CLI initialization
  - `cmd/login.go` - User authentication
  - `cmd/signup.go` - User registration
  - `cmd/subscribe.go` - Interactive subscription flow with order configuration
  - `cmd/manage.go` - Subscription management (pause, resume, update, cancel)

### API Layer (`api/`)
- **`api/client.go`**: Core HTTP client with automatic token refresh logic
  - `doRequest()` handles token expiration and auto-retry on 401
  - `handleResponse()` parses structured API errors from Django backend
- **`api/auth.go`**: Login, registration, and token refresh endpoints
- **`api/subscriptions.go`**: Subscription management API methods
  - Listing and retrieving subscriptions with order details
  - `Subscription` model includes `TotalQuantityKg` and `LineItems` for displaying current configuration
  - `GetTotalQuantity()` helper method converts string quantity to int
  - **Endpoints used**:
    - `GET /api/core/v1/subscriptions` - List all user subscriptions
    - `GET /api/core/v1/subscriptions/{id}/preferences` - Get subscription with order details
    - `POST /api/core/v1/subscriptions/{id}/pause` - Pause subscription (no request body)
    - `POST /api/core/v1/subscriptions/{id}/resume` - Resume paused subscription
    - `PATCH /api/core/v1/subscriptions/{id}/preferences` - Update subscription (expects `preferences` field, not `line_items`)
    - `POST /api/core/v1/subscriptions/{id}/cancel` - Cancel subscription
    - `GET /api/core/v1/subscriptions/available` - Get available subscription tiers
- **`api/orders.go`**: Order creation, checkout sessions, and order management

### Configuration (`config/`)
- **`config/config.go`**: Manages user config stored at `~/.butler-coffee/config.json`
  - Handles access tokens, refresh tokens, and expiration tracking
  - Token expiration uses 30-second safety margin
  - Supports `BASE_HOSTNAME` environment variable for API URL override
  - Default API URL: `https://api.butler.coffee`
  - **Quantity limits**: Configurable min/max kg per month for subscriptions
    - `min_quantity_kg`: Minimum kilos per month (default: 1)
    - `max_quantity_kg`: Maximum kilos per month (default: 10)
    - These values can be customized by editing `~/.butler-coffee/config.json`

### Templates (`templates/`)
- **`templates/templates.go`**: Text template rendering with custom functions
- **`templates/auth.go`**: Login, signup, and authentication UI templates
- **`templates/subscription.go`**: Subscription and order flow templates, including:
  - Subscription purchase and checkout templates
  - Subscription management templates (pause, resume, update, cancel)
  - Status-specific messaging and warnings

## Key Behaviors

### Token Management
- Access tokens and refresh tokens are automatically managed by `api/client.go`
- On 401 response, the client attempts token refresh and retries the request once
- Token expiration is checked proactively before authenticated requests
- If refresh token is expired, user must login again

### Subscription Purchase Flow
The subscription purchase flow in `cmd/subscribe.go` is a multi-step interactive process:
1. Display available subscriptions with active status for authenticated users
2. Allow user to select a tier and view detailed information
3. Configure order preferences:
   - Total quantity (kg per month)
   - Choice between uniform or split preferences
   - For each preference: grind type (whole bean vs ground) and brewing method
4. Show order summary and confirm
5. Create order via API
6. Generate Stripe checkout session and open browser
7. Poll for payment completion (5 minute timeout)

### Subscription Management Flow
The subscription management flow in `cmd/manage.go` provides comprehensive control:
1. List all user subscriptions with status indicators (active, paused, cancelled)
2. Select a subscription to manage
3. Display subscription details including:
   - Status, started date, and expiration/resume date
   - **Billing information**: Price, currency, and billing period (fetched from SubscriptionPlan)
   - Current order configuration (total kg/month and default preferences with grind settings)
   - Formatted display: "X kg → Whole beans for Espresso" or "X kg → Ground for V60 (medium)"
4. Show available actions based on subscription status:
   - **Active subscriptions**: pause, update preferences, or cancel
   - **Paused subscriptions**: resume, update preferences, or cancel
   - **Cancelled subscriptions**: no actions available
5. **Pause functionality**:
   - Pauses subscription indefinitely (backend doesn't support scheduled resume yet)
   - User can manually resume at any time
   - Note: See TEST_FAILURES.md for backend limitations
6. **Cancel flow**:
   - First offers pause as a better alternative
   - Requires double confirmation for permanent cancellation
7. **Update preferences**:
   - Shows current configuration before update
   - Reuses order configuration flow from subscribe command
   - Allows changing quantity and grind/brewing preferences

### Error Handling
- API errors follow Django REST structure with `meta.errors[]` containing field-level validation errors
- Client extracts and formats field errors for user-friendly display
- Generic errors fallback to `meta.message` or raw response

## Important Patterns

### Interactive Prompts
Uses `github.com/manifoldco/promptui` for interactive terminal UI:
- `promptui.Select` for menu selections
- `promptui.Prompt` for text input with validation
- Custom templates for styled output

### Decimal Field Handling
Django backend returns decimal fields as strings (e.g., `"2.5"`):
- API responses use string types for `quantity_kg`, `total_quantity_kg`
- Helper methods like `GetQuantity()` and `GetTotalQuantity()` parse to float64
- API requests use `int` types which backend accepts

### Timestamp Handling
`utils.FormatTimestamp()` handles multiple timestamp formats from the backend:
- Unix milliseconds as strings (e.g., `"1764427190000"`)
- ISO 8601 / RFC3339 formats
- Simple date strings (e.g., `"2006-01-02"`)
- All timestamps are displayed as: "January 2, 2006 at 3:04 PM"

### Browser Integration
`cmd/subscribe.go` includes `openBrowser()` function supporting macOS, Linux, and Windows for opening checkout URLs.

## Configuration Files

### go.mod
- Go version: 1.25.4
- Main framework: `github.com/spf13/cobra v1.10.1`
- UI library: `github.com/manifoldco/promptui v0.9.0`

### Makefile
Defines common tasks for building, installing dependencies, and managing pre-commit hooks.

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test ./... -v

# Run tests for a specific package
go test ./api -v
go test ./utils -v

# Run a specific test
go test ./api -v -run TestListSubscriptions

# Run tests with coverage
go test ./... -cover
```

### Test Structure

#### API Tests (`api/subscriptions_test.go`)
Comprehensive tests for subscription management API methods using httptest:
- **TestListSubscriptions** - Verifies GET /subscriptions endpoint
- **TestGetSubscription** - Tests GET /subscriptions/{id}/preferences with line items
- **TestPauseSubscription** - Tests pause with and without resume date
- **TestResumeSubscription** - Verifies resume functionality
- **TestCancelSubscription** - Tests cancellation
- **TestUpdateSubscription** - Verifies PATCH /preferences endpoint
- **TestSubscriptionGetTotalQuantity** - Tests quantity parsing (string to int)

Each test:
- Creates a mock HTTP server
- Verifies correct endpoint paths and HTTP methods
- Validates request/response structure
- Tests error handling

#### Utility Tests (`utils/date_test.go`)
Tests for timestamp formatting:
- **TestFormatTimestamp** - Tests multiple timestamp formats (unix ms, RFC3339, dates)
- **TestFormatTimestampBoundaries** - Validates year 2000-2100 range for unix milliseconds
- Tests edge cases (empty strings, invalid formats, boundary values)

### Test Coverage

Current test coverage focuses on:
- ✅ All subscription API methods
- ✅ Timestamp/date formatting
- ✅ Request/response parsing
- ✅ Error handling

### Writing New Tests

When adding new features:
1. Add tests in `*_test.go` files alongside the code
2. Use `httptest.NewServer` for API endpoint tests
3. Test both success and error cases
4. Verify request structure (method, path, body)
5. Validate response parsing
