import { fireEvent, screen } from "@testing-library/react";
import { Route, Routes } from "react-router";
import { describe, expect, it, vi } from "vitest";
import LogoutButton from "./LogoutButton";
import { renderWithProviders } from "../test/render";

describe("LogoutButton", () => {
  it("logs out and navigates to login", async () => {
    const logout = vi.fn().mockResolvedValue(undefined);
    renderWithProviders(
      <Routes>
        <Route path="/" element={<LogoutButton />} />
        <Route path="/login" element={<p>Login screen</p>} />
      </Routes>,
      { auth: { logout } }
    );

    fireEvent.click(screen.getByRole("button", { name: /logout/i }));

    expect(await screen.findByText("Login screen")).toBeInTheDocument();
    expect(logout).toHaveBeenCalledTimes(1);
  });
});
