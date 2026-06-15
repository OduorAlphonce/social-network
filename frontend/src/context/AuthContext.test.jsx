import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { AuthProvider } from "./AuthContext";
import { useAuth } from "./useAuth";
import { apiFetch } from "../utils/api";

vi.mock("../utils/api", async (importOriginal) => {
  const original = await importOriginal();

  return {
    ...original,
    apiFetch: vi.fn(),
  };
});

const AuthStatus = () => {
  const { currentUser, isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return <p>Loading session</p>;
  }

  return (
    <p>
      {isAuthenticated ? `Signed in as ${currentUser.username}` : "Signed out"}
    </p>
  );
};

describe("AuthProvider", () => {
  it("loads the current user and exposes an authenticated session", async () => {
    apiFetch.mockResolvedValueOnce({ id: 1, username: "ada" });

    render(
      <AuthProvider>
        <AuthStatus />
      </AuthProvider>
    );

    expect(screen.getByText("Loading session")).toBeInTheDocument();
    expect(await screen.findByText("Signed in as ada")).toBeInTheDocument();
    expect(apiFetch).toHaveBeenCalledWith("/api/users/me");
  });
});
