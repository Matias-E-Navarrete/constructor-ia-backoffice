import { Navigate, Route, Routes } from "react-router-dom";
import type { ReactNode } from "react";
import { useAuth } from "./context/AuthContext";
import AppLayout from "./components/AppLayout";
import AdminLayout from "./components/AdminLayout";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Habits from "./pages/app/Habits";
import Workouts from "./pages/app/Workouts";
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
  if (loading) return <div className="p-8 text-slate-500">Cargando...</div>;
  if (!user) return <Navigate to="/login" replace />;
  return <>{children}</>;
}

function RequireAdmin({ children }: { children: ReactNode }) {
  const { user, loading } = useAuth();
  if (loading) return <div className="p-8 text-slate-500">Cargando...</div>;
  if (!user) return <Navigate to="/login" replace />;
  if (user.role !== "admin") return <Navigate to="/app/habits" replace />;
  return <>{children}</>;
}

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<Navigate to="/app/habits" replace />} />
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
        <Route path="habits" element={<Habits />} />
        <Route path="workouts" element={<Workouts />} />
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
