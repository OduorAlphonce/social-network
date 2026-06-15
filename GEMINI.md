# Social Network Frontend Implementation Guide

This document outlines the frontend components and services required to integrate with the social network backend.

## Architecture Overview
- **Framework:** React (Vite)
- **State Management:** (To be determined, likely React Context or Redux)
- **Routing:** React Router

## Component Breakdown

### Authentication
- `RegisterForm`: Handles new user registration, including validation and submission to `/api/register`.
- `LoginForm`: Handles user login, manages JWT storage (localStorage/Cookies), and redirects on success.
- `LogoutButton`: Clears authentication tokens and redirects to the login page.
- `ProtectedRoute`: A HOC or Wrapper component to restrict access to authenticated routes.

### User Components
- `CurrentUserProfile`: Displays the authenticated user's profile information, posts, and settings.
- `UserCard`: A reusable card component to display brief user information (avatar, name, bio) in lists.

### Follow Components
- `FollowAction`: A button component that toggles between "Follow", "Unfollow", or "Requested" based on the relationship with the target user.
- `FollowRequestCard`: Displays a pending follow request with "Accept" and "Decline" actions.
- `FollowRequestsList`: A list view for managing incoming follow requests.
- `FollowersList`: Displays a list of users following the current user or a target user.
- `FollowingList`: Displays a list of users the current user or a target user is following.
- `FollowStats`: Displays counts for followers and following on profile pages.

## Implementation Tasks
- [x] Implement Authentication components (RegisterForm)
- [ ] Implement User profile and card components
- [ ] Implement Follow system components and logic
