import { MdOutlineEvent, MdOutlineGroup } from "react-icons/md";
import { AiOutlineMessage } from "react-icons/ai";
import { IoIosInformationCircleOutline } from "react-icons/io";
import { BiGroup, BiHome, BiUser } from "react-icons/bi";

const Sidebar = () => {
  return (
    <aside className="sidebar">
      <li className="links">
        <BiHome /> Home
      </li>
      <li className="links">
        <BiUser /> Profile
      </li>
      <li className="links">
        <BiGroup />
        Friends
      </li>
      <li className="links">
        <MdOutlineGroup />
        Groups
      </li>
      <li className="links">
        <AiOutlineMessage />
        Messages
      </li>
      <li className="links">
        <IoIosInformationCircleOutline />
        Notification
      </li>
      <li className="links">
        <MdOutlineEvent />
        Events
      </li>
    </aside>
  );
};

export default Sidebar;
