import { useState } from "react";
import "../styles/post.css";
import { Dislike, Like } from "./Reactions";

const Post = () => {
  const [likePost, setLikePost] = useState(false);
  const [dislikePost, setDislikePost] = useState(false);

  const dislikeReaction = document.getElementById("dislike");
  const likeReaction = document.getElementById("like");

  function like() {
    if (dislikePost) {
      setDislikePost(false);
    }
    setLikePost((prev) => !prev);
    if (likePost) {
      likeReaction.classList.add("reaction-active");
      dislikeReaction.classList.remove("reaction-active");
    } else {
      likeReaction.classList.remove("reaction-active");
    }
  }

  function dislike() {
    if (likePost) {
      setLikePost(false);
    }
    setDislikePost((prev) => !prev);
    if (dislikePost) {
      dislikeReaction.classList.add("reaction-active");
      likeReaction.classList.remove("reaction-active");
    } else {
      dislikeReaction.classList.remove("reaction-active");
    }
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
        <Like like={like} />
        <Dislike dislike={dislike} />
      </div>
    </div>
  );
};

export default Post;
