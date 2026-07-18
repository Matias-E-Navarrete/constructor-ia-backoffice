import { Navigate, Route, Routes } from "react-router-dom";
import type { ReactNode } from "react";
import { useAuth } from "./context/AuthContext";
import AppLayout from "./components/AppLayout";
import AdminLayout from "./components/AdminLayout";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Tasks from "./pages/app/Tasks";
import Inbox from "./pages/app/Inbox";
import Planner from "./pages/app/Planner";
import Habits from "./pages/app/Habits";
import Panorama from "./pages/app/Panorama";
import Workouts from "./pages/app/Workouts";
import Studies from "./pages/app/Studies";
import Finance from "./pages/app/Finance";
import Family from "./pages/app/Family";
import Coaching from "./pages/app/Coaching";
import Notes from "./pages/app/Notes";
import NotesGraph from "./pages/app/NotesGraph";
import Upgrade from "./pages/app/Upgrade";
import AdminDashboard from "./pages/admin/Dashboard";
import AdminFeatures from "./pages/admin/Features";
import AdminUsers from "./pages/admin/Users";
import AdminTestRunner from "./pages/admin/TestRunner";

function RequireAuth({ children }: { children: ReactNode }) {
  const { user, loading } = useAuth();
  if (loading) return <div className="p-8 text-neutral-500 dark:text-neutral-400 bg-white dark:bg-black min-h-screen">Cargando...</div>;
  if (!user) return <Navigate to="/login" replace />;
  return <>{children}</>;
}

function RequireAdmin({ children }: { children: ReactNode }) {
  const { user, loading } = useAuth();
  if (loading) return <div className="p-8 text-neutral-500 dark:text-neutral-400 bg-white dark:bg-black min-h-screen">Cargando...</div>;
  if (!user) return <Navigate to="/login" replace />;
  if (user.role !== "admin") return <Navigate to="/app/tasks" replace />;
  return <>{children}</>;
}

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<Navigate to="/app/tasks" replace />} />
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />

      <Route
        path="/app"
        element={
          <RequireAuth>
            <AppLayout />
          </RequireAuth>
        }
      >
        <Route path="tasks" element={<Tasks />} />
        <Route path="inbox" element={<Inbox />} />
        <Route path="planner" element={<Planner />} />
        <Route path="habits" element={<Habits />} />
        <Route path="habits/panorama" element={<Panorama />} />
        <Route path="workouts" element={<Workouts />} />
        <Route path="studies" element={<Studies />} />
        <Route path="finance" element={<Finance />} />
        <Route path="notes" element={<Notes />} />
        <Route path="notes/graph" element={<NotesGraph />} />
        <Route path="family" element={<Family />} />
        <Route path="coaching" element={<Coaching />} />
        <Route path="upgrade" element={<Upgrade />} />
      </Route>

      <Route
        path="/admin"
        element={
          <RequireAdmin>
            <AdminLayout />
          </RequireAdmin>
        }
      >
        <Route index element={<AdminDashboard />} />
        <Route path="features" element={<AdminFeatures />} />
        <Route path="users" element={<AdminUsers />} />
        <Route path="tests" element={<AdminTestRunner />} />
      </Route>

      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
