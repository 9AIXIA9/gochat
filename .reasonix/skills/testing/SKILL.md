# Testing Skill

## Purpose

Enforce testing conventions: file organization, testify patterns, mockgen, table-driven tests, and per-layer strategies.

## File Organization

- Test files co-located with source: `<source>_test.go`
- Same package as source (white-box for domain, application, infrastructure)
- Mock files in `mocks/` subdirectory of **wherever the interface is defined**: `mock_<interface>.go`
- The rule: `<package>/mocks/mock_<interface>.go` — mocks live next to the source interface, not in a central mock directory

```
# Interfaces defined in domain → mocks in domain/mocks/
internal/friendship/domain/
├── friendship_repository.go        # defines FriendshipRepository interface
├── friend_request.go
├── friend_request_test.go
└── mocks/
    └── mock_friendship_repository.go

# Interfaces defined in shared → mocks in shared/<pkg>/mocks/
internal/shared/command/
├── command.go                       # defines Command interface
├── ports.go                         # defines SyncPublisher interface
└── mocks/
    ├── mock_command.go
    └── mock_ports.go

internal/shared/contract/
├── gateway_service.go               # defines GatewayService interface
└── mocks/
    └── mock_gateway_service.go

# Infrastructure interfaces → mocks in infrastructure/<pkg>/mocks/
internal/infrastructure/gin/
├── ports.go                         # defines Bindable, Validator interfaces
└── mocks/
    └── mock_ports.go

# Gateway interfaces → mocks in gateway/core/mocks/
internal/gateway/core/
├── ports.go                         # defines SessionIDGenerator interface
└── mocks/
    └── mock_ports.go
```

## Test Framework

### testify/assert

```go
import "github.com/stretchr/testify/assert"

assert.Equal(t, expected, actual)      // expected first
assert.NoError(t, err)
assert.ErrorIs(t, err, targetErr)      // sentinel error check
assert.True(t, condition)
assert.NotNil(t, value)
```

### testify/require

Use `require.*` when the test cannot continue on failure:

```go
require.NoError(t, err)    // stop test if err
require.NotNil(t, result)  // stop test if nil
```

### testify/mock

```go
import "github.com/stretchr/testify/mock"

mockRepo := new(mocks.MockFriendRequestRepository)
mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*domain.FriendRequest")).Return(nil)

// Execute...
mockRepo.AssertExpectations(t)
```

## Mock Generation

Use `mockgen` with `//go:generate` directives on the **interface source file**:

```go
// domain/friendship_repository.go
//go:generate mockgen -source=friendship_repository.go -destination=./mocks/mock_friendship_repository.go -package=mocks

// shared/command/command.go
//go:generate mockgen -source=command.go -destination=./mocks/mock_command.go -package=mocks

// shared/contract/gateway_service.go
//go:generate mockgen -source=gateway_service.go -destination=./mocks/mock_gateway_service.go -package=mocks
```

### Rules

- The `//go:generate` directive goes on the **interface definition file**, not on the mock
- `-source` points to the current file (the interface)
- `-destination` is always `./mocks/mock_<source_file>.go`
- `-package` is always `mocks`
- Run `go generate ./...` from repo root to regenerate all mocks

### Mock Locations (Complete)

| Interface defined in | Mock location |
|---|---|
| `<context>/domain/<file>.go` | `<context>/domain/mocks/mock_<file>.go` |
| `shared/command/<file>.go` | `shared/command/mocks/mock_<file>.go` |
| `shared/event/<file>.go` | `shared/event/mocks/mock_<file>.go` |
| `shared/contract/<file>.go` | `shared/contract/mocks/mock_<file>.go` |
| `shared/kernel/<file>.go` | `shared/kernel/mocks/mock_<file>.go` |
| `gateway/core/<file>.go` | `gateway/core/mocks/mock_<file>.go` |
| `infrastructure/gin/<file>.go` | `infrastructure/gin/mocks/mock_<file>.go` |

## Table-Driven Tests

Preferred pattern for multiple scenarios:

```go
func TestXxx(t *testing.T) {
    tests := []struct {
        name    string
        input   Input
        want    Output
        wantErr error
    }{
        {"valid case", validInput, expectedOutput, nil},
        {"empty input", emptyInput, Output{}, ErrEmptyInput},
        {"already exists", dupInput, Output{}, ErrAlreadyExists},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := doSomething(tt.input)
            if tt.wantErr != nil {
                assert.ErrorIs(t, err, tt.wantErr)
                return
            }
            assert.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

## Per-Layer Strategy

### Domain Tests

Test business logic directly. No mocks needed for entities and value objects:

```go
func TestFriendRequest_Accept(t *testing.T) {
    req := domain.NewFriendRequest(fromID, toID, msg)
    err := req.Accept()
    assert.NoError(t, err)
    assert.Equal(t, domain.Accepted, req.Status())
}
```

Also test domain errors in `consts_test.go` or separate error test files.

### Application Tests

Mock domain interfaces (from `<context>/domain/mocks/`) to test use case orchestration:

```go
import "gochat/internal/friendship/domain/mocks"

func TestSendFriendRequestUseCase_Execute(t *testing.T) {
    mockRepo := new(mocks.MockFriendRequestRepository)
    mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
    
    uc := application.NewSendFriendRequestUseCase(mockRepo, eventRepo)
    output, err := uc.Execute(ctx, input)
    
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

When the use case depends on shared interfaces, import from the corresponding `shared/<pkg>/mocks/`:

```go
import cmdMocks "gochat/internal/shared/command/mocks"
import evtMocks "gochat/internal/shared/event/mocks"
```

### Port/HTTP Tests

Use `pkg/httptest` and Gin test mode. Mock the use case (from application layer) and the validator (from `infrastructure/gin/mocks/`):

```go
import ginMocks "gochat/internal/infrastructure/gin/mocks"

func TestXxxHandler_Handle(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    // Mock the use case
    mockUC := new(applicationMocks.MockXxxUseCase)
    mockUC.On("Execute", mock.Anything, mock.Anything).Return(output, nil)
    
    // Mock the validator
    mockValidator := new(ginMocks.MockValidator)
    mockValidator.On("Validate", mock.Anything, mock.Anything).Return("", nil)
    
    handler := NewXxxHandler(mockUC, mockValidator)
    router := gin.New()
    router.POST("/api/v1/xxx", handler)
    
    w := httptest.ExecuteRouteTest(t, router, &httptest.TestCase{
        Method:   "POST",
        Path:     "/api/v1/xxx",
        Body:     requestBody,
        Response: expectedResponse,
    })
    assert.Equal(t, 200, w.Code)
}
```

### Repository Tests

Use a real test database or skip (integration tests). Repository unit tests are typically not required for GORM wrappers.

## What to Test

- ✅ Domain entities and value objects (business rules)
- ✅ Use cases (orchestration with mocks)
- ✅ HTTP handlers (request parsing, response formatting, error mapping)
- ✅ Domain errors (correct sentinel values)
- ❌ Simple getters/setters
- ❌ Wire-generated code (`wire_gen.go`)
- ❌ Config struct `Validate()` if it's just field checks
- ❌ Repository method bodies (GORM wrappers — integration test instead)

## Running Tests

| Command | Scope |
|---|---|
| `make test` | All tests |
| `make test-race` | All tests with race detector |
| `go test ./internal/friendship/...` | Single context |
| `go test -v -run TestName ./path/` | Single test verbose |
| `go test -cover ./...` | With coverage |

## Test Helpers in pkg/

| Package | Purpose |
|---|---|
| `pkg/httptest` | `NewTestRouter(t)`, `ExecuteRouteTest(t, router, testCase)` |
| `pkg/validate` | `NotNil(values...)` — nil check helper |
| `pkg/ctxutil` | `WithUserID(ctx, id)` — inject user into context for tests |
