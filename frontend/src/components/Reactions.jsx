import { BiDislike, BiLike } from "react-icons/bi";

const Like = ({ onClick }) => {
  return <BiLike onClick={onClick} />;
};

const Dislike = ({ onClick }) => {
  return <BiDislike onClick={onClick} />;
};

export { Like, Dislike };
