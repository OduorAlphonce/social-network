import "../styles/comment.css";

const Comment = ({ comment }) => {
  return (
    <div classname="comment ">
      <img
        src={comment?.author?.avatar ? comment.author.avatar : avatar}
        alt="avatar"
        className="profile-photo"
      />
      <div>
        <strong>{comment?.name}</strong>
        <p>{comment?.content}</p>
      </div>
      <span>{comment.time}</span> <span>like</span> <span>reply</span>
    </div>
  );
};

export default Comment;
