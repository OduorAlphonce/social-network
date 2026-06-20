import { fireEvent, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import ProfileUpdateForm from "./ProfileUpdateForm";
import { renderWithProviders } from "../test/render";
import { apiFetch } from "../utils/api";

vi.mock("../utils/api", async (importOriginal) => {
  const original = await importOriginal();

  return {
    ...original,
    apiFetch: vi.fn(),
  };
});

const currentUser = {
  id: "user-1",
  email: "amina@example.com",
  first_name: "Amina",
  last_name: "Njeri",
  date_of_birth: "1998-04-12",
  nickname: "amina",
  about_me: "Weekend hiker",
  is_public: true,
};

describe("ProfileUpdateForm", () => {
  beforeEach(() => {
    apiFetch.mockReset();
  });

  it("preloads the current profile and submits changed fields", async () => {
    const refresh = vi.fn().mockResolvedValue(currentUser);
    apiFetch.mockResolvedValueOnce({ id: "user-1" });
    renderWithProviders(<ProfileUpdateForm />, {
      auth: { currentUser, isAuthenticated: true, refresh },
    });

    expect(screen.getByLabelText(/email/i)).toHaveValue("amina@example.com");
    expect(screen.getByLabelText(/first name/i)).toHaveValue("Amina");
    fireEvent.change(screen.getByLabelText(/first name/i), {
      target: { value: "Grace" },
    });
    fireEvent.change(screen.getByLabelText(/new password/i), {
      target: { value: "secret2" },
    });
    fireEvent.click(screen.getByRole("button", { name: /update profile/i }));

    await waitFor(() => {
      expect(apiFetch).toHaveBeenCalledWith(
        "/api/users/update",
        expect.objectContaining({ method: "PATCH", body: expect.any(FormData) })
      );
    });
    const formData = apiFetch.mock.calls[0][1].body;
    expect(formData.get("email")).toBe("amina@example.com");
    expect(formData.get("first_name")).toBe("Grace");
    expect(formData.get("new_password")).toBe("secret2");
    expect(formData.has("current_password")).toBe(false);
    expect(formData.get("is_public")).toBe("true");
    expect(refresh).toHaveBeenCalledTimes(1);
    expect(
      await screen.findByText(/profile updated successfully/i)
    ).toBeInTheDocument();
    expect(screen.getByLabelText(/new password/i)).toHaveValue("");
  });

  it("shows loading while the current profile is absent", () => {
    renderWithProviders(<ProfileUpdateForm />, {
      auth: { currentUser: null, isAuthenticated: true },
    });

    expect(screen.getByText(/loading/i)).toBeInTheDocument();
  });

  it("shows API errors without refreshing the profile", async () => {
    const refresh = vi.fn();
    apiFetch.mockRejectedValueOnce(new Error("invalid current password"));
    renderWithProviders(<ProfileUpdateForm />, {
      auth: { currentUser, isAuthenticated: true, refresh },
    });

    fireEvent.change(screen.getByLabelText(/email/i), {
      target: { value: "new@example.com" },
    });
    fireEvent.click(screen.getByRole("button", { name: /update profile/i }));

    expect(
      await screen.findByText("invalid current password")
    ).toBeInTheDocument();
    expect(refresh).not.toHaveBeenCalled();
  });
});
