# IFTP Class Management System

A backend API and web portal for managing student enrollment and class scheduling for IFTP improv classes. Built with Go (standard `net/http`), PostgreSQL, and Bootstrap.

---

## Tech Stack

- **Backend:** Go 1.24, `net/http`
- **Database:** PostgreSQL (`pgx/v5`)
- **Frontend:** Bootstrap, HTML, CSS, JavaScript

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
BASE_URL=http://localhost:8080
DATABASE_URL=postgres://user:password@localhost:5432/iftp
```

### Run
```bash
go run main.go
```

---

## API Reference

### Classes
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/classes/all` | List all classes |
| GET | `/classes` | List classes by month |
| GET | `/classes/{student_id}` | List classes by month for a student |
| POST | `/classes` | Create a class |
| PATCH | `/classes/{class_id}` | Update a class |

### Roster
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/roster/{class_id}` | Get roster for a class |
| GET | `/roster/enrollment/{student_id}` | Get enrollments for a student |

### Enrollment
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/enrollment_request` | Submit an enrollment request |

### Calendar
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/calendarEvents` | Get all calendar events |
| GET | `/calendarEvents/{student_id}` | Get calendar events for a student |

### Students
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/students` | List all students |
| GET | `/students/enrollment` | List students with enrollment data |
| POST | `/students` | Add a student |
| PATCH | `/students/{student_id}` | Update a student |

---

## Notes

Currently undergoing refactor from Gin to standard `net/http`.