import { BiGroup, BiHome, BiMessage, BiUser } from "react-icons/bi";
import "../../styles/layout.css"
import { MdGroupWork, MdOutlineEvent, MdOutlineGroup } from "react-icons/md";
import { AiOutlineMessage } from "react-icons/ai";
import { IoIosInformationCircleOutline } from "react-icons/io";

const Layout = () => {
  return (
    <div className="main-container">
      <aside className="sidebar">
        <li><BiHome/>Home</li>
        <li><BiUser/> Profile</li>
        <li><BiGroup/>Friends</li>
        <li><MdOutlineGroup/>Groups</li>
        <li><AiOutlineMessage/>Messages</li>
        <li><IoIosInformationCircleOutline/>Notification</li>
        <li><MdOutlineEvent/>Events</li>
      </aside>
      <div className="header">Header</div>
      <div className="main-content">Body</div>
    </div>
  );
};

export default Layout;
