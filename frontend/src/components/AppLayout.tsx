import { NavLink, Outlet, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { useTheme } from "../context/ThemeContext";
import QuickAddFab from "./QuickAddFab";
import {
  IconTasks,
  IconPlanner,
  IconInbox,
  IconHabits,
  IconWorkouts,
  IconStudies,
  IconFinance,
  IconNotes,
  IconFamily,
  IconCoaching,
  IconAdmin,
  IconUpgrade,
  IconSun,
  IconMoon,
  IconLogout,
} from "./icons";

const navItems = [
  { to: "/app/tasks", label: "Tareas", icon: IconTasks },
  { to: "/app/planner", label: "Planner", icon: IconPlanner },
  { to: "/app/inbox", label: "Bandeja", icon: IconInbox },
  { to: "/app/habits", label: "Hábitos", icon: IconHabits },
  { to: "/app/workouts", label: "Entrenamientos", icon: IconWorkouts },
  { to: "/app/studies", label: "Estudios", icon: IconStudies },
  { to: "/app/finance", label: "Finanzas", icon: IconFinance },
  { to: "/app/notes", label: "Notas", icon: IconNotes },
  { to: "/app/family", label: "Familia", icon: IconFamily },
  { to: "/app/coaching", label: "Coaching", icon: IconCoaching },
];

export default function AppLayout() {
  const { user, logout } = useAuth();
  const { theme, toggleTheme } = useTheme();
  const navigate = useNavigate();
  const isDark = theme === "dark";

  function handleLogout() {
    logout();
    navigate("/login");
  }

  return (
    <div className="min-h-screen flex flex-col md:flex-row bg-white dark:bg-black">
      <nav className="md:w-60 md:h-screen md:sticky md:top-0 bg-neutral-50/60 dark:bg-neutral-950/60 border-b md:border-b-0 md:border-r border-neutral-200 dark:border-neutral-800 p-4 flex md:flex-col gap-1 overflow-x-auto">
        <div className="mb-5 hidden md:flex items-center gap-2.5">
          <div className="w-8 h-8 rounded-lg bg-neutral-900 dark:bg-white flex items-center justify-center flex-shrink-0 shadow-glow-sm">
            <span className="text-white dark:text-neutral-900 font-bold text-sm">r</span>
          </div>
          <div className="min-w-0">
            <p className="font-semibold text-neutral-900 dark:text-neutral-50 leading-tight truncate">rimu</p>
            <p className="text-xs text-neutral-500 dark:text-neutral-400 truncate">{user?.email}</p>
          </div>
        </div>

        <div className="mb-3 hidden md:block">
          <span
            className={`inline-flex items-center text-xs px-2 py-0.5 rounded-full font-semibold tracking-wide ${
              user?.plan === "pro"
                ? "bg-accent text-white shadow-glow-sm"
                : "bg-neutral-100 text-neutral-600 dark:bg-neutral-800 dark:text-neutral-400"
            }`}
          >
            {user?.plan === "pro" ? "PRO" : "FREE"}
          </span>
        </div>

        <div className="flex md:flex-col gap-1">
          {navItems.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              className={({ isActive }) =>
                `group flex items-center gap-2.5 whitespace-nowrap px-3 py-2 rounded-lg text-sm font-medium transition-colors duration-150 ${
                  isActive
                    ? "bg-neutral-900 text-white dark:bg-white dark:text-neutral-900"
                    : "text-neutral-600 dark:text-neutral-400 hover:bg-neutral-200/60 dark:hover:bg-neutral-800/60 hover:text-neutral-900 dark:hover:text-neutral-100"
                }`
              }
            >
              <item.icon className="flex-shrink-0" />
              {item.label}
            </NavLink>
          ))}
        </div>

        <NavLink
          to="/app/upgrade"
          className={({ isActive }) =>
            `mt-1 flex items-center gap-2.5 whitespace-nowrap px-3 py-2 rounded-lg text-sm font-medium transition-colors duration-150 ${
              isActive
                ? "bg-accent text-white shadow-glow-sm"
                : "text-accent-strong dark:text-accent-soft hover:bg-fuchsia-50 dark:hover:bg-fuchsia-950/40"
            }`
          }
        >
          <IconUpgrade className="flex-shrink-0" />
          {user?.plan === "pro" ? "Mi plan" : "Mejorar a Pro"}
        </NavLink>

        {user?.role === "admin" && (
          <NavLink
            to="/admin"
            className="flex items-center gap-2.5 whitespace-nowrap px-3 py-2 rounded-lg text-sm font-medium text-neutral-600 dark:text-neutral-400 hover:bg-neutral-200/60 dark:hover:bg-neutral-800/60 hover:text-neutral-900 dark:hover:text-neutral-100 transition-colors duration-150"
          >
            <IconAdmin className="flex-shrink-0" />
            Admin
          </NavLink>
        )}

        <div className="mt-auto flex flex-col gap-1 pt-2 md:border-t border-neutral-200 dark:border-neutral-800">
          <button
            onClick={toggleTheme}
            aria-label={isDark ? "Cambiar a modo claro" : "Cambiar a modo oscuro"}
            className="w-full flex items-center gap-2.5 px-3 py-2 rounded-lg text-sm font-medium text-neutral-600 dark:text-neutral-400 hover:bg-neutral-200/60 dark:hover:bg-neutral-800/60 hover:text-neutral-900 dark:hover:text-neutral-100 transition-colors duration-150"
          >
            {isDark ? <IconMoon className="flex-shrink-0" /> : <IconSun className="flex-shrink-0" />}
            <span className="hidden md:inline">{isDark ? "Modo oscuro" : "Modo claro"}</span>
          </button>
          <button
            onClick={handleLogout}
            className="w-full flex items-center gap-2.5 px-3 py-2 rounded-lg text-sm font-medium text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-950/40 transition-colors duration-150 text-left"
          >
            <IconLogout className="flex-shrink-0" />
            <span className="hidden md:inline">Cerrar sesión</span>
          </button>
        </div>
      </nav>
      <main className="flex-1 p-6 md:p-8 bg-white dark:bg-black">
        <Outlet />
      </main>
      <QuickAddFab />
    </div>
  );
}
