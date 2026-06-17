import { useLocation } from "react-router";
import Post from "../components/Post.jsx";
import "../styles/post-detail.css";
import Comment from "../components/Comment.jsx";

const PostDetail = () => {
  const mockComment = [
    {
      id: "c101",
      name: "Sarah Jenkins",
      content:
        "This is a fantastic breakdown! Thanks for sharing this solution.",
      time: "12 minutes ago", // Matches your DateFormatter output scale
      author: {
        avatar: "https://i.pravatar.cc/150?img=32", // URL string or imported asset
      },
    },
  ];
  const location = useLocation();
  return (
    <div className="post-detail">
      <Post post={location.state} />;
      <div className="comments card">
        <h3>Comments</h3>
        {mockComment.map((it, idx) => (
          <Comment comment={it} key={idx} />
        ))}
      </div>
    </div>
  );
};

export default PostDetail;
