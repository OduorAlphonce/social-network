import { IoArrowRedo } from "react-icons/io5";
import "../styles/header.css";
import { FaPlusCircle } from "react-icons/fa";
import { AiOutlineMessage } from "react-icons/ai";
import { IoIosNotificationsOutline } from "react-icons/io";
import { MdOutlineGroup } from "react-icons/md";
import avatar from "../assets/user.svg";
import LogoutButton from "./LogoutButton";

const Header = () => {
  return (
    <div className="header">
      <div className="pretitle">
        <IoArrowRedo />
        <strong>CoreConnect</strong>
        <input
          type="text"
          className="search"
          placeholder="Search CoreConnect..."
        />
      </div>
      <div className="icons">
        <FaPlusCircle size={24} />
        <AiOutlineMessage size={24} />
        <IoIosNotificationsOutline size={24} />
        <MdOutlineGroup size={24} />
        <img src={avatar} className="profile-photo" />
        <LogoutButton />
      </div>
    </div>
  );
};

export default Header;
