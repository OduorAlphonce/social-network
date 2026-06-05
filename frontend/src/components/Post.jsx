import { useState } from "react";
import "../styles/post.css";
import { Dislike, Like } from "./Reactions";

const Post = () => {
  const [likePost, setLikePost] = useState(false);
  const [dislikePost, setDislikePost] = useState(false);

  function like() {
    setDislikePost(false);
    setLikePost((prev) => !prev);
  }

  function dislike() {
    setLikePost(false);
    setDislikePost((prev) => !prev);
  }

  return (
    <div className="post-container">
      <div className="post-header">
        <img src="" alt="avatar" />
        <div className="post-bio">
          <h5>User Name</h5>
          <small>16:06</small>
        </div>
      </div>
      <div className="post-body">
        <img src="" alt="post-image" />
        <p>Post body</p>
      </div>
      <div className="post-footer">
        <Like like={like} isActive={likePost} />
        <Dislike dislike={dislike} isActive={dislikePost} />
      </div>
    </div>
  );
};

export default Post;
