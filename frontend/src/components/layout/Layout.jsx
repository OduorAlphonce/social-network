import { Outlet } from 'react-router';
import '../../styles/layout.css';
import Header from '../Header';
import Sidebar from '../Sidebar';

const Layout = () => {
  return (
    <div className="main-container">
      <Sidebar />
      <Header />
      <div className="main-content">
        <Outlet />
      </div>
    </div>
  );
};

export default Layout;
