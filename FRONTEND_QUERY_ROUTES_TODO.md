# Frontend Query Routes TODO

This document outlines the missing query and management routes needed to support a web or mobile frontend, following the existing project routing style and module boundaries.

Status legend:
- Implemented: handler + usecase exist and are already wired.
- Missing: handler and/or usecase need implementation, then route wiring in `cmd/di/providers_http.go`.

## Profile Module

Implemented:
- PUT `/api/v1/profile/user` — Update current user profile.
- PUT `/api/v1/profile/room` — Update room profile.

Missing (to implement):
- GET `/api/v1/profile/user/me` — Get current logged-in user's profile.
- GET `/api/v1/profile/user/:user_id` — Get a user's public profile by ID.
- GET `/api/v1/profile/room/:room_id` — Get a room's public profile by ID.
- GET `/api/v1/profile/user/search?q=...` — Search users by keyword.

Suggested contracts:
- Auth: required.
- Params: path `:user_id`, `:room_id`; query `q` for search.
- Responses: standard envelope `{ data, error }` following your HTTP port conventions.

## Roomship Module

Implemented:
- POST `/api/v1/roomship/room` — Create a room.
- GET `/api/v1/roomship/` — List rooms the user is in.
- GET `/api/v1/roomship/request` — List incoming/outgoing member requests.
- POST `/api/v1/roomship/request` — Send a member request.
- PUT `/api/v1/roomship/request/:request_id/agree` — Agree a member request.
- PUT `/api/v1/roomship/request/:request_id/refuse` — Refuse a member request.

Missing (to implement):
- DELETE `/api/v1/roomship/:room_id` — Leave room.
- GET `/api/v1/roomship/:room_id/member` — List room members.

Suggested contracts:
- Auth: required.
- Params: path `:room_id`, `:request_id`.
- Responses: members list with pagination (optional: `page`, `page_size`).

## Friendship Module

Implemented:
- GET `/api/v1/friendship/` — List friendships.
- GET `/api/v1/friendship/request` — List friend requests.
- POST `/api/v1/friendship/request` — Send friend request.
- PUT `/api/v1/friendship/request/:request_id/agree` — Agree friend request.
- PUT `/api/v1/friendship/request/:request_id/refuse` — Refuse friend request.

Missing (to implement):
- DELETE `/api/v1/friendship/:friendship_id` — Delete friendship.

Suggested contracts:
- Auth: required.
- Params: path `:friendship_id`.
- Responses: `{ success: true }` or standard envelope.

## Chat Module (Reference)

Implemented:
- GET `/api/v1/chat/private` — List private messages.
- GET `/api/v1/chat/room` — List room messages.
- POST `/api/v1/chat/private` — Send private message.
- POST `/api/v1/chat/room` — Send room message.

## Notification Module (Reference)

Implemented:
- GET `/api/v1/notification/system` — List system messages.

## Wiring Guidance

Once usecases and handlers are implemented under `internal/<module>/application` and `internal/<module>/port/http`, wire the routes in `cmd/di/providers_http.go` using the established pattern:

```go
profileGroup.GET("/user/me", profileHTTP.NewGetCurrentUserProfileHandler(getCurrentUserProfile, validator))
profileGroup.GET("/user/:user_id", profileHTTP.NewGetUserProfileHandler(getUserProfile, validator))
profileGroup.GET("/room/:room_id", profileHTTP.NewGetRoomProfileHandler(getRoomProfile, validator))
profileGroup.GET("/user/search", profileHTTP.NewSearchUsersHandler(searchUsers, validator))

roomshipGroup.DELETE("/:room_id", roomshipHTTP.NewLeaveRoomHandler(leaveRoom, validator))
roomshipGroup.GET("/:room_id/member", roomshipHTTP.NewListRoomMembersHandler(listRoomMembers, validator))

friendshipGroup.DELETE("/:friendship_id", friendshipHTTP.NewDeleteFriendshipHandler(deleteFriendship, validator))
```

Ensure corresponding usecases are added to the DI providers in `cmd/di/providers_usecase.go` and injected into `provideHttpRouter` parameters.

## Implementation Checklist

- Profile
  - [ ] Usecases: GetCurrentUserProfile, GetUserProfile, GetRoomProfile, SearchUsers
  - [ ] Handlers: NewGetCurrentUserProfileHandler, NewGetUserProfileHandler, NewGetRoomProfileHandler, NewSearchUsersHandler
  - [ ] Wire routes in `providers_http.go`

- Roomship
  - [ ] Usecases: LeaveRoom, ListRoomMembers
  - [ ] Handlers: NewLeaveRoomHandler, NewListRoomMembersHandler
  - [ ] Wire routes in `providers_http.go`

- Friendship
  - [ ] Usecases: DeleteFriendship
  - [ ] Handlers: NewDeleteFriendshipHandler
  - [ ] Wire route in `providers_http.go`

## Notes

- Keep validation via `validator` consistent with existing handlers.
- Reuse `authorizationMiddleware` for protected endpoints.
- Follow response and error conventions in `internal/delivery/http`.
- Add tests under each module's `application` and `port/http` as done elsewhere.

