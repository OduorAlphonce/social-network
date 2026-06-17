import { BiSolidDislike, BiSolidLike } from "react-icons/bi";

const Like = ({ like, isActive }) => {
  if (!like) {
    return <BiSolidLike aria-hidden="true" size={24} />;
  }

  const handleClick = (event) => {
    event.stopPropagation();
    like?.();
  };

  return (
    <button
      type="button"
      aria-label="Like"
      aria-pressed={Boolean(isActive)}
      className={`reaction-button ${isActive ? "reaction-like" : ""}`}
      onClick={handleClick}
    >
      <BiSolidLike aria-hidden="true" size={24} />
    </button>
  );
};

const Dislike = ({ dislike, isActive }) => {
  if (!dislike) {
    return <BiSolidDislike aria-hidden="true" size={24} />;
  }

  const handleClick = (event) => {
    event.stopPropagation();
    dislike?.();
  };

  return (
    <button
      type="button"
      aria-label="Dislike"
      aria-pressed={Boolean(isActive)}
      className={`reaction-button ${isActive ? "reaction-dislike" : ""}`}
      onClick={handleClick}
    >
      <BiSolidDislike aria-hidden="true" size={24} />
    </button>
  );
};

export { Like, Dislike };
