# Snippetbox Project Summary

A web application for sharing code/text snippets, built following the **"Let's Go"** book by Alex Edwards.

---

## 📂 Project Structure

```
snippetbox/
├── cmd/
│   └── web/
│       ├── main.go         # Entry point, DI setup
│       ├── handlers.go     # HTTP handlers
│       ├── routes.go       # Routing configuration
│       ├── helpers.go      # Helper functions
│       ├── middleware.go   # Middleware functions
│       └── templates.go    # Template caching & functions
├── internal/
│   ├── models/
│   │   ├── snippets.go     # Snippet CRUD operations
│   │   └── errors.go       # Custom errors
│   └── validator/
│       └── validator.go    # Form validation
├── ui/
│   ├── html/
│   │   ├── base.tmpl       # Base template
│   │   ├── pages/          # Page templates
│   │   └── partials/       # Partial templates
│   └── static/             # CSS, JS, images
└── notes/                  # Learning notes
```

---

## 🛠️ Skills & Techniques Learned

### 1. Dependency Injection (DI)

**File:** `cmd/web/main.go`

```go
type application struct {
    logger         *slog.Logger
    snippets       *models.SnippetModel
    templateCache  map[string]*template.Template
    formDecoder    *form.Decoder
}
```

All handlers are methods on `*application`, providing access to shared dependencies without global variables.

---

### 2. Middleware Pattern

**File:** `cmd/web/middleware.go`

| Middleware     | Purpose                                                                    |
| -------------- | -------------------------------------------------------------------------- |
| `commonHeader` | Sets security headers (CSP, X-Frame-Options, X-Content-Type-Options, etc.) |
| `logRequest`   | Logs incoming requests using structured logging                            |
| `recoverPanic` | Graceful panic recovery to prevent server crashes                          |

**Middleware chaining with Alice:**

```go
standard := alice.New(app.recoverPanic, app.logRequest, commonHeader)
return standard.Then(mux)
```

---

### 3. Helper Functions

**File:** `cmd/web/helpers.go`

| Function         | Purpose                                                           |
| ---------------- | ----------------------------------------------------------------- |
| `serverError`    | Logs 500 errors with stack trace                                  |
| `clientError`    | Returns client-facing HTTP errors                                 |
| `render`         | Template rendering with buffer (prevents partial writes on error) |
| `newTemplates`   | Creates template data with common fields                          |
| `decodePostForm` | Generic form decoding with error handling                         |

---

### 4. Template System

**File:** `cmd/web/templates.go`

- **Template caching** at startup for performance
- **Custom template functions** (e.g., `humanDate`)
- **Template inheritance** using base + pages + partials

```go
var functions = template.FuncMap{
    "humanDate": humanDate,
}
```

---

### 5. Form Validation

**File:** `internal/validator/validator.go`

Reusable validator that can be embedded in form structs:

```go
type Validator struct {
    FieldErrors map[string]string
}
```

**Validation helpers:**

- `NotBlank()` – Checks for non-empty values
- `MaxChars()` – Character limit validation
- `PermittedValue()` – Generic "must be one of" validation (uses generics)

**Usage in handlers:**

```go
type snippetCreateForm struct {
    Title               string `form:"title"`
    Content             string `form:"content"`
    Expires             int    `form:"expires"`
    validator.Validator `form:"-"`  // Embedded
}

form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
```

---

### 6. Database Layer (Repository Pattern)

**File:** `internal/models/snippets.go`

```go
type SnippetModel struct {
    DB *sql.DB
}
```

| Method     | Purpose                              |
| ---------- | ------------------------------------ |
| `Insert()` | Create a new snippet                 |
| `Get()`    | Retrieve a single snippet by ID      |
| `Latest()` | Retrieve the 10 most recent snippets |

---

### 7. Error Handling

**File:** `internal/models/errors.go`

- **Sentinel errors:** `ErrNoRecord` for not-found cases
- Using `errors.Is()` for error comparison

```go
if errors.Is(err, models.ErrNoRecord) {
    http.NotFound(w, r)
}
```

---

### 8. HTTP Routing (Go 1.22+)

**File:** `cmd/web/routes.go`

```go
mux.HandleFunc("GET /{$}", app.home)              // Exact match for root
mux.HandleFunc("GET /snippet/view/{id}", ...)     // Path parameters
mux.HandleFunc("GET /snippet/create", ...)        // GET for form
mux.HandleFunc("POST /snippet/create", ...)       // POST for submission
```

---

### 9. Structured Logging (slog)

**File:** `cmd/web/main.go`

```go
logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
app.logger.Info("starting server", "addr", *addr)
app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
```

---

### 10. Configuration with CLI Flags

**File:** `cmd/web/main.go`

```go
addr := flag.String("addr", ":4000", "HTTP network address")
dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL DSN")
flag.Parse()
```

---

## 📊 Progress Checklist

| Topic                                         | Status     |
| --------------------------------------------- | ---------- |
| Basic HTTP Server                             | ✅ Done    |
| Project Structure                             | ✅ Done    |
| Configuration & CLI Flags                     | ✅ Done    |
| Structured Logging (slog)                     | ✅ Done    |
| Dependency Injection                          | ✅ Done    |
| Database (MySQL) Setup                        | ✅ Done    |
| Repository/Model Layer                        | ✅ Done    |
| HTML Templates & Caching                      | ✅ Done    |
| Middleware (Logging, Headers, Panic Recovery) | ✅ Done    |
| Form Processing & Validation                  | ✅ Done    |
| RESTful Routing                               | ✅ Done    |
| Sessions                                      | ⬜ Not yet |
| User Authentication                           | ⬜ Not yet |
| CSRF Protection                               | ⬜ Not yet |
| HTTPS/TLS                                     | ⬜ Not yet |
| Testing                                       | ⬜ Not yet |

---

## 📚 Key Packages Used

| Package                            | Purpose               |
| ---------------------------------- | --------------------- |
| `net/http`                         | HTTP server & routing |
| `database/sql`                     | Database interface    |
| `github.com/go-sql-driver/mysql`   | MySQL driver          |
| `html/template`                    | HTML templating       |
| `log/slog`                         | Structured logging    |
| `github.com/justinas/alice`        | Middleware chaining   |
| `github.com/go-playground/form/v4` | Form decoding         |

---

## 🔗 Quick Reference

- **Start server:** `go run ./cmd/web`
- **With custom port:** `go run ./cmd/web -addr=":8080"`
- **Access app:** `http://localhost:4000`
