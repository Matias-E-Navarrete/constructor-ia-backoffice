import { NavLink, Outlet, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import ThemeToggle from "./ThemeToggle";
import QuickAddFab from "./QuickAddFab";

const navItems = [
  { to: "/app/tasks", label: "Tareas" },
  { to: "/app/planner", label: "Planner" },
  { to: "/app/inbox", label: "Bandeja" },
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
      <nav className="md:w-56 bg-white dark:bg-black border-b md:border-b-0 md:border-r border-neutral-200 dark:border-neutral-800 p-4 flex md:flex-col gap-1 overflow-x-auto">
        <div className="mb-4 hidden md:block">
          <p className="font-semibold text-neutral-900 dark:text-neutral-50">rimu</p>
          <p className="text-xs text-neutral-500 dark:text-neutral-400">{user?.email}</p>
          <span
            className={`inline-block mt-1 text-xs px-2 py-0.5 rounded-full font-semibold tracking-wide ${
              user?.plan === "pro"
                ? "bg-accent text-white"
                : "bg-neutral-100 text-neutral-600 dark:bg-neutral-800 dark:text-neutral-400"
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
                  ? "bg-neutral-900 text-white dark:bg-white dark:text-neutral-900"
                  : "text-neutral-700 dark:text-neutral-300 hover:bg-neutral-100 dark:hover:bg-neutral-800"
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
                ? "bg-accent text-white"
                : "text-accent-strong dark:text-accent-soft hover:bg-fuchsia-50 dark:hover:bg-fuchsia-950/40"
            }`
          }
        >
          {user?.plan === "pro" ? "Mi plan" : "Mejorar a Pro"}
        </NavLink>
        {user?.role === "admin" && (
          <NavLink
            to="/admin"
            className="whitespace-nowrap px-3 py-2 rounded-md text-sm font-medium text-neutral-700 dark:text-neutral-300 hover:bg-neutral-100 dark:hover:bg-neutral-800"
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
      <QuickAddFab />
    </div>
  );
}
