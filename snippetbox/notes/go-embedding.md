# Go Struct Embedding (Composition Pattern)

Go doesn't have inheritance like Java/C#. Instead, it uses **embedding** to achieve composition.

---

## Basic Syntax

```go
type Validator struct {
    FieldErrors map[string]string
}

func (v *Validator) Valid() bool {
    return len(v.FieldErrors) == 0
}

// Embedding - no field name!
type snippetCreateForm struct {
    Title   string
    Content string
    Validator  // ← Embedded (no field name)
}
```

---

## With vs Without Embedding

| Syntax | How to Call |
|--------|-------------|
| `Validator` (embedded) | `form.Valid()` |
| `Validator Validator` (named field) | `form.Validator.Valid()` |

When embedded, the inner struct's methods and fields are **promoted** to the outer struct.

---

## Multiple Embedding

Go allows embedding multiple structs — something Java/C# can't do!

```go
type Validator struct {
    FieldErrors map[string]string
}
func (v *Validator) Valid() bool { return len(v.FieldErrors) == 0 }

type Timestamper struct {
    CreatedAt time.Time
}
func (t *Timestamper) SetCreatedNow() { t.CreatedAt = time.Now() }

// Embed BOTH structs!
type snippetCreateForm struct {
    Title       string
    Content     string
    Validator   // ← Embedded
    Timestamper // ← Also embedded
}

// Usage - call methods from both directly
func main() {
    form := snippetCreateForm{Title: "Hello"}
    form.Valid()         // From Validator
    form.SetCreatedNow() // From Timestamper
}
```

---

## Embedding vs Inheritance Comparison

| Feature | Java/C# Inheritance | Go Embedding |
|---------|---------------------|--------------|
| Syntax | `class A extends B` | `type A struct { B }` |
| Relationship | "is-a" | "has-a" (composition) |
| Multiple inheritance | ❌ Not supported | ✅ Can embed multiple |
| Override | `@Override` | No real override, uses "shadowing" |
| Access parent | `super.method()` | `form.Validator.Valid()` |
| Polymorphism | Via class hierarchy | Via **interfaces** |

---

## Shadowing Explained

### What is Shadowing?

In Java/C#, when you override a method, the child's method **replaces** the parent's method in the inheritance chain. The parent method is still there but accessed via `super`.

In Go, there's **no real override** because there's no inheritance. Instead, Go uses **shadowing**:

- If the outer struct defines a method with the **same name** as the embedded struct's method
  
- The outer method **shadows** (hides) the embedded method
  
- The embedded method still exists and can be called explicitly

### Shadowing Example

```go
type Animal struct{}

func (a *Animal) Speak() string {
    return "..."
}

type Dog struct {
    Animal  // Embed Animal
}

// This SHADOWS Animal.Speak(), it does NOT override it
func (d *Dog) Speak() string {
    return "Woof!"
}

func main() {
    dog := Dog{}
    
    fmt.Println(dog.Speak())        // "Woof!" - calls Dog.Speak()
    fmt.Println(dog.Animal.Speak()) // "..."   - calls the embedded Animal.Speak()
}
```

### Override vs Shadow - Key Difference

| Aspect | Java Override | Go Shadow |
|--------|---------------|-----------|
| Parent method | Replaced in dispatch | Still exists |
| Access original | `super.method()` | `outer.Embedded.Method()` |
| Polymorphism | Works with class type | Doesn't affect interface satisfaction |
| Runtime dispatch | Virtual method table | Static at compile time |

### Shadowing with Interfaces

Here's where it gets interesting. Shadowing **does** affect interface satisfaction:

```go
type Speaker interface {
    Speak() string
}

type Animal struct{}
func (a *Animal) Speak() string { return "..." }

type Dog struct {
    Animal
}
func (d *Dog) Speak() string { return "Woof!" }

func main() {
    var s Speaker
    
    s = &Dog{}
    fmt.Println(s.Speak())  // "Woof!" - Dog.Speak() is used
    
    s = &Animal{}
    fmt.Println(s.Speak())  // "..." - Animal.Speak() is used
}
```

When `Dog` implements `Speak()`, it shadows `Animal.Speak()`. When you treat `Dog` as a `Speaker` interface, it uses `Dog.Speak()`.

### Calling the Embedded Method from Shadow

You can call the original embedded method from your shadowing method:

```go
type Dog struct {
    Animal
}

func (d *Dog) Speak() string {
    original := d.Animal.Speak()  // Call embedded method
    return "Dog says: " + original + " Woof!"
}
```

---

## Embedding vs DI (Dependency Injection)

| Concept | Purpose | Example |
|---------|---------|---------|
| **Embedding** | Composition — decided at compile time | `snippetCreateForm` has validation methods |
| **DI** | Inject dependencies — usually at runtime | `application` receives `*models.SnippetModel` |

```go
// Embedding example
type snippetCreateForm struct {
    validator.Validator // Embedded - gives form validation methods
}

// DI example
type application struct {
    snippets *models.SnippetModel // Injected - provides database access
}
```

---

## Real-World Usage in Snippetbox

```go
// internal/validator/validator.go
type Validator struct {
    FieldErrors map[string]string
}

func (v *Validator) CheckField(ok bool, key, message string) {
    if !ok {
        v.AddFieldError(key, message)
    }
}

// cmd/web/handlers.go
type snippetCreateForm struct {
    Title               string `form:"title"`
    Content             string `form:"content"`
    Expires             int    `form:"expires"`
    validator.Validator `form:"-"` // Embedded!
}

// Usage in handler
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    var form snippetCreateForm
    // ...
    
    // Call embedded methods directly
    form.CheckField(validator.NotBlank(form.Title), "title", "Cannot be blank")
    
    if !form.Valid() {
        // Validation failed
    }
}
```

---

## Key Takeaways

1. **Embedding = Composition**, not inheritance

2. **No field name** = embedding

3. Embedded methods and fields are **promoted** to the outer struct

4. Can **embed multiple** structs

5. **Shadowing** hides embedded methods but doesn't replace them

6. Use `outer.Embedded.Method()` to access shadowed methods

7. Go prefers **"Composition over Inheritance"**
