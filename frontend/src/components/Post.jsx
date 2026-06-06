import { useState } from "react";
import "../styles/post.css";
import { Dislike, Like } from "./Reactions";
import avatar from "../assets/user.svg";

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
        <img src={avatar} alt="avatar" />
        <div className="post-bio">
          <h5>User Name</h5>
          <small>16:06</small>
        </div>
      </div>
      <div className="post-body">
        <img src="" alt="post-image" />
        <p>Post body</p>
      </div>
      <div className="reaction-count">
        <div><Like/> 45</div>
        <div>13 Comments</div>
      </div>
      <div className="post-footer">
        <div className="reaction-container">
          <Like like={like} isActive={likePost} />
          <Dislike dislike={dislike} isActive={dislikePost} />
        </div>
        <div>
          <p className="links">Comment</p>
        </div>
        <div>
          <p className="links">Share</p>
        </div>
      </div>
    </div>
  );
};

export default Post;
