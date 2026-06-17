import { fireEvent, screen, waitFor } from "@testing-library/react";
import { Route, Routes } from "react-router";
import { describe, expect, it, vi } from "vitest";
import LoginForm from "./LoginForm";
import { renderWithProviders } from "../test/render";

const renderLogin = ({ auth = {}, route = "/login" } = {}) =>
  renderWithProviders(
    <Routes>
      <Route path="/login" element={<LoginForm />} />
      <Route path="/" element={<p>Home feed</p>} />
      <Route path="/profile" element={<p>Profile page</p>} />
    </Routes>,
    { route, auth }
  );

describe("LoginForm", () => {
  it("submits credentials through auth and returns to the requested page", async () => {
    const login = vi.fn().mockResolvedValue({ id: "user-1" });
    renderLogin({
      route: {
        pathname: "/login",
        state: { from: { pathname: "/profile" } },
      },
      auth: { login },
    });

    fireEvent.change(screen.getByLabelText(/email address/i), {
      target: { value: "amina@example.com" },
    });
    fireEvent.change(screen.getByLabelText(/password/i), {
      target: { value: "secret1" },
    });
    fireEvent.click(screen.getByRole("button", { name: /login/i }));

    await waitFor(() => {
      expect(login).toHaveBeenCalledWith({
        email: "amina@example.com",
        password: "secret1",
      });
    });
    expect(await screen.findByText("Profile page")).toBeInTheDocument();
  });

  it("shows the authentication error and keeps the form available", async () => {
    const login = vi.fn().mockRejectedValue(new Error("Invalid credentials"));
    renderLogin({ auth: { login } });

    fireEvent.change(screen.getByLabelText(/email address/i), {
      target: { value: "amina@example.com" },
    });
    fireEvent.change(screen.getByLabelText(/password/i), {
      target: { value: "wrong-password" },
    });
    fireEvent.click(screen.getByRole("button", { name: /login/i }));

    expect(await screen.findByText("Invalid credentials")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /login/i })).toBeEnabled();
  });
});
