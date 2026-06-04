import { BiDislike, BiLike } from "react-icons/bi";

const Like = ({ like }) => {
  return <BiLike onClick={like} size={24} />;
};

const Dislike = ({ dislike }) => {
  return <BiDislike onClick={dislike} size={24}/>;
};

export { Like, Dislike };
