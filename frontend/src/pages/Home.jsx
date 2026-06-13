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
      <div className="quick-links">
        <div>
          <strong>Sponsored</strong>
        </div>
        <div>
          <strong>Upcoming Events</strong>
        </div>
        <div>
          <strong>Contacts</strong>
        </div>
      </div>
    </div>
  );
}

export default Home;
