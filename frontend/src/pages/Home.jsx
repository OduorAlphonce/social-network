import React from "react";
// import { posts } from "../assets/posts-data.js";
// import { comments } from "../assets/comments-data.js";
import Post from "../components/Post.jsx";

function Home() {
  const posts = new Array()

  return (
    <div className="posts">
      {posts.map((it) => {
        return <Post />;
      })}
    </div>
  );
}

export default Home;
