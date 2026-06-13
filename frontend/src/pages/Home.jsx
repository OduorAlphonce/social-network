import React from "react";
import { posts } from "../assets/posts-data.js";
import { comments } from "../assets/comments-data.js";
import Post from "../components/Post.jsx";
import "../styles/home.css";

function Home() {
  return (
    <div className="home-container">
    <div className="posts">
      {posts.map((it, idx) => {
        return <Post key={idx} post={it} />;
      })}
    </div>
    <div className="quick-links">Hello</div>
    </div>
  );
}

export default Home;
