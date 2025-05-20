# Explanation of `server_test.go`

This document explains the structure and purpose of the `server_test.go` file, which is used to test every function and route defined in the `server.go` file. The tests ensure that the HTTP endpoints are correctly set up and return the expected responses using a table-driven approach.

---

## Overview

The `server_test.go` file is designed to validate the routing logic of the server by:
- Simulating HTTP requests to each endpoint.
- Verifying that the correct handler functions are invoked.
- Comparing the actual response output with the expected output.
- Displaying the results in a visual, ASCII table format.

---

## Key Components

### 1. Interface Definition and Dummy Implementations

- **FullHandlerInterface**  
  An interface is defined to list all the handler methods required by the routes. This interface ensures that our dummy handler implements every function that is used by the routing logic.

- **dummyHandler**  
  This is a dummy implementation of `FullHandlerInterface`. Each method in `dummyHandler` writes a fixed string (e.g., `"register"`, `"signin"`, etc.) to the HTTP response. This allows tests to verify that the correct endpoint is called without needing to invoke the real business logic.

- **dummyWSServer**  
  A dummy implementation for the WebSocket server (or a placeholder for a `ws.Server`). It provides a dummy instance so that routes requiring a WebSocket server can function correctly during tests.

### 2. Route Setup for Testing

- **setRoutesForTest**  
  This function sets up the HTTP routes using the dummy handler and dummy WebSocket server. It mirrors the production route configuration in `server.go` but is tailored for testing purposes. By registering the routes with a simple HTTP multiplexer (`http.NewServeMux()`), it isolates the routing logic from the rest of the application.

### 3. Table-Driven Testing

- **Table-Driven Test Cases**  
  The test function `TestAllEndpointsVisual` defines a slice of test cases, where each case specifies:
  - A descriptive **name** for the endpoint test.
  - The **endpoint URL** (route path).
  - The **expected output** from the dummy handler.
  
- **Iterating over Test Cases**  
  For each test case:
  - An HTTP GET request is created and sent to the corresponding endpoint.
  - The response is captured using an HTTP recorder.
  - The actual output is compared with the expected output.
  - The results (endpoint, expected, and actual outputs) are logged in a formatted ASCII table for visual clarity.

### 4. Visual Output

- **ASCII Table Logging**  
  The test logs a header, each test result, and a footer forming an ASCII table. This table provides a clear, visual summary of:
  - Which endpoints were tested.
  - What output was expected.
  - What output was actually returned.

---

## How the Tests Work

1. **Setup:**  
   The dummy implementations of the handler and WebSocket server are created. Then, `setRoutesForTest` configures all the endpoints using these dummy implementations.

2. **Execution:**  
   A table of test cases is defined for every route defined in `server.go`. The test function iterates over each case, sends a simulated GET request, and records the response.

3. **Verification:**  
   Each endpoint's response is compared with the expected output. If they match, the test passes; if not, the test reports an error and logs the discrepancy.

4. **Visual Feedback:**  
   An ASCII table is printed in the test logs (visible when running `go test -v`), summarizing the endpoint tested, expected output, and actual output. This visual format makes it easier to spot any mismatches.

---

## Running the Tests

- **Command:**  
  To run the tests and see the detailed output, use the command:
  ```bash
  go test -v
