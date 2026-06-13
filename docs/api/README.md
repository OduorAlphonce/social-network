# API

[`openapi.json`](openapi.json) contains the OpenAPI 3.1 contract for the target
Social Network API.

The document is a design specification, not a guarantee that every path is
currently available. The running backend currently registers:

- `POST /api/users/register`
- `POST /api/users/login`
- `POST /api/users/logout`
- `GET /api/users/me`
- Follow, unfollow, accept, reject, followers, and following operations under
  `/api/followers/`

Consult `backend/internal/api/routers/router.go` when checking runtime endpoint
availability. The local API base URL is `http://localhost:8080/api`.
