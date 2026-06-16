# HTTP Handler Skill

## Purpose

Enforce the HTTP handler conventions: generic handler adapter, request binding, validation, response formatting, and Swagger documentation.

## The Generic Handler Adapter

The project uses a **generic handler adapter** instead of writing boilerplate handler code. This is the primary pattern for all HTTP handlers.

### AdaptUseCaseToHandler

Defined in `internal/infrastructure/gin/handler_adapter.go`:

```go
func AdaptUseCaseToHandler[
    Request any,
    RequestPointer RequestPointers[Request],
    Input kernel.Validatable,
    Output any,
](
    useCase kernel.UseCase[Input, Output],
    validator Validator,
    convertRequestToInput func(RequestPointer) Input,
    handleOutput func(*gin.Context, Output),
) gin.HandlerFunc
```

### How It Works

1. **Bind** → `request.Bind(ginContext)` — parse JSON body/query into Request struct
2. **Validate request** → `validator.Validate(ctx, request)` — run validator tags, return translated error messages
3. **Validate input** → `input.Validate()` — domain-level validation
4. **Execute use case** → `useCase.Execute(ctx, input)`
5. **Handle output** → `handleOutput(ctx, output)` — format success response
6. **Error handling** → `BusinessError` → `CodeBusinessError`; other errors → `CodeServerError`

### Request Pointers Interface

```go
type RequestPointers[Request any] interface {
    *Request
    Bindable
}

type Bindable interface {
    Bind(ginContext *gin.Context) error
}
```

Request structs satisfy `Bindable` by implementing `Bind`:

```go
func (r *SendFriendRequestRequest) Bind(c *gin.Context) error {
    return c.ShouldBindJSON(r)
}
```

### Usage Example

```go
func NewSendFriendRequestHandler(uc *application.SendFriendRequestUseCase, v gin.Validator) gin.HandlerFunc {
    return ginutils.AdaptUseCaseToHandler[
        SendFriendRequestRequest,
        *SendFriendRequestRequest,
        application.SendFriendRequestInput,
        application.SendFriendRequestOutput,
    ](
        uc,
        v,
        func(req *SendFriendRequestRequest) application.SendFriendRequestInput {
            return application.SendFriendRequestInput{
                FromUserID: req.FromUserID,
                ToUserID:   req.ToUserID,
                Message:    req.Message,
            }
        },
        func(c *gin.Context, output application.SendFriendRequestOutput) {
            ginutils.ResponseSuccessWithData(c, output)
        },
    )
}
```

## Manual Handler Pattern (for non-standard cases)

When the generic adapter doesn't fit, use the manual pattern:

```go
type XxxHandler struct {
    useCase *application.XxxUseCase
}

func (h *XxxHandler) Handle(c *gin.Context) {
    // 1. Parse
    var req XxxRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        ginutils.Response(c, api.CodeInvalidParam)
        return
    }
    // 2. Build input
    input := application.XxxInput{...}
    // 3. Validate
    if err := input.Validate(); err != nil {
        ginutils.Response(c, api.CodeInvalidParam)
        return
    }
    // 4. Execute
    output, err := h.useCase.Execute(c.Request.Context(), input)
    if err != nil {
        if errors.IsBusinessError(err) {
            ginutils.ResponseWithMessage(c, api.CodeBusinessError, err.Error())
        } else {
            zap.L().Error("handler failed", zap.Error(err))
            ginutils.Response(c, api.CodeServerError)
        }
        return
    }
    // 5. Respond
    ginutils.ResponseSuccessWithData(c, output)
}
```

## Validator

Implemented in `internal/infrastructure/validator/validator.go`:

```go
type Validator struct { ... }

func (v *Validator) Validate(ctx context.Context, model any) (string, error)
```

- Registers custom English translations for validation tags
- Returns translated error messages as a string (empty if valid)
- Supports: `required`, `min`, `max`, and custom rules

### Validation Tags

Request structs use `binding` tags:

```go
type SendFriendRequestRequest struct {
    ToUserID string `json:"to_user_id" binding:"required,uuid"`
    Message  string `json:"message"    binding:"max=200"`
}
```

Common bindings: `required`, `uuid`, `email`, `min=N`, `max=N`, `oneof=a b c`.

## Response Format

Always use `ginutils.Response*` helpers, never raw `c.JSON`:

| Helper | Use case |
|---|---|
| `ginutils.Response(c, code)` | Standard response from Code enum |
| `ginutils.ResponseSuccess(c)` | Success with no data |
| `ginutils.ResponseSuccessWithData(c, data)` | Success with payload |
| `ginutils.ResponseWithMessage(c, code, msg)` | Error with custom message |

### Standard Response Shape

```json
{
    "code": 200,
    "message": "success",
    "data": {}
}
```

## Route Registration

Routes are registered in `cmd/api/di/providers_http.go`:

```go
v1 := router.Group("/api/v1")
v1.Use(rateLimitMiddleware, breakerMiddleware, timeoutMiddleware)

{
    auth := v1.Group("/auth")
    auth.POST("/sign-up", signUpHandler)
    auth.POST("/login/user_number", loginByNumberHandler)
    // ...
}
{
    friends := v1.Group("/friendships")
    friends.Use(authMiddleware)
    friends.GET("/", listFriendshipsHandler)
    // ...
}
```

Prefix: `/api/v1`. Auth routes are unauthenticated. Most others require `authMiddleware`.

## Swagger Annotations

### Top-Level (in `cmd/api/main.go`)

```go
// @title           GoChat Backend API
// @version         1.0
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
```

### Per-Handler

```go
// @Summary      Send friend request
// @Description  Send a friend request to another user
// @Tags         Friendship
// @Accept       json
// @Produce      json
// @Param        request body SendFriendRequestRequest true "Friend request payload"
// @Success      200 {object} api.Response
// @Failure      400 {object} api.Response
// @Security     BearerAuth
// @Router       /api/v1/friend-requests [post]
```

### Generation

```bash
make swagger      # Generate docs
make swagger-fix   # Fix example fields
```

Output: `docs/docs.go`, `docs/swagger.json`, `docs/swagger.yaml`.

## Context Propagation

Always pass `c.Request.Context()` to use cases — never `context.Background()`:

```go
output, err := h.useCase.Execute(c.Request.Context(), input)
```

This enables timeout propagation from the `Timeout` middleware and trace context from OTEL.
