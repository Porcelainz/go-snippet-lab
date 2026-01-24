# Go Middleware Chaining with Alice

**Package:** `github.com/justinas/alice`

Alice is a small package that helps chain HTTP middleware together in a clean, readable way.

---

## The Problem: Without Alice

Without Alice, chaining middleware looks like this:

```go
func (app *application) routes() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("/", app.home)

    // Nested middleware - hard to read! 😵
    return app.recoverPanic(app.logRequest(commonHeader(mux)))
}
```

**Issues:**

- Reads **right-to-left** (confusing)

- Hard to add/remove middleware

- Gets messy with more middleware

- Adding per-route middleware is awkward

---

## The Solution: With Alice

```go
func (app *application) routes() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("/", app.home)

    // Clean chain - reads left-to-right! ✅
    standard := alice.New(app.recoverPanic, app.logRequest, commonHeader)
    return standard.Then(mux)
}
```

**Benefits:**

- Reads **left-to-right** (natural order)

- Easy to add/remove middleware

- Clean and declarative

- Reusable chains

---

## Key Benefits of Alice

### 1. Readable Execution Order

```go
// Execution order: recoverPanic → logRequest → commonHeader → handler
standard := alice.New(app.recoverPanic, app.logRequest, commonHeader)
```

The order you write is the order they execute. No more mental gymnastics!

---

### 2. Reusable Middleware Chains

You can create different chains for different route groups:

```go
func (app *application) routes() http.Handler {
    mux := http.NewServeMux()

    // Chain for ALL routes
    standard := alice.New(app.recoverPanic, app.logRequest, commonHeader)

    // Chain for routes that need sessions
    dynamic := alice.New(app.sessionManager.LoadAndSave)

    // Chain for protected routes (future: auth)
    protected := dynamic.Append(app.requireAuthentication)

    // Apply different chains to different routes
    mux.Handle("GET /static/", http.FileServer(...))           // No middleware
    mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))         // With session
    mux.Handle("GET /snippet/create", protected.ThenFunc(...)) // With auth

    return standard.Then(mux)
}
```

---

### 3. Composable Chains

You can extend existing chains with `Append()`:

```go
// Base chain
base := alice.New(app.recoverPanic, app.logRequest)

// Extend for different needs
withSession := base.Append(app.sessionManager.LoadAndSave)
withAuth := withSession.Append(app.requireAuthentication)
withCSRF := withAuth.Append(noSurf)

// Each chain builds on the previous
// withCSRF = recoverPanic → logRequest → session → auth → csrf
```

---

### 4. Two Ways to Apply

```go
// ThenFunc - for http.HandlerFunc
mux.Handle("GET /", chain.ThenFunc(app.home))

// Then - for http.Handler
mux.Handle("GET /static/", chain.Then(fileServer))
```

---

## Real Example from Snippetbox

**File:** `cmd/web/routes.go`

```go
func (app *application) routes() http.Handler {
    mux := http.NewServeMux()

    // Static files - no session needed
    fileServer := http.FileServer(http.Dir("./ui/static"))
    mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

    // Dynamic routes - need session for flash messages
    dynamic := alice.New(app.sessionManager.LoadAndSave)
    mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
    mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
    mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))
    mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

    // Global middleware - applies to everything
    standard := alice.New(app.recoverPanic, app.logRequest, commonHeader)
    return standard.Then(mux)
}
```

**Request flow:**

```
Request
   ↓
┌─────────────────────────────────────────────────────┐
│ standard chain (global)                             │
│   recoverPanic → logRequest → commonHeader          │
│              ↓                                      │
│   ┌─────────────────────────────────────────────┐   │
│   │ mux (router)                                │   │
│   │         ↓                                   │   │
│   │   /static/* → fileServer (no extra chain)  │   │
│   │   /         → dynamic chain → home handler │   │
│   │   /snippet  → dynamic chain → handler      │   │
│   └─────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────┘
   ↓
Response
```

---

## Without vs With Alice Comparison

| Aspect           | Without Alice      | With Alice                   |
| ---------------- | ------------------ | ---------------------------- |
| Readability      | `a(b(c(handler)))` | `New(a, b, c).Then(handler)` |
| Order            | Right-to-left      | Left-to-right                |
| Reusability      | Copy-paste chains  | Named, reusable chains       |
| Composability    | Manual nesting     | `Append()` method            |
| Per-route chains | Messy              | Clean with `ThenFunc()`      |

---

## Common Patterns

### Pattern 1: Global + Per-Route Chains

```go
standard := alice.New(recoverPanic, logRequest, headers)  // Global
dynamic := alice.New(session.LoadAndSave)                 // Per-route

mux.Handle("/api", dynamic.ThenFunc(apiHandler))
return standard.Then(mux)
```

### Pattern 2: Building Up Chains

```go
base := alice.New(recoverPanic, logRequest)
withSession := base.Append(sessionMiddleware)
withAuth := withSession.Append(authMiddleware)
withAdmin := withAuth.Append(adminOnlyMiddleware)

mux.Handle("/", withSession.ThenFunc(home))
mux.Handle("/dashboard", withAuth.ThenFunc(dashboard))
mux.Handle("/admin", withAdmin.ThenFunc(admin))
```

### Pattern 3: Conditional Middleware

```go
chain := alice.New(recoverPanic, logRequest)

if config.Debug {
    chain = chain.Append(debugMiddleware)
}

return chain.Then(mux)
```

---

## Summary

| Feature               | Description                       |
| --------------------- | --------------------------------- |
| `alice.New(...)`      | Create a new chain                |
| `chain.Append(...)`   | Add middleware to existing chain  |
| `chain.Then(handler)` | Apply chain to `http.Handler`     |
| `chain.ThenFunc(fn)`  | Apply chain to `http.HandlerFunc` |

**Why use Alice?**

1. ✅ Clean, readable middleware chains
2. ✅ Left-to-right execution order
3. ✅ Reusable and composable chains
4. ✅ Easy to add per-route middleware
5. ✅ No magic, just simple composition
