# IFTP Class Management System
 
A backend API and web portal for managing student enrollment and class scheduling for IFTP improv classes. Built with Go (standard `net/http`), PostgreSQL, and Bootstrap.
 
---
 
## Tech Stack
 
- **Backend:** Go 1.24, `net/http`
- **Database:** PostgreSQL (`pgx/v5`)
- **Frontend:** Bootstrap, HTML, CSS, JavaScript
- **Auth:** Google OAuth via `goth` / `gothic`, cookie-based sessions
- **Hosting:** Railway
---
 
## Getting Started
 
### Prerequisites
- Go 1.24+
- PostgreSQL
### Installation
```bash
git clone https://github.com/yourhandle/IFTP.git
cd IFTP
go mod download
```
 
### Configuration
 
Create a `.env` file in the project root:
```env
PORT=8080
CONN_STR=postgres://user:password@localhost:5432/iftp
GOOGLE_KEY=your_google_oauth_key
GOOGLE_SECRET=your_google_oauth_secret
GOOGLE_SCOPES=openid email profile
SESSION_SECRET=your_session_secret
API_KEY=your_api_key
```
 
### Run
```bash
go run main.go
```
 
---
 
## Project Structure
 
```
IFTP/
├── main.go           # Server setup, routing
├── utils.go          # Auth, middleware, shared handlers
├── db.go             # Database connection pool struct
├── timeutils.go      # Date/weekday helpers
├── class/
│   ├── handlers.go   # Class and calendar HTTP handlers
│   └── database.go   # Class DB queries
├── roster/
│   ├── handlers.go   # Roster and enrollment HTTP handlers
│   └── database.go   # Roster DB queries
├── students/
│   ├── handlers.go   # Student HTTP handlers
│   └── database.go   # Student DB queries
├── makeup/
│   ├── handlers.go   # Makeup request HTTP handlers
│   └── database.go   # Makeup request DB queries
├── static/           # Static assets
└── templates/
    └── index.html    # Main SPA template
```
 
---
 
## Database Schema
 
### Core Tables
- **students** — student records with `makeup_credits` counter
- **classes** — recurring class definitions (name, teacher, day, time, capacity)
- **class_schedule** — individual session dates per class per month
- **roster** — student enrollment per session date, with status (`Enrolled` / `Away`)
### Request Tables
- **enrollment_requests** — student requests to join a new class, pending admin approval
- **schedule_approvals** — auto-generated monthly schedules pending admin review
- **schedule_approval_dates** — proposed dates for each schedule approval
- **makeup_requests** — student requests to be marked absent for a session
- **makeup_request_dates** — individual missed dates per makeup request
- **makeup_redemptions** — student requests to use a makeup credit for a specific session
---
 
## Authentication
 
- Google OAuth via `/auth/google` → `/auth/google/callback`
- Sessions stored in signed cookies via `gorilla/sessions`
- All routes under `/` require authentication via `RequireAuth` middleware
- API key auth supported via `X-API-Key` header for server-to-server or testing
---
 
## API Reference
 
### Auth
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/auth/google` | Initiate Google OAuth login |
| GET | `/auth/google/callback` | Google OAuth callback |
| GET | `/auth/user` | Get current session user |
| GET | `/logout` | Clear session and log out |
| GET | `/health` | Health check + DB ping |
 
### Classes
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/classes/all` | List all classes across all months |
| GET | `/classes` | List classes filtered by month |
| GET | `/classes/{student_id}` | List enrolled classes for a student by month |
| POST | `/classes` | Create a new class |
| PATCH | `/classes/{class_id}` | Update a class |
 
### Schedule Approvals
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/classes/schedule_approval/generate` | Trigger auto-generation of next month's schedules |
| GET | `/classes/schedule_approval` | Get pending schedule approvals |
| PATCH | `/classes/schedule_approval/confirm/{approval_id}` | Approve or reject a schedule |
 
### Roster
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/roster/{class_id}` | Get roster for a class session |
| GET | `/roster/enrollment/{student_id}` | Get all enrollments for a student |
| POST | `/roster/{class_id}/enroll` | Directly enroll a student in a class |
 
### Enrollment Requests
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/enrollment_requests` | Submit an enrollment request |
| GET | `/enrollment_requests` | Get all pending enrollment requests |
| PATCH | `/enrollment_requests/{request_id}` | Approve or deny an enrollment request |
 
### Makeup Requests
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/makeup_requests` | Submit a makeup request for missed sessions |
| GET | `/makeup_requests` | Get all pending makeup requests |
| PATCH | `/makeup_requests/{request_id}` | Approve or deny a makeup request |
 
### Makeup Redemptions
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/makeup_redemptions` | Request to use a makeup credit for a session |
| GET | `/makeup_redemptions` | Get all pending redemption requests |
| PATCH | `/makeup_redemptions/{request_id}` | Approve or deny a redemption request |
 
### Calendar
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/calendarEvents` | Get all scheduled calendar events |
| GET | `/calendarEvents/{student_id}` | Get calendar events for a specific student |
 
### Students
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/students` | List all active students |
| GET | `/students/enrollment` | List students with their enrolled classes |
| POST | `/students` | Add a new student |
| PATCH | `/students/{student_id}` | Update a student |
 
---
 
## Key Flows
 
### Student Enrollment
1. Student browses available classes for the month
2. Submits an enrollment request with a reason
3. Admin reviews and approves/denies
4. On approval, student is enrolled in all sessions for that month
### Monthly Schedule
1. Cron or manual trigger hits `POST /classes/schedule_approval/generate`
2. System calculates session dates for all active classes for next month
3. Admin reviews proposed dates, can edit on a calendar view
4. On approval, dates are written to `class_schedule`
### Makeup Request
1. Student selects upcoming sessions they will miss from their enrolled class
2. Submits a makeup request with optional reason
3. Admin approves/denies
4. On approval: roster status flipped to `Away` for each missed date, student receives one makeup credit per missed session
### Makeup Redemption *(in progress)*
1. Student sees their available makeup credits
2. Selects an available class and session date to attend
3. Submits a redemption request with optional note
4. Admin approves/denies
5. On approval: student enrolled in that session, credit decremented by 1
---
 
## Notes
 
- Admin and student views are currently served from the same SPA with a toggle button — role-based view separation and admin route protection is planned
- Roster `Away` status is used for approved makeup absences to preserve enrollment history