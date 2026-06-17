import { posts } from "../assets/posts-data.js";
import Post from "../components/Post.jsx";
import "../styles/home.css";

function Home() {
  const Allposts = posts;

  return (
    <div className="home-container">
      <div className="posts">
        {Allposts?.map((it, idx) => {
          return <Post key={idx} post={it} />;
        })}
      </div>
    </div>
  );
}

export default Home;
