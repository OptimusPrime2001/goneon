# Golang Learning Roadmap

## 1. Basic Review

### Variables & Constants

* `var`
* `:=`
* `const`

### Data Types

* Numeric types
* String
* Boolean
* Type inference
* Type conversion
* Zero value

### Custom Types

* Type definition
* Type alias
* Underlying type

#### Exercises

* Temperature converter
* Currency converter
* Type alias practice

---

## 2. Aggregate Types

### Array

* Declaration
* Length
* Comparison
* Multidimensional Array

### Struct

* Declaration
* Anonymous Struct
* Nested Struct
* Struct Tag
* Struct Comparison

#### Exercises

* Student Management
* Employee Information
* Nested Address Struct

---

## 3. Reference Types

### Slice

* `len()`
* `cap()`
* `append()`
* `copy()`
* slicing
* underlying array
* nil slice
* memory sharing
* deep copy

### Map

* Create
* Read
* Update
* Delete
* Check existence
* Iteration

### Pointer

* Address operator (`&`)
* Dereference (`*`)
* Pass by value vs pointer
* Pointer receiver

### Function

* Parameters
* Return values
* Multiple return values
* Variadic functions
* Anonymous functions
* Closures

#### Exercises

* Slice utilities
* Word counter
* Pointer swap
* Closure counter

---

## 4. Control Flow

### Conditional

* if / else
* switch

### Loop

* for
* range
* break
* continue
* goto (optional)

#### Exercises

* FizzBuzz
* Prime number checker
* Fibonacci
* Palindrome checker
* Factorial

---

## 5. Methods

### Method Basics

* Method declaration
* Value receiver
* Pointer receiver
* Method sets

#### Exercises

* Rectangle Area/Perimeter
* Circle Area
* Bank Account Deposit/Withdraw

---

## 6. Interfaces

### Interface Basics

* Interface declaration
* Implicit implementation
* Empty interface (`any`)
* Type assertion
* Type switch

#### Exercises

* Shape interface
* Payment interface
* Notification interface

---

## 7. Error Handling

### Errors

* `errors.New()`
* `fmt.Errorf()`
* Error wrapping
* `errors.Is()`
* `errors.As()`

### Panic & Recover

* panic
* recover

#### Exercises

* Divide function
* File validation
* Custom errors

---

## 8. Packages & Modules

### Package System

* Package structure
* Exported identifiers
* Internal package

### Modules

* `go mod init`
* `go mod tidy`
* Dependency management

#### Exercises

* Utility package
* Calculator package

---

## 9. Generics

### Generic Basics

* Type parameters
* Constraints
* Generic functions
* Generic structs

#### Exercises

* Generic Sum
* Generic Stack
* Generic Repository

---

## 10. Memory Model

### Stack vs Heap

* Variable lifetime
* Heap allocation

### Escape Analysis

* Escape analysis basics
* `go build -gcflags="-m"`

### Value Semantics

* Copy behavior
* Pointer behavior

#### Exercises

* Analyze allocations
* Compare value vs pointer performance

---

## 11. Concurrency

### Goroutines

* Creating goroutines
* Synchronization basics

### Channels

* Unbuffered channel
* Buffered channel
* Directional channel
* Closing channels

### Select

* Multiple channel handling
* Timeout pattern

### Synchronization

* WaitGroup
* Mutex
* RWMutex
* Atomic

### Context

* Context basics
* Cancellation
* Timeout

### Concurrency Patterns

* Worker Pool
* Fan-Out
* Fan-In
* Pipeline
* Producer Consumer

### Common Problems

* Race Condition
* Deadlock
* Goroutine Leak

#### Exercises

* Concurrent API calls
* Worker Pool
* Rate Limiter
* Pipeline processing

---

## 12. Runtime Internals

### Scheduler

* G
* M
* P
* Work stealing

### Garbage Collector

* Mark & Sweep
* GC behavior

### Memory Allocation

* `new()`
* `make()`
* Allocation strategies

#### Exercises

* GC observation
* Scheduler experiments

---

## 13. Standard Library

### Text Processing

* strings
* strconv
* bytes
* bufio

### Collections

* slices
* maps
* sort

### File System

* os
* io
* filepath

### Time

* time

### Context & Sync

* context
* sync
* sync/atomic

### Logging

* log
* log/slog

#### Exercises

* Log parser
* File scanner
* Text processor

---

## 14. File Handling

### File Operations

* Read file
* Write file
* Append file

### Data Formats

* JSON
* CSV
* XML

#### Exercises

* Todo CLI
* CSV parser
* JSON config loader

---

## 15. Testing

### Unit Testing

* Test functions
* Assertions

### Table Driven Test

* Test cases
* Subtests

### Benchmark

* Benchmark functions
* Performance analysis

### Mocking

* Interface-based mocks

### Race Detection

* `go test -race`

#### Exercises

* Test Calculator
* Test Repository
* Benchmark String Builder

---

## 16. HTTP Server

### net/http

* Handlers
* ServeMux

### HTTP Concepts

* Request
* Response
* Headers
* Status Codes

### Middleware

* Logging
* Authentication
* Recovery

### JSON API

* Encode
* Decode

#### Exercises

* Todo API
* User API
* Authentication Middleware

---

## 17. Database

### database/sql

* Connection
* Query
* Exec

### PostgreSQL

* CRUD
* Prepared Statements

### Transactions

* Commit
* Rollback

### Performance

* Connection Pool
* Query Optimization
* N+1 Query

### Pagination

* Offset Pagination
* Cursor Pagination

#### Exercises

* CRUD User
* CRUD Product
* Pagination API

---

## 18. Profiling & Optimization

### Benchmark Analysis

* `go test -bench .`

### Profiling

* pprof
* CPU profile
* Memory profile
* Goroutine profile

### Optimization

* Reduce allocations
* Efficient string handling

#### Exercises

* Optimize JSON processing
* Benchmark comparison

---

## 19. Tooling

### Go Commands

* `go run`
* `go build`
* `go install`
* `go test`
* `go fmt`
* `go vet`

### Dependency Management

* `go mod`
* `go work`

### Linting

* golangci-lint

#### Exercises

* Setup project tooling
* Configure linter

---

## 20. Advanced Topics

### Reflection

* Type inspection
* Dynamic values

### Dependency Injection

* Constructor Injection
* Interface Injection

### Design Patterns

* Factory
* Strategy
* Repository

### Architecture

* Clean Architecture
* Hexagonal Architecture

### Communication

* gRPC
* WebSocket
* Message Queue

#### Exercises

* Reflection validator
* gRPC service
* WebSocket chat

---

## 21. Production Engineering

### Configuration

* Environment variables
* Config files

### Logging

* Structured logging
* slog
* zap

### Docker

* Dockerfile
* Multi-stage build

### CI/CD

* GitHub Actions

### Observability

* Metrics
* Tracing
* Prometheus

#### Exercises

* Dockerize API
* Add Metrics
* CI/CD Pipeline

---

# Build Projects

## Beginner

* Calculator CLI
* Todo CLI
* Notes Manager
* Expense Tracker

## Intermediate

* REST API
* URL Shortener
* Blog API
* Authentication Service

## Advanced

* E-commerce Backend
* Chat Application
* Task Scheduler
* Distributed Worker
* Notification Service

---

# Recommended Learning Order

1. Basic Review
2. Aggregate Types
3. Reference Types
4. Control Flow
5. Methods
6. Interfaces
7. Error Handling
8. Packages
9. Generics
10. Memory Model
11. Concurrency
12. Runtime Internals
13. Standard Library
14. File Handling
15. Testing
16. HTTP Server
17. Database
18. Profiling
19. Tooling
20. Advanced Topics
21. Production Engineering
22. Build Real Projects
