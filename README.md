# ⏰ Clocked — Attendance & Tasks App

A lightweight attendance and task tracking app.

- 💚 Frontend: Nuxt 3 (SPA, Tailwind)
- ⚡️ Backend: Go + Fiber
- 🐘 Database: PostgreSQL (via GORM)
- 🐳 Production: Docker Compose (nginx + Fiber + Postgres)

---

## ✨ Features

- ✅ Tasks: create, view, and mark completed (with completion date)
- 🕒 Attendances: clock in/out, associate a task or create on clock-out
- 📅 Dashboard: Indonesian holidays + tasks for the selected day
- 🔐 JWT auth (login with username or email)

---

## 🧭 Architecture

```
┌───────────┐     /api/*      ┌──────────────┐
│  Nuxt SPA │  ───────────▶   │   Fiber API  │  ─▶  PostgreSQL
│ (nginx)   │  ◀─ static ◀──  │   (Go v2)    │
└───────────┘                 └──────────────┘
```

- The SPA calls `/api/...` (nginx proxies to the backend).
- Backend serves JSON and uses GORM for DB access.

---

## 📁 Project Structure

```
Clocked/
├── backend/             # Fiber app, models, handlers, router
│   ├── handler/         # Auth, dashboard, tasks, attendance
│   ├── model/           # User, Product, Task, Attendance
│   ├── router/          # Route groups & middleware
│   └── database/        # GORM connection & automigrations
├── frontend/            # Nuxt 3 SPA
│   ├── app/pages/       # index, login, tasks, attendances
│   └── app/composables/ # useAuth, useHolidays, useTasks, useAttendance
├── docker-compose.yml   # Frontend + Backend + DB
└── README.md
```

---

## 🚀 Quick Start (Docker Compose)

Make sure Docker is installed on your VPS.

```
# From the project root
docker compose build
JWT_SECRET=super_secret_token docker compose up -d

# Visit the app
http://YOUR_SERVER_IP/
```

Default services

- Frontend (nginx): port 80, serves SPA and proxies `/api` to backend
- Backend (Fiber): internal port 8000 (not exposed by default)
- Postgres: internal, persisted in `db_data` volume

---

## 🔧 Environment Variables

Backend (Fiber)

- `DB_HOST` (default: `db` in docker)
- `DB_PORT` (default: `5432`)
- `DB_USER` (default: `clocked`)
- `DB_PASSWORD` (default: `clocked`)
- `DB_NAME` (default: `clocked`)
- `SECRET` (JWT signing secret)

Frontend (Nuxt SPA)

- `NUXT_PUBLIC_API_BASE`: public API base used by the SPA (default `/api`)
  - In docker, nginx proxies `/api` → backend

---

## 🛠️ Local Development

Backend

```
# Terminal A
cd backend
export DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=clocked SECRET=dev_secret
# Run postgres however you prefer (e.g., local or docker)
go run ./...
```

Frontend

```
# Terminal B
cd frontend
# API base for local dev, e.g., to call localhost:8000
echo 'NUXT_PUBLIC_API_BASE=http://localhost:8000/api' > .env
npm install
npm run dev
```

Then open http://localhost:3000/

---

## 🔗 API Overview (selected)

Auth

- `POST /api/auth/login` — body: `{ identity, password }` → returns JWT

Dashboard

- `GET /api/dashboard` — Holidays map `{ YYYY-MM-DD: { summary } }`

Tasks (auth)

- `GET /api/tasks` — list user tasks
- `POST /api/tasks` — create `{ name, subject?, description?, how_its_done? }`
- `PATCH /api/tasks/:id` — `{ completed?: boolean, completed_at?: ISO }`
- `GET /api/tasks/day?date=YYYY-MM-DD` — tasks created or completed that day

Attendances (auth)

- `POST /api/attendance/clockin` — `{ clock_in?: ISO }` (defaults to now)
- `POST /api/attendance/clockout` — `{ clock_out?: ISO, task_id?: number, new_task?: {...}, mark_completed?: boolean }`
- `GET /api/attendance` — history

Users (auth)

- `GET /api/user/` — current user
- `PATCH /api/user/:id` — update names
- `DELETE /api/user/:id` — delete user
- `POST /api/user/:id/avatar` — upload avatar
- `GET /api/user/:id/avatar` — get avatar image

---

## 🧩 Implementation Notes

- SPA mode (`ssr: false`) and Nuxt public runtime config ensure the app can run behind a reverse proxy cleanly.
- GORM automigrates `User`, `Product`, `Task`, `Attendance` on startup.
- Attendance clock-out can link an existing task or create a new one and optionally mark it completed at the clock-out time.

---

## 🧪 Smoke Test

- Register a user or insert one in DB.
- Log in via `/login` page.
- Create a task on `/tasks`.
- Clock in/out on `/attendances`, choose a task (or create one), mark completed.
- On the home page, select a date to see tasks for that day.

---

Made with ❤️ using Nuxt + Fiber.
