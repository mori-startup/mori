# Overview of server.go

The **server.go** file is the entry point for the application’s backend server. It handles the initialization of various components and sets up all the HTTP routes. Below is an explanation of its key parts:

## Main Function

- **Database Initialization:**  
  The file starts by initializing the database connection using `sqlite.InitDB()`. It then defers closing the database connection to ensure that resources are released when the application exits.

- **Repository & Handler Setup:**  
  After the database is ready, the code initializes repositories with `sqlite.InitRepositories(db)`. These repositories abstract data access operations and are passed to `handlers.InitHandlers(repos)` to create a handler instance. The handlers contain the business logic for each endpoint.

- **WebSocket Server Initialization:**  
  The WebSocket server is started via `ws.StartServer(repos)`, which is later used by some of the routes to handle real-time communication.

- **HTTP Server Configuration:**  
  A new `http.Server` is created with:
  - An address (`:8081`), meaning the server listens on port 8081.
  - A handler returned by the `setRoutes` function, which defines all routes.

- **Server Start:**  
  Finally, the server prints a startup message and calls `ListenAndServe()`, which starts the HTTP server. Any error during startup is printed to the console.

## setRoutes Function

The `setRoutes` function configures all the HTTP endpoints (routes) used by the application:

- **Static File Serving:**  
  Routes starting with `/imageUpload/` serve static files (e.g., images) using a file server with custom headers configured via `utils.ConfigFSHeader`.

- **Authentication & Verification Routes:**  
  Endpoints such as `/register`, `/signin`, `/logout`, `/captcha`, `/verified`, `/sessionActive`, `/request-password-reset`, and `/reset-password` handle user authentication, registration, and account verification.

- **LLM and Chat Routes:**  
  The endpoint `/llmConvo` is protected (wrapped by an authentication middleware) and handles interactions with the chatbot (LLM). Similarly, several routes manage user data, group functionalities, notifications, chat messages, and conversations.

- **WebSocket Route:**  
  The `/ws` endpoint sets up a WebSocket connection for real-time communication.

- **File Upload Endpoints:**  
  Routes for file upload, listing, and deletion are defined under `/api/upload`, `/api/files`, and `/api/files/`.

The modular design of the route setup ensures that different aspects of the application (authentication, chat, groups, etc.) are organized and easily maintainable.

## Summary

- **Initialization:**  
  The server initializes the database, repositories, handlers, and a WebSocket server.
- **Route Setup:**  
  All HTTP routes are defined in a centralized function (`setRoutes`), ensuring a clear separation of concerns.
- **Execution:**  
  The server starts listening on port 8081 and handles incoming requests based on the defined routes.

This architecture facilitates scalability and modularity, allowing individual components (like authentication or chat functionalities) to be developed and tested separately.
