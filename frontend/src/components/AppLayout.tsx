import { NavLink, Outlet, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

const navItems = [
  { to: "/app/habits", label: "Hábitos" },
  { to: "/app/workouts", label: "Entrenamientos" },
  { to: "/app/finance", label: "Finanzas" },
  { to: "/app/notes", label: "Notas" },
  { to: "/app/family", label: "Familia" },
  { to: "/app/coaching", label: "Coaching" },
];

export default function AppLayout() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  function handleLogout() {
    logout();
    navigate("/login");
  }

  return (
    <div className="min-h-screen flex flex-col md:flex-row bg-slate-50">
      <nav className="md:w-56 bg-white border-b md:border-b-0 md:border-r border-slate-200 p-4 flex md:flex-col gap-1 overflow-x-auto">
        <div className="mb-4 hidden md:block">
          <p className="font-semibold text-slate-900">rimu</p>
          <p className="text-xs text-slate-500">{user?.email}</p>
          <span
            className={`inline-block mt-1 text-xs px-2 py-0.5 rounded-full ${
              user?.plan === "pro" ? "bg-emerald-100 text-emerald-700" : "bg-slate-100 text-slate-600"
            }`}
          >
            {user?.plan === "pro" ? "PRO" : "FREE"}
          </span>
        </div>
        {navItems.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            className={({ isActive }) =>
              `whitespace-nowrap px-3 py-2 rounded-md text-sm font-medium ${
                isActive ? "bg-slate-900 text-white" : "text-slate-700 hover:bg-slate-100"
              }`
            }
          >
            {item.label}
          </NavLink>
        ))}
        <NavLink
          to="/app/upgrade"
          className={({ isActive }) =>
            `whitespace-nowrap px-3 py-2 rounded-md text-sm font-medium ${
              isActive ? "bg-emerald-600 text-white" : "text-emerald-700 hover:bg-emerald-50"
            }`
          }
        >
          {user?.plan === "pro" ? "Mi plan" : "Mejorar a Pro"}
        </NavLink>
        {user?.role === "admin" && (
          <NavLink
            to="/admin"
            className="whitespace-nowrap px-3 py-2 rounded-md text-sm font-medium text-slate-700 hover:bg-slate-100"
          >
            Admin
          </NavLink>
        )}
        <button
          onClick={handleLogout}
          className="mt-auto whitespace-nowrap px-3 py-2 rounded-md text-sm font-medium text-red-600 hover:bg-red-50 text-left"
        >
          Cerrar sesión
        </button>
      </nav>
      <main className="flex-1 p-6">
        <Outlet />
      </main>
    </div>
  );
}
