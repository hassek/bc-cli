# Integration Tests

This document describes how to run end-to-end integration tests for the Butler Coffee CLI against a real backend API.

## Quick Start

**Want to test everything? Run the interactive test:**

```bash
# 1. Start your backend
cd /path/to/backend && python manage.py runserver

# 2. Run the interactive test (in CLI directory)
export BASE_HOSTNAME=http://localhost:8000
export RUN_INTEGRATION_TESTS=true
export INTERACTIVE=true
go test -v -run TestE2EInteractiveFullLifecycle -timeout 30m
```

This will test the complete lifecycle: registration → order → payment → subscription management → cancellation.
The test will open Stripe checkout in your browser - use test card `4242 4242 4242 4242`.

## Overview

The integration tests are designed to test the entire user flow without mocking any API calls:

1. **User Registration & Authentication**: Creating new users and logging in
2. **Subscription Discovery**: Fetching available subscription plans
3. **Order Creation**: Creating orders with coffee preferences
4. **Checkout & Payment**: Completing Stripe checkout with test cards
5. **Subscription Management**: Pausing, resuming, updating, and cancelling subscriptions

## Prerequisites

### For Local Testing

1. **Backend Running**: Start the Butler Coffee backend locally
   ```bash
   # In your backend repository
   python manage.py runserver
   ```

2. **Environment Variables**: Set the API endpoint
   ```bash
   export BASE_HOSTNAME=http://localhost:8000
   ```

3. **Database**: Ensure your local database is running and migrations are applied

### For Production/Staging Testing

1. **Environment Variables**: Set the API endpoint (or omit to use production)
   ```bash
   export BASE_HOSTNAME=https://api.butler.coffee  # or your staging URL
   ```

2. **Warning**: Be cautious when running tests against production

## Running the Tests

### Test 1: Interactive Full Lifecycle (TestE2EInteractiveFullLifecycle) ⭐ RECOMMENDED

**This is the most comprehensive test that verifies the complete subscription lifecycle including payment!**

This interactive test walks through the entire user journey from registration to cancellation, including completing a real Stripe checkout. The test will:
1. Create a new user account
2. Create an order with coffee preferences
3. Generate a Stripe checkout session
4. **Open your browser for you to complete payment with a test card**
5. Poll the API to verify the subscription was created
6. Test all subscription management operations (pause, resume, update)
7. Cancel the subscription

**Requirements:**
- Local backend running with Stripe test mode configured
- About 5-10 minutes to complete (includes manual payment step)

**Run the test:**
```bash
# Set environment variables
export BASE_HOSTNAME=http://localhost:8000
export RUN_INTEGRATION_TESTS=true
export INTERACTIVE=true

# Run the test (note the 30m timeout to allow time for manual steps)
go test -v -run TestE2EInteractiveFullLifecycle -timeout 30m
```

**What to expect:**
1. The test will create a new user and order
2. Your browser will open to Stripe checkout
3. Use the Stripe test card:
   - Card: `4242 4242 4242 4242`
   - Expiry: Any future date (e.g., `12/34`)
   - CVC: Any 3 digits (e.g., `123`)
   - ZIP: Any 5 digits (e.g., `12345`)
4. Complete the payment
5. The test will automatically detect the subscription and continue testing all operations
6. Finally, it will cancel the subscription

**Expected output:**
```
=== RUN   TestE2EInteractiveFullLifecycle
=======================================================================
INTERACTIVE INTEGRATION TEST - FULL SUBSCRIPTION LIFECYCLE
=======================================================================

This test will:
  1. Create a new user account
  2. Create an order and checkout session
  3. Open Stripe checkout in your browser
  4. Wait for you to complete payment
  5. Test subscription management (pause, resume, update)
  6. Cancel the subscription

Test credentials: test_user_1234567890_1234 / TestPassword123!@#
=======================================================================

▶ Step 1: Creating user account...
  User ID: c3e84c12-67f4-4975-97dc-ab5a4a3ca9bc
  Username: test_user_1234567890_1234
  Email: test_1234567890_1234@butler.test
✓ User created: c3e84c12-67f4-4975-97dc-ab5a4a3ca9bc

▶ Step 2: Fetching available subscriptions...
  Available tiers:
    [1] Butler - 49.00 EUR/monthly
    [2] Collection - 79.00 EUR/monthly
✓ Selected tier: butler

▶ Step 3: Creating order with coffee preferences...
  Order ID: 7133b1a5-5efb-432f-8481-f56fbcebe9f4
  Tier: butler
  Total quantity: 5.0 kg/month
  Line items:
    [1] v60 → 2.0 kg (ground)
    [2] espresso → 3.0 kg (whole_bean)
✓ Order created: 7133b1a5-5efb-432f-8481-f56fbcebe9f4

▶ Step 4: Creating checkout session and processing payment...
  Checkout session created
  Session ID: cs_test_...

  ┌─────────────────────────────────────────────────────────────┐
  │ OPENING BROWSER FOR STRIPE CHECKOUT...                     │
  └─────────────────────────────────────────────────────────────┘

  ✓ Browser opened

  Use Stripe test card: 4242 4242 4242 4242
  Expiry: Any future date (e.g., 12/34)
  CVC: Any 3 digits (e.g., 123)
  ZIP: Any 5 digits (e.g., 12345)

  Waiting for payment completion...
  (Polling every 3 seconds, timeout in 5 minutes)

  [Attempt 1] No active subscription yet, retrying in 3 seconds...
  [Attempt 2] No active subscription yet, retrying in 3 seconds...

  ✓ Payment verified!
  ✓ Subscription created: sub_xxxxx
✓ Payment completed! Subscription ID: sub_xxxxx

▶ Step 5: Verifying subscription was created...
  ID: sub_xxxxx
  Tier: butler
  Status: active
  Total quantity: 5 kg/month
  Preferences:
    [1] v60 → 2 kg (ground)
    [2] espresso → 3 kg (whole_bean)
✓ Subscription verified

▶ Step 6: Pausing subscription...
  Status changed to: paused
✓ Subscription paused

▶ Step 7: Resuming subscription...
  Status changed to: active
✓ Subscription resumed

▶ Step 8: Updating subscription preferences...
  Updated total quantity: 8 kg/month
  New preferences:
    [1] french_press → 5 kg (ground)
    [2] aeropress → 3 kg (whole_bean)
✓ Subscription updated

▶ Step 9: Restoring original preferences...
  Restored to 5 kg/month
✓ Preferences restored

▶ Step 10: Cancelling subscription...
  Status changed to: cancelled
  Subscription is now permanently cancelled
✓ Subscription cancelled

=======================================================================
✓ ALL TESTS PASSED - COMPLETE LIFECYCLE VERIFIED
=======================================================================
--- PASS: TestE2EInteractiveFullLifecycle (XX.XXs)
```

### Test 2: Automated User Flow (TestE2EFullFlow)

This automated test creates a new user and runs through the subscription purchase flow up to checkout (but doesn't complete payment).

**What it tests:**
- User registration
- User login
- Fetching available subscriptions
- Creating an order with multiple preferences
- Creating a checkout session
- Listing subscriptions (verifies empty list since payment wasn't completed)

**Run the test:**
```bash
# Set environment variables
export BASE_HOSTNAME=http://localhost:8000
export RUN_INTEGRATION_TESTS=true

# Run the test
go test -v -run TestE2EFullFlow
```

**Expected output:**
```
=== RUN   TestE2EFullFlow
    integration_test.go:XX: Testing against: http://localhost:8000
    integration_test.go:XX: Generated test credentials: username=test_user_1234567890_1234, email=test_1234567890_1234@butler.test
=== RUN   TestE2EFullFlow/UserRegistration
    integration_test.go:XX: ✓ User registered successfully: ID=usr_xxxxx
    integration_test.go:XX: ✓ Received access token: eyJhbGciOiJIUzI1NiIs...
=== RUN   TestE2EFullFlow/UserLogin
    integration_test.go:XX: ✓ User logged in successfully: UserID=usr_xxxxx
=== RUN   TestE2EFullFlow/GetAvailableSubscriptions
    integration_test.go:XX: ✓ Found 3 available subscription(s)
    integration_test.go:XX:   [1] Tier: butler, Name: Butler, Price: 29.99 USD/month
=== RUN   TestE2EFullFlow/CreateOrder
    integration_test.go:XX: ✓ Order created successfully: ID=ord_xxxxx
    integration_test.go:XX:   [1] espresso → 3.0 kg (whole_bean)
    integration_test.go:XX:   [2] v60 → 2.0 kg (ground)
=== RUN   TestE2EFullFlow/CreateCheckoutSession
    integration_test.go:XX: ✓ Checkout session created successfully
    integration_test.go:XX: Note: Automated tests cannot complete Stripe payment flow
--- PASS: TestE2EFullFlow (X.XXs)
```

### Test 3: Subscription Management (TestE2ESubscriptionManagement)

This test requires an existing active subscription and tests all management operations.

**What it tests:**
- Fetching subscription details with preferences
- Pausing a subscription
- Resuming a paused subscription
- Updating subscription preferences
- Restoring original preferences
- (Optional) Cancelling a subscription

**Setup:**
1. First, you need a subscription with completed payment. You can:
   - Manually complete the Stripe checkout from TestE2EFullFlow
   - Use Stripe test mode with test cards
   - Create a subscription directly in the database for testing

2. Get the subscription ID:
   ```bash
   # Login to the CLI
   ./bc-cli login

   # List your subscriptions
   ./bc-cli subscriptions
   # Copy the subscription ID
   ```

**Run the test:**
```bash
# Set environment variables
export BASE_HOSTNAME=http://localhost:8000
export RUN_INTEGRATION_TESTS=true
export TEST_SUBSCRIPTION_ID=sub_xxxxxxxxxxxxx

# Run the test
go test -v -run TestE2ESubscriptionManagement
```

**Expected output:**
```
=== RUN   TestE2ESubscriptionManagement
    integration_test.go:XX: Testing against: http://localhost:8000
=== RUN   TestE2ESubscriptionManagement/GetSubscription
    integration_test.go:XX: ✓ Subscription retrieved: ID=sub_xxxxx
    integration_test.go:XX: ✓ Tier: butler, Status: active
    integration_test.go:XX: ✓ Total quantity: 5 kg/month
=== RUN   TestE2ESubscriptionManagement/PauseSubscription
    integration_test.go:XX: ✓ Subscription paused successfully: ID=sub_xxxxx
=== RUN   TestE2ESubscriptionManagement/ResumeSubscription
    integration_test.go:XX: ✓ Subscription resumed successfully: ID=sub_xxxxx
=== RUN   TestE2ESubscriptionManagement/UpdateSubscription
    integration_test.go:XX: ✓ Subscription updated successfully: ID=sub_xxxxx
    integration_test.go:XX: ✓ New total quantity: 8 kg/month
    integration_test.go:XX: ✓ Verified updated preferences:
    integration_test.go:XX:   [1] french_press → 5 kg (ground)
    integration_test.go:XX:   [2] aeropress → 3 kg (whole_bean)
=== RUN   TestE2ESubscriptionManagement/RestoreOriginalPreferences
    integration_test.go:XX: ✓ Subscription restored to original preferences: ID=sub_xxxxx
--- PASS: TestE2ESubscriptionManagement (X.XXs)
```

## Running All Integration Tests

```bash
export BASE_HOSTNAME=http://localhost:8000
export RUN_INTEGRATION_TESTS=true
export TEST_SUBSCRIPTION_ID=sub_xxxxxxxxxxxxx  # Optional, only for subscription management test

go test -v -run "TestE2E"
```

## Debugging Failed Tests

### Enable Verbose Logging

The API client includes request/response logging. To see detailed output:

```bash
go test -v -run TestE2EFullFlow 2>&1 | tee test_output.log
```

### Common Issues

1. **Connection Refused**
   - Verify backend is running: `curl http://localhost:8000/api/core/v1/health`
   - Check BASE_HOSTNAME is set correctly

2. **401 Unauthorized**
   - Token may have expired between test runs
   - Delete config and re-run: `rm -rf ~/.butler-coffee/config.json`

3. **Validation Errors**
   - Check backend schema matches expected request format
   - Review API schema: `curl http://localhost:8000/api/schema/`

4. **Test User Already Exists**
   - Tests generate random usernames, but collisions can happen
   - Manually delete test users from database if needed

## Test Coverage

### What IS Tested
- ✅ User registration with email/password
- ✅ User login and token retrieval
- ✅ Fetching available subscription plans
- ✅ Creating orders with multiple line items
- ✅ Generating Stripe checkout sessions
- ✅ Subscription retrieval with preferences
- ✅ Pausing/resuming subscriptions
- ✅ Updating subscription preferences
- ✅ Cancelling subscriptions

### What IS NOT Tested
- ❌ Completing Stripe payment flow (requires manual intervention)
- ❌ Webhook processing (Stripe payment confirmation)
- ❌ Email verification
- ❌ Password reset flows
- ❌ Concurrent operations
- ❌ Rate limiting
- ❌ Token refresh on expiration

## CI/CD Integration

To run these tests in a CI/CD pipeline:

1. **Set up test database**: Use a separate test database
2. **Configure secrets**: Store API URLs and credentials securely
3. **Run backend**: Start backend server in CI environment
4. **Execute tests**: Run with appropriate environment variables
5. **Cleanup**: Remove test data after test completion

Example GitHub Actions workflow:

```yaml
name: Integration Tests

on: [push, pull_request]

jobs:
  integration-tests:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:14
        env:
          POSTGRES_PASSWORD: postgres
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.25'

      - name: Start Backend
        run: |
          # Clone and setup backend
          # Run migrations
          # Start server
          python manage.py runserver &
          sleep 5

      - name: Run Integration Tests
        env:
          BASE_HOSTNAME: http://localhost:8000
          RUN_INTEGRATION_TESTS: true
        run: |
          go test -v -run TestE2EFullFlow
```

## Best Practices

1. **Isolation**: Each test run creates new users to avoid conflicts
2. **Cleanup**: Tests should be idempotent and clean up after themselves
3. **Environment**: Always use dedicated test environments, never production
4. **Data**: Use test data that's easily identifiable (e.g., `test_user_*` prefix)
5. **Assertions**: Verify both success responses and data integrity
6. **Documentation**: Keep this document updated as tests evolve

## Troubleshooting

### View API Logs

If tests fail, check the backend logs for detailed error messages:

```bash
# Backend logs (Django)
tail -f /path/to/backend/logs/django.log

# Or if using docker-compose
docker-compose logs -f backend
```

### Inspect Database

To verify test data was created correctly:

```bash
# Connect to local database
psql butler_coffee_dev

# Check test users
SELECT id, username, email, created_on FROM core_user WHERE username LIKE 'test_user_%';

# Check test orders
SELECT o.id, o.tier, o.status, u.username
FROM core_order o
JOIN core_user u ON o.user_id = u.id
WHERE u.username LIKE 'test_user_%';
```

### Manual Testing

Before running automated tests, verify endpoints manually:

```bash
# Register a user
curl -X POST http://localhost:8000/api/core/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username":"manual_test","email":"test@test.com","password":"test123"}'

# Login
curl -X POST http://localhost:8000/api/core/v1/users/token \
  -H "Content-Type: application/json" \
  -d '{"username":"manual_test","password":"test123"}'
```

## Contributing

When adding new integration tests:

1. Follow existing test naming conventions: `testXxx()`
2. Add descriptive logging with `✓` checkmarks for success
3. Validate all response fields, not just HTTP status
4. Update this README with new test documentation
5. Consider test execution time and optimize where possible
