# Snippetbox Project Summary

A web application for sharing code/text snippets, built following the **"Let's Go"** book by Alex Edwards.

---

## 📂 Project Structure

```
snippetbox/
├── cmd/
│   └── web/
│       ├── main.go         # Entry point, DI setup, TLS config
│       ├── handlers.go     # HTTP handlers
│       ├── routes.go       # Routing & middleware chains
│       ├── helpers.go      # Helper functions
│       ├── middleware.go   # Middleware functions
│       └── templates.go    # Template caching & functions
├── internal/
│   ├── models/
│   │   ├── snippets.go     # Snippet CRUD operations
│   │   ├── users.go        # User CRUD & authentication
│   │   └── errors.go       # Custom errors
│   └── validator/
│       └── validator.go    # Form validation
├── tls/
│   ├── cert.pem            # TLS certificate
│   └── key.pem             # TLS private key
├── ui/
│   ├── html/
│   │   ├── base.tmpl       # Base template
│   │   ├── pages/          # Page templates (home, view, create, signup)
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
    sessionManager *scs.SessionManager
}
```

All handlers are methods on `*application`, providing access to shared dependencies without global variables.

---

### 2. Middleware Pattern

**File:** `cmd/web/middleware.go`

| Middleware     | Purpose                                            |
| -------------- | -------------------------------------------------- |
| `commonHeader` | Sets security headers (CSP, X-Frame-Options, etc.) |
| `logRequest`   | Logs incoming requests using structured logging    |
| `recoverPanic` | Graceful panic recovery to prevent server crashes  |

**File:** `cmd/web/routes.go`

```go
// Global middleware chain
standard := alice.New(app.recoverPanic, app.logRequest, commonHeader)

// Dynamic routes middleware (with sessions)
dynamic := alice.New(app.sessionManager.LoadAndSave)
```

---

### 3. Helper Functions

**File:** `cmd/web/helpers.go`

| Function         | Purpose                                                   |
| ---------------- | --------------------------------------------------------- |
| `serverError`    | Logs 500 errors with stack trace                          |
| `clientError`    | Returns client-facing HTTP errors                         |
| `render`         | Template rendering with buffer (prevents partial writes)  |
| `newTemplates`   | Creates template data with common fields + flash messages |
| `decodePostForm` | Generic form decoding with error handling                 |

---

### 4. Template System

**File:** `cmd/web/templates.go`

- **Template caching** at startup for performance
- **Custom template functions** (e.g., `humanDate`)
- **Template inheritance** using base + pages + partials

```go
type TemplatesData struct {
    CurrentYear int
    Snippet     models.Snippet
    Snippets    []models.Snippet
    Form        any
    Flash       string  // For flash messages
}

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
mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))              // Exact match
mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(...))     // Path parameters
mux.Handle("GET /snippet/create", dynamic.ThenFunc(...))        // GET for form
mux.Handle("POST /snippet/create", dynamic.ThenFunc(...))       // POST for submission
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

### 11. Session Management

**Package:** `github.com/alexedwards/scs/v2` with `mysqlstore`

**File:** `cmd/web/main.go`

```go
sessionManager := scs.New()
sessionManager.Store = mysqlstore.New(db)  // Store sessions in MySQL
sessionManager.Lifetime = 12 * time.Hour   // Session expires after 12 hours
```

**How sessions work:**

1. `LoadAndSave` middleware loads session from MySQL

2. Session data attached to request context

3. Handler can read/write session data

4. Middleware saves session back to MySQL

---

### 12. Flash Messages

**File:** `cmd/web/handlers.go`

```go
// After creating a snippet, set a flash message
app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
```

**File:** `cmd/web/helpers.go`

```go
func (app *application) newTemplates(r *http.Request) TemplatesData {
    return TemplatesData{
        CurrentYear: time.Now().Year(),
        Flash:       app.sessionManager.PopString(r.Context(), "flash"),
    }
}
```

**How flash messages work:**

1. `Put()` stores message in session
2. `PopString()` retrieves AND removes message
3. Message only shows once (one-time notifications)

---

### 13. HTTPS/TLS Support

**File:** `cmd/web/main.go`

```go
tlsConfig := &tls.Config{
    CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
}

srv := &http.Server{
    Addr:         *addr,
    Handler:      app.routes(),
    ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
    TLSConfig:    tlsConfig,
    IdleTimeout:  60 * time.Second,
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 10 * time.Second,
}

err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
```

**Server timeouts:**
| Timeout | Value | Purpose |
|---------|-------|---------|
| `IdleTimeout` | 60s | Close idle keep-alive connections |
| `ReadTimeout` | 5s | Prevent slow client attacks |
| `WriteTimeout` | 10s | Prevent slow response attacks |

**TLS Curve Preferences:**

- `X25519` - Modern, fast elliptic curve (preferred)
- `CurveP256` - NIST P-256 curve (fallback)

---

### 14. User Authentication

**Files:** `internal/models/users.go`, `cmd/web/handlers.go`

#### User Model

```go
type User struct {
    ID             int
    Name           string
    Email          string
    HashedPassword []byte
    Created        time.Time
}

type UserModel struct {
    DB *sql.DB
}
```

#### UserModel Methods

| Method             | Purpose                                     |
| ------------------ | ------------------------------------------- |
| `Insert()`         | Create new user with bcrypt-hashed password |
| `Authenticate()`   | Verify email/password, return user ID       |
| `Exists()`         | Check if user exists by ID                  |
| `Get()`            | Retrieve user by ID                         |
| `Update()`         | Modify user's name and email                |
| `UpdatePassword()` | Change password (verifies current first)    |
| `Delete()`         | Remove user from database                   |
| `List()`           | Retrieve all users (admin functionality)    |

#### Password Security with bcrypt

```go
// Hashing password on signup (cost factor = 12)
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)

// Verifying password on login
err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
```

**Why bcrypt?**

- Automatically handles salting
- Configurable cost factor (work factor)
- Designed to be slow (prevents brute-force attacks)
- Industry standard for password hashing

#### Authentication Errors

**File:** `internal/models/errors.go`

```go
var (
    ErrNoRecord           = errors.New("models: no matching record found")
    ErrInvalidCredentials = errors.New("models: invalid credentials")
    ErrDuplicateEmail     = errors.New("models: duplicate email")
)
```

#### User Signup Flow

**File:** `cmd/web/handlers.go`

```go
type userSignupForm struct {
    Name                string `form:"name"`
    Email               string `form:"email"`
    Password            string `form:"password"`
    validator.Validator `form:"-"`
}
```

**Validation rules:**

- Name: Required (not blank)
- Email: Required, valid email format (regex)
- Password: Required, minimum 8 characters

**Signup process:**

1. Parse form data using `decodePostForm()`
2. Validate all fields with `validator`
3. Hash password with bcrypt
4. Insert user into database
5. Handle duplicate email errors gracefully
6. Set flash message and redirect to login

#### User Login Flow

**File:** `cmd/web/handlers.go`

```go
type userLoginForm struct {
    Email               string `form:"email"`
    Password            string `form:"password"`
    validator.Validator `form:"-"`
}
```

**Login process:**

1. Parse form data using `decodePostForm()`
2. Validate email format and password presence
3. Call `Authenticate()` to verify credentials with bcrypt
4. On failure: add non-field error (security: no hint which field is wrong)
5. On success: `RenewToken()` to prevent session fixation
6. Store `authenticatedUserID` in session
7. Redirect to `/snippet/create`

**Logout process:**

1. `RenewToken()` for security
2. `Remove("authenticatedUserID")` from session
3. Flash message "You've been logged out successfully!"
4. Redirect to home

#### Authentication Helpers

**File:** `cmd/web/helpers.go`

```go
// isAuthenticated checks if user is logged in
func (app *application) isAuthenticated(r *http.Request) bool {
    return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
}
```

#### Dynamic Navigation

**File:** `cmd/web/templates.go`

```go
type TemplatesData struct {
    // ... other fields
    IsAuthenticated bool  // Set automatically by newTemplates()
}
```

**File:** `ui/html/partials/nav.tmpl`

| Authenticated User | Unauthenticated User |
| ------------------ | -------------------- |
| Home               | Home                 |
| Create snippet     | Signup               |
| Logout (button)    | Login                |

---

### 15. Authorization (Route Protection)

**File:** `cmd/web/middleware.go`

```go
// requireAuthentication redirects unauthenticated users to login
func (app *application) requireAuthentication(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !app.isAuthenticated(r) {
            http.Redirect(w, r, "/user/login", http.StatusSeeOther)
            return
        }
        // Prevent caching of protected pages
        w.Header().Add("Cache-Control", "no-store")
        next.ServeHTTP(w, r)
    })
}
```

**File:** `cmd/web/routes.go`

```go
// Unprotected routes (available to all)
dynamic := alice.New(app.sessionManager.LoadAndSave)

// Protected routes (require login)
protected := dynamic.Append(app.requireAuthentication)
mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))
```

**Protected Routes:**
| Route | Method | Handler |
|-------|--------|--------|
| `/snippet/create` | GET | Show create form |
| `/snippet/create` | POST | Save new snippet |
| `/user/logout` | POST | Log out user |

#### Extended Validation Helpers

**File:** `internal/validator/validator.go`

| Function             | Purpose                                   |
| -------------------- | ----------------------------------------- |
| `MinChars()`         | Minimum character count validation        |
| `Matches()`          | Regex pattern matching (for email)        |
| `EmailRX`            | Compiled regex for email validation       |
| `AddNonFieldError()` | Add form-level error (not field-specific) |

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
| Sessions (MySQL Store)                        | ✅ Done    |
| Flash Messages                                | ✅ Done    |
| HTTPS/TLS                                     | ✅ Done    |
| Server Timeouts                               | ✅ Done    |
| User Authentication (Signup)                  | ✅ Done    |
| User Authentication (Login/Logout)            | ✅ Done    |
| Authorization (Route Protection)              | ✅ Done    |
| Dynamic Navigation                            | ✅ Done    |
| CSRF Protection                               | ⬜ Not yet |
| Testing                                       | ⬜ Not yet |

---

## 📚 Key Packages Used

| Package                                 | Purpose                   |
| --------------------------------------- | ------------------------- |
| `net/http`                              | HTTP server & routing     |
| `database/sql`                          | Database interface        |
| `github.com/go-sql-driver/mysql`        | MySQL driver              |
| `html/template`                         | HTML templating           |
| `log/slog`                              | Structured logging        |
| `github.com/justinas/alice`             | Middleware chaining       |
| `github.com/go-playground/form/v4`      | Form decoding             |
| `github.com/alexedwards/scs/v2`         | Session management        |
| `github.com/alexedwards/scs/mysqlstore` | MySQL session store       |
| `crypto/tls`                            | TLS configuration         |
| `golang.org/x/crypto/bcrypt`            | Password hashing (bcrypt) |

---

## 🔗 Quick Reference

- **Start server (HTTPS):** `go run ./cmd/web`
- **With custom port:** `go run ./cmd/web -addr=":8080"`
- **Access app:** `https://localhost:4000`

---

## 📝 Related Notes

- [Go Embedding](./go-embedding.md) - Struct embedding & composition

- [Go Panic/Recover](./go-panic-recover.md) - Panic handling
