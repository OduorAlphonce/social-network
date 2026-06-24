import { useState } from "react";
import { posts } from "../assets/posts-data.js";
import Post from "../components/Post.jsx";
import "../styles/home.css";
import NewPost from "../components/NewPost.jsx";

function Home() {
  const initialPosts = posts || [];
  const [Allposts, setAllposts] = useState(initialPosts);

  function handleNewPost(created) {
    // send the new post to the backend and update the state with the new post
    setAllposts((prev) => [created, ...prev]);
    
  }

  return (
    <div className="home-container">
      <div className="posts">
        <NewPost onCreate={handleNewPost} />
        {Allposts?.map((it, idx) => {
          return <Post key={idx} post={it} />;
        })}
      </div>
    </div>
  );
}

export default Home;
