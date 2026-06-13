# Social Network Frontend

The frontend is a React 19 application built with Vite, React Router,
React Icons and plain JavaScript/JSX.
See [`../docs/openapi.json`](../docs/openapi.json) for request and response shapes.

## Commands

```bash
npm install
npm run dev
npm run lint
npm run build
```

The app runs at `http://localhost:5173` during development and calls the backend at `http://localhost:8080/api`. Requests to protected routes must include credentials so the `session_token` cookie is sent.

## Conventions

- Keep API field names aligned with `../docs/openapi.json`.
- Use ordinary JavaScript objects for post and comment DTOs.
  JSDoc may be used where it improves editor support.

## Posting Behavior

- Deleted posts are minimal tombstones.
  Deleted comments remain in the recursive thread so replies keep their context.
- Authorized users may read an existing comment thread after its post is deleted,
  but the UI must hide creation, edit, and vote controls for that deleted post.
