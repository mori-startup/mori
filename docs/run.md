# Running the Project in Test Phase

This document explains the steps needed to run the project in its test phase. The project is divided into three main parts: the backend server, the frontend application, and the LLM logic server.

---

## Step 1: Run the Backend Server

1. Open a terminal.

2. Navigate to the **/backend** folder:
   cd /path/to/project/backend

    Start the backend server by running:

    go run .

    This command compiles and launches the Go backend server.

Step 2: Run the Frontend Application

    Open another terminal window.

    Navigate to the /frontend folder:

cd /path/to/project/frontend

Start the frontend server by running:

    npm run serve

    This command starts the Vue.js development server. Typically, the application will be available at http://localhost:8080.

Step 3: Run the LLM Server

    Open a third terminal window.

    Navigate to the /backend/logicllm folder:

cd /path/to/project/backend/logicllm

Start the LLM server using Uvicorn:

uvicorn server:app --host 127.0.0.1 --port 8000

This command runs the LLM logic server on http://127.0.0.1:8000.