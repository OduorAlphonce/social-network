import "../styles/comment.css";

const Comment = ({ comment }) => {

 
  return (
    <div id="comment-container">
      <img
        src={comment?.author?.avatar ? comment.author.avatar : avatar}
        alt="avatar"
        className="profile-photo"
      />
      <div className="comment-body">
        <div className="comment-details">
          <strong>{comment?.name}</strong>
          <p>{comment?.content}</p>
        </div>
        <div>{comment?.time} <span>Like</span>  <span>Reply</span></div>
      </div>
    </div>
  );
};

export default Comment;
