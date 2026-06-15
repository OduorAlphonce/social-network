// import { posts } from "../assets/posts-data.js";
// import { comments } from "../assets/comments-data.js";
import Post from "../components/Post.jsx";

function Home() {
  const posts = [];

  return (
    <div className="posts">
      {posts.map((post) => {
        return <Post key={post.id} />;
      })}
    </div>
  );
}

export default Home;
