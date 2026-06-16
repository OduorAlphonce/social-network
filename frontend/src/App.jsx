import { Route, Routes } from "react-router";
import "./App.css";
import Layout from "./components/layout/Layout";
import Events from "./pages/Events.jsx";
import Groups from "./pages/Groups.jsx";
import Home from "./pages/Home.jsx";
import Friends from "./pages/Friends.jsx";
import Profile from "./pages/Profile.jsx";
import Messages from "./pages/Messages.jsx";
import Notifications from "./pages/Notifications.jsx";
import RegisterForm from "./components/RegisterForm";
import LoginForm from "./components/LoginForm";

function App() {
  return (
    <Routes>
      <Route path="/register" element={<RegisterForm />} />
      <Route path="/login" element={<LoginForm />} />
      <Route path="/" element={<Layout />}>
        <Route index element={<Home />} />
        <Route path="/profile" element={<Profile />} />
        <Route path="/friends" element={<Friends />} />
        <Route path="/events" element={<Events />} />
        <Route path="/groups" element={<Groups />} />
        <Route path="/messages" element={<Messages />} />
        <Route path="/notifications" element={<Notifications />} />
      </Route>
    </Routes>
  );
}

export default App;
