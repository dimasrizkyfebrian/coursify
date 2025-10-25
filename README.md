# Coursify - Learning Management System 📚

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25-blue.svg)](https://golang.org/)
[![Vue.js](https://img.shields.io/badge/Vue.js-3.x-brightgreen.svg)](https://vuejs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-Used-blue.svg)](https://www.typescriptlang.org/)
[![Docker](https://img.shields.io/badge/Docker-Used-blueviolet.svg)](https://www.docker.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-darkblue.svg)](https://www.postgresql.org/)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-Used-38B2AC.svg)](https://tailwindcss.com/)
[![Vite](https://img.shields.io/badge/Vite-4.x-blue.svg)](https://vitejs.dev/)

Coursify is a full-stack web application serving as a **Learning Management System (LMS)**, similar to platforms like Google Classroom. This application allows instructors to create and manage courses along with their materials, while students can enroll and access learning content.

## ✨ Key Features

- **Authentication & Authorization**: JWT-based login system with distinct user roles (Admin, Instructor, Student).
- **User Management (Admin)**:
  - View user statistics (total, active, pending).
  - Manage users (approve/reject pending status, view list, delete).
- **Course Management (Instructor)**:
  - Create new courses by selecting from provided cover images.
  - Update course details.
  - Delete courses.
  - View the list of students enrolled in their courses.
- **Material Management (Instructor)**:
  - Add learning materials (text, video links, PDF uploads).
  - Upload PDF files for materials.
  - Edit and delete materials.
- **Student Functionality**:
  - Enroll in available courses.
  - View the list of enrolled courses.
  - Access details and materials of enrolled courses.
- **Profile Management**: Users can upload and update their profile avatars (JPG/PNG converted to WebP).
- **Containerization**: The entire application (backend, frontend, database) runs within Docker containers for a consistent development environment.

## 🛠️ Tech Stack

- **Backend**:
  - Go (Golang)
  - Chi (Router)
  - PostgreSQL (Database)
  - `golang-migrate/migrate` (Database Migrations)
  - JWT (Authentication)
  - `swaggo/swag` (API Documentation)
  - `kolesa-team/go-webp` (Image Conversion)
- **Frontend**:
  - Vue.js 3 (Composition API)
  - TypeScript
  - Vite (Build Tool)
  - Tailwind CSS (Styling)
  - shadcn-vue (UI Components)
  - TanStack Table (Data Tables)
  - Pinia (State Management - _implied_)
  - Axios (HTTP Client)
- **DevOps**:
  - Docker
  - Docker Compose

## ⚠️ Current Status & Known Limitations

Please note that Coursify is currently under **active development**. As such, you may encounter bugs or incomplete features.

- **Storage**: File uploads (avatars, course materials) are currently stored on the **local filesystem** within the `backend/uploads` directory. Cloud storage integration (like AWS S3 or Google Cloud Storage) is planned for future development.
- **Testing**: While core functionalities are tested, comprehensive test coverage is still in progress.

## 🚀 Running the Project Locally (Using Docker)

Ensure you have **Docker Desktop** installed.

1.  **Clone the Repository**:

    ```bash
    git clone https://github.com/dimasrizkyfebrian/coursify.git
    cd coursify
    ```

2.  **Configure Environment Variables**:

    - Copy the `.env.example` file (if present) to `.env` inside the `backend/` directory.
    - Fill in the required environment variables in `backend/.env`, especially:
      - `DB_HOST=db` (Keep this for Docker)
      - `DB_PORT=5432` (Keep this for Docker)
      - `DB_USER`
      - `DB_PASSWORD`
      - `DB_NAME`
      - `DB_SSLMODE=disable`
      - `JWT_SECRET_KEY` (Generate a strong secret key)
    - Copy the `.env.example` file (if present) to `.env` inside the `frontend/` directory.
    - Fill in the required environment variables in `frontend/.env`:
      - `VITE_API_BASE_URL=http://localhost:8080/api` (Use the exposed backend port)

3.  **Run Docker Compose**:

    - Open a terminal in the project root directory (`coursify/`).
    - Run the following command to build the images and start all containers (database, backend, frontend):
    - ```bash
      docker-compose up -d --build
      ```

4.  **Run Database Migrations**:

    - Open a new terminal or use the same one.
    - Navigate to the `backend/` directory:
    - ```powershell
      cd backend
      ```

    ````
    - Run the migrations using `golang-migrate/migrate`. Replace `YOUR_PASSWORD` and `YOUR_DB_NAME` according to your `.env`:
    - ```powershell
    # Using PowerShell
    migrate -database "postgres://postgres:YOUR_PASSWORD@localhost:5433/YOUR_DB_NAME?sslmode=disable" -path ./migrations up

    # Or using Bash/Cmd
    # export DB_URL="postgres://postgres:YOUR_PASSWORD@localhost:5433/YOUR_DB_NAME?sslmode=disable"
    # migrate -database ${DB_URL} -path ./migrations up
    ````

    - _(Note: The migration command might differ based on your `golang-migrate/migrate` setup)_

5.  **Run Seeder (Optional but Recommended)**:

    - Still inside the `backend/` directory, run the seeder script:
    - ```powershell
      go run ./cmd/seeder/main.go
      ```

    ```

    - This will create default users (admin, instructor, student) and some random data.

    ```

6.  **Access the Application**:
    - **Frontend**: Open your browser and navigate to `http://localhost:5173` (as per `docker-compose.yml` port mapping).
    - **Backend API**: Available at `http://localhost:8080`.
    - **API Documentation (Swagger)**: Access `http://localhost:8080/swagger/index.html`.
    - **Database**: Can be accessed using tools like TablePlus/DBeaver at `localhost:5433`.

```

```
