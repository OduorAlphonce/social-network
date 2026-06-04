import "../styles/post.css";
import { Dislike, Like } from "./Reactions";

const Post = ({ like, dislike }) => {
    
  return (
    <div className="post-container">
      <div className="post-header">
        <img src="" alt="avatar" />
        <h5>User Name</h5>
      </div>
      <div className="post-body">
        <img src="" alt="post-image" />
        <p>Post body</p>
      </div>
      <div className="post-footer">
        <Like like={like} />
        <Dislike dislike={dislike}/>
      </div>
    </div>
  );
};

export default Post;
