# Repository Skill

## Purpose

Enforce the repository implementation pattern used across all domain contexts in GoChat backend.

## The Triad Pattern

Every entity persisted to MySQL follows a strict three-part structure:

```
infrastructure/persistence/
├── model/
│   └── <entity>.go          # GORM model struct
├── converter/               # Entity ↔ Model conversion (optional, can be on repo)
└── repository/
    └── <entity>_repository.go  # Repository implementation
```

## GORM Model

### Rules

- Models live in `infrastructure/persistence/model/`
- Use domain value object types as fields, NOT primitives
- Every model has a `TableName()` method
- Use anonymous struct fields with GORM tags for constraints

```go
// model/friend_request.go
type FriendRequest struct {
    ID           kernel.ID        `gorm:"column:id;type:char(36);primaryKey"`
    SenderID     kernel.UserID    `gorm:"column:sender_id;type:char(36);not null;index"`
    ReceiverID   kernel.UserID    `gorm:"column:receiver_id;type:char(36);not null;index"`
    Status       string           `gorm:"column:status;type:varchar(20);not null"`
    Message      string           `gorm:"column:message;type:varchar(200)"`
    CreatedAt    time.Time        `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt    time.Time        `gorm:"column:updated_at;autoUpdateTime"`
}

func (FriendRequest) TableName() string { return "friend_requests" }
```

### Key Conventions

| Convention | Example |
|---|---|
| PK type | `CHAR(36)` for UUIDs |
| FK type | `CHAR(36)` matching the referenced PK |
| Timestamps | `autoCreateTime` / `autoUpdateTime` |
| Collation | `utf8mb4_unicode_ci` (set at table level in migrations) |

## Converter

### Entity → Model (`toModel`)

Always a private method. Uses `LoadXxx()` factory for the inverse:

```go
func (r *FriendRequestRepository) toModel(entity *domain.FriendRequest) *model.FriendRequest {
    return &model.FriendRequest{
        ID:         kernel.ID(entity.ID()),
        SenderID:   entity.SenderID(),
        ReceiverID: entity.ReceiverID(),
        Status:     entity.Status().String(),
        Message:    entity.Message(),
        CreatedAt:  entity.CreatedAt(),
    }
}
```

### Model → Entity (`toDomain`)

Always uses the domain entity's `LoadXxx()` factory:

```go
func (r *FriendRequestRepository) toDomain(m *model.FriendRequest) *domain.FriendRequest {
    return domain.LoadFriendRequest(
        m.ID, m.SenderID, m.ReceiverID,
        domain.FriendRequestStatusFromString(m.Status),
        m.Message, m.CreatedAt, m.UpdatedAt,
    )
}
```

If multiple models return, provide a `toDomains` plural variant:

```go
func (r *FriendRequestRepository) toDomains(models []*model.FriendRequest) []*domain.FriendRequest {
    result := make([]*domain.FriendRequest, len(models))
    for i, m := range models { result[i] = r.toDomain(m) }
    return result
}
```

## Repository Implementation

### Structure

```go
var _ domain.FriendRequestRepository = (*FriendRequestRepository)(nil)

type FriendRequestRepository struct {
    db        *gorm.DB
    eventRepo event.Repository
}

func NewFriendRequestRepository(db *gorm.DB, eventRepo event.Repository) *domain.FriendRequestRepository {
    return &FriendRequestRepository{db: db, eventRepo: eventRepo}
}
```

### Compile-Time Interface Check

ALWAYS include the compile-time assertion:

```go
var _ domain.XxxRepository = (*XxxRepository)(nil)
```

This catches interface mismatches at build time, not runtime.

### Write Operations (Transactional + Outbox)

Every write wraps in `db.Transaction()` and publishes events via `eventRepo`:

```go
func (r *FriendRequestRepository) Save(ctx context.Context, req *domain.FriendRequest) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        model := r.toModel(req)
        if err := tx.Create(model).Error; err != nil {
            return gormutils.TranslateError(err)
        }
        // Publish events in same transaction (Outbox pattern)
        events := req.Events() // drain from entity's event manager
        return r.eventRepo.CreateUnpublishedEvents(ctx, events...)
    })
}
```

### Read Operations

Use `db.WithContext(ctx)` for context propagation:

```go
func (r *FriendRequestRepository) FindByID(ctx context.Context, id kernel.ID) (*domain.FriendRequest, error) {
    var m model.FriendRequest
    err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error
    if err != nil {
        return nil, gormutils.TranslateError(err)
    }
    return r.toDomain(&m), nil
}
```

### Cursor Pagination

Use `WHERE id < ? ORDER BY id DESC LIMIT ?` pattern:

```go
func (r *PrivateMessageRepository) ListByUserID(
    ctx context.Context, userID kernel.UserID, cursor kernel.ID, limit int,
) ([]*domain.PrivateMessage, error) {
    var models []*model.PrivateMessage
    query := r.db.WithContext(ctx).
        Where("sender_id = ? OR recipient_id = ?", userID, userID)
    if cursor != "" {
        query = query.Where("id < ?", cursor)
    }
    err := query.Order("id DESC").Limit(limit).Find(&models).Error
    if err != nil {
        return nil, gormutils.TranslateError(err)
    }
    return r.toDomains(models), nil
}
```

### Error Translation

Always wrap GORM errors through `gormutils.TranslateError()`:

```go
func TranslateError(err error) error {
    switch {
    case errors.Is(err, gorm.ErrRecordNotFound):
        return errors.ErrNotFound
    case errors.Is(err, gorm.ErrDuplicatedKey):
        return errors.ErrDuplicatedKey
    // ... etc
    }
}
```

## Domain Entity LoadXxx() Factory

Domain entities use a `LoadXxx()` factory for reconstruction from persistence:

```go
// domain/friend_request.go
func LoadFriendRequest(
    id kernel.ID,
    senderID kernel.UserID,
    receiverID kernel.UserID,
    status FriendRequestStatus,
    message string,
    createdAt time.Time,
    updatedAt time.Time,
) *FriendRequest {
    return &FriendRequest{
        id:         id,
        senderID:   senderID,
        receiverID: receiverID,
        status:     status,
        message:    message,
        createdAt:  createdAt,
        updatedAt:  updatedAt,
    }
}
```

`LoadXxx()` is the ONLY way to reconstruct an entity from persistence. The entity's `NewXxx()` constructor is for creating NEW entities only.

## Common Mistakes

- ❌ Using primitives instead of domain value objects in GORM models
- ❌ Missing `var _ domain.XxxRepository = (*XxxRepository)(nil)` compile-time check
- ❌ Publishing events outside the transaction
- ❌ Forgetting `gormutils.TranslateError()` on DB errors
- ❌ Using `NewXxx()` for persistence reconstruction (use `LoadXxx()`)
- ❌ Raw SQL instead of GORM query builder (use GORM unless performance-critical)
