# Frontend Service Components

This document outlines the essential frontend components required for the various services within the social network platform.

## Authentication Components
These components handle user access and session management.

*   **RegisterForm**: A form for new users to create an account, including fields for username, email, and password.
*   **LoginForm**: A form for existing users to authenticate with their credentials.
*   **LogoutButton**: A component that triggers the logout process and clears the user session.
*   **ProtectedRoute**: A higher-order component or route wrapper that ensures only authenticated users can access specific pages.

## User Components
Components focused on displaying and managing user-specific information.

*   **CurrentUserProfile**: Displays the profile details of the currently logged-in user, including their bio, avatar, and account settings.
*   **UserCard**: A reusable summary component for displaying a user's basic information (e.g., in search results or lists).

## Follow Components
Components that manage user-to-user social connections and follow logic.

*   **FollowAction**: A button or action component that allows a user to follow or unfollow another user.
*   **FollowRequestCard**: Displays an individual follow request, typically with options to accept or decline.
*   **FollowRequestsList**: A list or page showing all pending follow requests for the current user.
*   **FollowersList**: Displays a list of users who are following a specific profile.
*   **FollowingList**: Displays a list of users that a specific profile is following.
*   **FollowStats**: Displays counts of followers and following for a user profile.
