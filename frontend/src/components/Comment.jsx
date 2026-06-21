import "../styles/comment.css";
import avatar from "../assets/user.svg";
import { Like } from "./Reactions";

const Comment = ({ comment }) => {
  const authorName = comment?.author
    ? (comment.author.nickname || `${comment.author.first_name || ""} ${comment.author.last_name || ""}`.trim())
    : comment?.name;

  return (
    <div id="comment-container">
      <img
        src={comment?.author?.avatar ? comment.author.avatar : avatar}
        alt="avatar"
        className="profile-photo"
      />
      <div className="comment-body">
        <div className="comment-details">
          <strong>{authorName || "Anonymous"}</strong>
          <p>{comment?.content}</p>
        </div>
        <div className="comment-footer">
          {comment?.time} 
          <div className="comment-reaction">
          <span className="comment-like"><Like/></span> <span className="comment-reply">Reply</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Comment;
