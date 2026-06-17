import { fireEvent, screen } from "@testing-library/react";
import { Route, Routes } from "react-router";
import { describe, expect, it } from "vitest";
import Post from "./Post";
import { renderWithProviders } from "../test/render";

const post = {
  id: "post-1",
  author: { name: "Amina Njeri" },
  content: "A careful post about testing behavior.",
  privacy: "public",
  created_at: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(),
  like_count: 2,
  comment_count: 3,
};

describe("Post", () => {
  it("opens the post when the post content is selected", async () => {
    renderWithProviders(
      <Routes>
        <Route path="/" element={<Post post={post} />} />
        <Route path="/post/:id" element={<p>Post details opened</p>} />
      </Routes>
    );

    fireEvent.click(screen.getByText(post.content));

    expect(await screen.findByText("Post details opened")).toBeInTheDocument();
  });

  it("toggles reactions without opening the post", () => {
    renderWithProviders(
      <Routes>
        <Route path="/" element={<Post post={post} />} />
        <Route path="/post/:id" element={<p>Post details opened</p>} />
      </Routes>
    );

    const like = screen.getByRole("button", { name: /^like$/i });
    const dislike = screen.getByRole("button", { name: /^dislike$/i });

    fireEvent.click(like);
    expect(like).toHaveAttribute("aria-pressed", "true");
    expect(dislike).toHaveAttribute("aria-pressed", "false");
    expect(screen.queryByText("Post details opened")).not.toBeInTheDocument();

    fireEvent.click(dislike);
    expect(like).toHaveAttribute("aria-pressed", "false");
    expect(dislike).toHaveAttribute("aria-pressed", "true");
    expect(screen.queryByText("Post details opened")).not.toBeInTheDocument();
  });
});
