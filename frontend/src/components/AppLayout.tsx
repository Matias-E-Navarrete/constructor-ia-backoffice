import { NavLink, Outlet, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import ThemeToggle from "./ThemeToggle";

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
    <div className="min-h-screen flex flex-col md:flex-row bg-white dark:bg-black">
      <nav className="md:w-56 bg-white dark:bg-black border-b md:border-b-0 md:border-r border-slate-200 dark:border-neutral-800 p-4 flex md:flex-col gap-1 overflow-x-auto">
        <div className="mb-4 hidden md:block">
          <p className="font-semibold text-slate-900 dark:text-slate-50">rimu</p>
          <p className="text-xs text-slate-500 dark:text-slate-400">{user?.email}</p>
          <span
            className={`inline-block mt-1 text-xs px-2 py-0.5 rounded-full ${
              user?.plan === "pro"
                ? "bg-emerald-100 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-400"
                : "bg-slate-100 text-slate-600 dark:bg-neutral-800 dark:text-slate-400"
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
                isActive
                  ? "bg-slate-900 text-white dark:bg-white dark:text-slate-900"
                  : "text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-neutral-800"
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
              isActive
                ? "bg-emerald-600 text-white"
                : "text-emerald-700 dark:text-emerald-400 hover:bg-emerald-50 dark:hover:bg-emerald-950"
            }`
          }
        >
          {user?.plan === "pro" ? "Mi plan" : "Mejorar a Pro"}
        </NavLink>
        {user?.role === "admin" && (
          <NavLink
            to="/admin"
            className="whitespace-nowrap px-3 py-2 rounded-md text-sm font-medium text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-neutral-800"
          >
            Admin
          </NavLink>
        )}
        <div className="mt-auto flex flex-col gap-1">
          <ThemeToggle />
          <button
            onClick={handleLogout}
            className="whitespace-nowrap px-3 py-2 rounded-md text-sm font-medium text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-950 text-left"
          >
            Cerrar sesión
          </button>
        </div>
      </nav>
      <main className="flex-1 p-6 bg-white dark:bg-black">
        <Outlet />
      </main>
    </div>
  );
}
