import { fireEvent, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import RegisterForm from "./RegisterForm";
import { renderWithProviders } from "../test/render";
import { apiFetch } from "../utils/api";

vi.mock("../utils/api", async (importOriginal) => {
  const original = await importOriginal();

  return {
    ...original,
    apiFetch: vi.fn(),
  };
});

describe("RegisterForm", () => {
  beforeEach(() => {
    apiFetch.mockReset();
  });

  it("submits registration fields as multipart data and resets after success", async () => {
    apiFetch.mockResolvedValueOnce({ id: "user-1" });
    renderWithProviders(<RegisterForm />);

    fireEvent.change(screen.getByLabelText(/email/i), {
      target: { value: "amina@example.com" },
    });
    fireEvent.change(screen.getByLabelText(/password/i), {
      target: { value: "secret1" },
    });
    fireEvent.change(screen.getByLabelText(/first name/i), {
      target: { value: "Amina" },
    });
    fireEvent.change(screen.getByLabelText(/last name/i), {
      target: { value: "Njeri" },
    });
    fireEvent.change(screen.getByLabelText(/date of birth/i), {
      target: { value: "1998-04-12" },
    });
    fireEvent.change(screen.getByLabelText(/nickname/i), {
      target: { value: "amina" },
    });
    fireEvent.change(screen.getByLabelText(/about me/i), {
      target: { value: "Weekend hiker" },
    });
    fireEvent.click(screen.getByLabelText(/public profile/i));
    const avatar = new File(["avatar"], "avatar.png", { type: "image/png" });
    fireEvent.change(screen.getByLabelText(/avatar/i), {
      target: { files: [avatar] },
    });
    fireEvent.click(screen.getByRole("button", { name: /register/i }));

    await waitFor(() => {
      expect(apiFetch).toHaveBeenCalledWith(
        "/api/users/register",
        expect.objectContaining({ method: "POST", body: expect.any(FormData) })
      );
    });
    const formData = apiFetch.mock.calls[0][1].body;
    expect(formData.get("email")).toBe("amina@example.com");
    expect(formData.get("password")).toBe("secret1");
    expect(formData.get("first_name")).toBe("Amina");
    expect(formData.get("last_name")).toBe("Njeri");
    expect(formData.get("date_of_birth")).toBe("1998-04-12");
    expect(formData.get("nickname")).toBe("amina");
    expect(formData.get("about_me")).toBe("Weekend hiker");
    expect(formData.get("is_public")).toBe("true");
    expect(formData.get("avatar")).toBe(avatar);
    expect(
      await screen.findByText(/registration successful/i)
    ).toBeInTheDocument();
    expect(screen.getByLabelText(/email/i)).toHaveValue("");
  });

  it("shows the API error without clearing the entered fields", async () => {
    apiFetch.mockRejectedValueOnce(new Error("email already registered"));
    renderWithProviders(<RegisterForm />);

    fireEvent.change(screen.getByLabelText(/email/i), {
      target: { value: "amina@example.com" },
    });
    fireEvent.change(screen.getByLabelText(/password/i), {
      target: { value: "secret1" },
    });
    fireEvent.change(screen.getByLabelText(/first name/i), {
      target: { value: "Amina" },
    });
    fireEvent.change(screen.getByLabelText(/last name/i), {
      target: { value: "Njeri" },
    });
    fireEvent.change(screen.getByLabelText(/date of birth/i), {
      target: { value: "1998-04-12" },
    });
    fireEvent.click(screen.getByRole("button", { name: /register/i }));

    expect(
      await screen.findByText("email already registered")
    ).toBeInTheDocument();
    expect(screen.getByLabelText(/email/i)).toHaveValue("amina@example.com");
  });
});
