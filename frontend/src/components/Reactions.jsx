import { BiSolidDislike, BiSolidLike } from 'react-icons/bi';

const Like = ({ like, isActive }) => {
  return (
    <BiSolidLike
      id="like"
      onClick={like}
      size={24}
      className={isActive ? 'reaction-like' : ''}
    />
  );
};

const Dislike = ({ dislike, isActive }) => {
  return (
    <BiSolidDislike
      id="dislike"
      onClick={dislike}
      size={24}
      className={isActive ? 'reaction-dislike' : ''}
    />
  );
};

export { Like, Dislike };
