# GenAI Platform - Local Setup Guide

This guide provides comprehensive instructions to set up and run the GenAI Platform locally on your machine. It covers the backend (Go), frontend (React), and database (PostgreSQL) setup.

## 1. Prerequisites

Before you begin, ensure you have the following installed on your system:

*   **Go**: Version 1.18 or higher. You can download it from [https://golang.org/dl/](https://golang.org/dl/).
*   **Node.js and npm**: Node.js version 14 or higher, and npm version 6 or higher. You can download them from [https://nodejs.org/en/download/](https://nodejs.org/en/download/).
*   **PostgreSQL**: Version 12 or higher. You can download it from [https://www.postgresql.org/download/](https://www.postgresql.org/download/).
*   **Git**: For cloning the repository. Download from [https://git-scm.com/downloads](https://git-scm.com/downloads).
*   **jq**: A lightweight and flexible command-line JSON processor. Install using your system's package manager (e.g., `sudo apt-get install jq` on Ubuntu, `brew install jq` on macOS).
*   **pandoc**: A universal document converter. Install using your system's package manager (e.g., `sudo apt-get install pandoc` on Ubuntu, `brew install pandoc` on macOS).

## 2. Clone the Repository

First, clone the GenAI Platform repository to your local machine:

```bash
git clone <repository_url> # Replace with the actual repository URL
cd genai-platform
```

## 3. Backend Setup (Go)

### 3.1. Database Configuration

The backend uses PostgreSQL. You need to set up a database and configure the connection.

1.  **Start PostgreSQL Service**: Ensure your PostgreSQL service is running. The command varies based on your operating system. For Ubuntu, it's typically:
    ```bash
sudo systemctl start postgresql
sudo systemctl enable postgresql
    ```

2.  **Create Database User and Database**: By default, the application expects a user `postgres` with password `password` and a database named `genai_platform`. You can create them using the following commands:
    ```bash
sudo -u postgres psql -c "ALTER USER postgres PASSWORD 'password';"
sudo -u postgres createdb genai_platform
    ```
    If you wish to use different credentials or database name, you will need to update the `DATABASE_URL` environment variable accordingly.

3.  **Database Migrations**: The application will automatically run database migrations on startup. No manual migration steps are required.

### 3.2. Build and Run the Backend

Navigate to the `genai-platform` directory (if you are not already there) and build the Go application:

```bash
cd /path/to/your/genai-platform # Replace with your actual path
go mod tidy
go build -o bin/server_local ./cmd/server
```

Now, you can run the backend server:

```bash
./bin/server_local
```

The backend server will start on `http://localhost:8080` by default. You can change the port by setting the `PORT` environment variable (e.g., `PORT=9000 ./bin/server_local`).

## 4. Frontend Setup (React)

### 4.1. Install Dependencies

Navigate to the `genai-frontend` directory and install the Node.js dependencies:

```bash
cd /path/to/your/genai-frontend # Replace with your actual path
npm install
```

### 4.2. Configure API Endpoint

The frontend needs to know where your backend API is located. Open the `genai-frontend/.env` file (create it if it doesn't exist) and add the following line:

```
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

### 4.3. Run the Frontend

From the `genai-frontend` directory, start the development server:

```bash
npm run dev
```

The frontend application will typically open in your browser at `http://localhost:5173`.

## 5. AI Service Setup (Python)

The AI services (LLM integration, RAG, etc.) are implemented in Python and called by the Go backend. These services require specific Python libraries.

1.  **Install Python Dependencies**: Navigate to the `genai-platform` directory and install the required Python packages:
    ```bash
cd /path/to/your/genai-platform # Replace with your actual path
pip3 install langchain langchain-openai langchain-google-genai faiss-cpu pypdf2 python-docx
    ```

2.  **API Keys**: For full functionality with real LLMs (OpenAI GPT, Google Gemini), you will need to set up your API keys as environment variables. The Python scripts will look for these. For example:
    ```bash
export OPENAI_API_KEY="your_openai_api_key"
export GOOGLE_API_KEY="your_google_api_key"
    ```
    If these are not set, the AI services will use mock implementations.

## 6. Accessing the Platform

Once both the backend and frontend servers are running, open your web browser and navigate to the frontend URL (e.g., `http://localhost:5173`). You can then register a new user and start using the GenAI Platform.

## 7. Troubleshooting

*   **Port already in use**: If you encounter an error like `listen tcp 0.0.0.0:8080: bind: address already in use`, it means another process is using the required port. You can find and kill the process using:
    ```bashrun
sudo lsof -t -i:8080
sudo kill -9 <PID>
    ```
    (Replace `8080` with the port number and `<PID>` with the process ID).
*   **Database connection issues**: Double-check your PostgreSQL service status, user credentials, and database name. Ensure the `DATABASE_URL` environment variable (if set) is correct.
*   **Frontend not connecting to backend**: Verify that the `VITE_API_BASE_URL` in your `.env` file points to the correct backend address and port.
*   **AI service errors**: Ensure all Python dependencies are installed and that your API keys (if using real LLMs) are correctly set as environment variables.

## 8. Development Notes

*   **Backend (Go)**: Changes to Go files require recompilation (`go build`) and restarting the server.
*   **Frontend (React)**: Changes to React files will typically trigger a hot reload in the development server.
*   **Database Schema**: The Go application handles database migrations automatically on startup. If you make changes to the database schema in the Go models, ensure your migrations are correctly defined.

---



