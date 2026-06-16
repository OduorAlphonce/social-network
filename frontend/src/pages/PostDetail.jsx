import { useLocation } from "react-router";
import Post from "../components/Post.jsx";
import "../styles/post-detail.css";

const PostDetail = () => {
  const location = useLocation();
  return (
    <div className="post-detail">
      <Post post={location.state} />;
    </div>
  );
};

export default PostDetail;
