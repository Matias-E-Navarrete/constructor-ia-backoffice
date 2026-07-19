import { useState } from "react";
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
  IconAssistant,
  IconSun,
  IconMoon,
  IconLogout,
  IconMenu,
  IconClose,
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
  { to: "/app/assistant", label: "Asistente", icon: IconAssistant },
];

function SidebarContent({ onNavigate }: { onNavigate?: () => void }) {
  const { user, logout } = useAuth();
  const { theme, toggleTheme } = useTheme();
  const navigate = useNavigate();
  const isDark = theme === "dark";

  function handleLogout() {
    logout();
    navigate("/login");
  }

  return (
    <>
      <div className="mb-5 flex items-center gap-2.5">
        <div className="w-8 h-8 rounded-lg bg-neutral-900 dark:bg-white flex items-center justify-center flex-shrink-0 shadow-glow-sm">
          <span className="text-white dark:text-neutral-900 font-bold text-sm">r</span>
        </div>
        <div className="min-w-0">
          <p className="font-semibold text-neutral-900 dark:text-neutral-50 leading-tight truncate">rimu</p>
          <p className="text-xs text-neutral-500 dark:text-neutral-400 truncate">{user?.email}</p>
        </div>
      </div>

      <div className="mb-3">
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

      <div className="flex flex-col gap-1">
        {navItems.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            onClick={onNavigate}
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
        onClick={onNavigate}
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
          onClick={onNavigate}
          className="flex items-center gap-2.5 whitespace-nowrap px-3 py-2 rounded-lg text-sm font-medium text-neutral-600 dark:text-neutral-400 hover:bg-neutral-200/60 dark:hover:bg-neutral-800/60 hover:text-neutral-900 dark:hover:text-neutral-100 transition-colors duration-150"
        >
          <IconAdmin className="flex-shrink-0" />
          Admin
        </NavLink>
      )}

      <div className="mt-auto flex flex-col gap-1 pt-2 border-t border-neutral-200 dark:border-neutral-800">
        <button
          onClick={toggleTheme}
          aria-label={isDark ? "Cambiar a modo claro" : "Cambiar a modo oscuro"}
          className="w-full flex items-center gap-2.5 px-3 py-2 rounded-lg text-sm font-medium text-neutral-600 dark:text-neutral-400 hover:bg-neutral-200/60 dark:hover:bg-neutral-800/60 hover:text-neutral-900 dark:hover:text-neutral-100 transition-colors duration-150"
        >
          {isDark ? <IconMoon className="flex-shrink-0" /> : <IconSun className="flex-shrink-0" />}
          {isDark ? "Modo oscuro" : "Modo claro"}
        </button>
        <button
          onClick={handleLogout}
          className="w-full flex items-center gap-2.5 px-3 py-2 rounded-lg text-sm font-medium text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-950/40 transition-colors duration-150 text-left"
        >
          <IconLogout className="flex-shrink-0" />
          Cerrar sesión
        </button>
      </div>
    </>
  );
}

export default function AppLayout() {
  const [mobileOpen, setMobileOpen] = useState(false);

  return (
    <div className="min-h-screen flex flex-col md:flex-row bg-white dark:bg-black">
      <header className="md:hidden sticky top-0 z-30 flex items-center justify-between px-4 py-3 bg-white/95 dark:bg-black/95 backdrop-blur border-b border-neutral-200 dark:border-neutral-800">
        <div className="flex items-center gap-2.5">
          <div className="w-7 h-7 rounded-lg bg-neutral-900 dark:bg-white flex items-center justify-center flex-shrink-0">
            <span className="text-white dark:text-neutral-900 font-bold text-xs">r</span>
          </div>
          <span className="font-semibold text-neutral-900 dark:text-neutral-50">rimu</span>
        </div>
        <button
          onClick={() => setMobileOpen(true)}
          aria-label="Abrir menú"
          className="w-9 h-9 flex items-center justify-center rounded-lg text-neutral-600 dark:text-neutral-400 hover:bg-neutral-100 dark:hover:bg-neutral-900 transition-colors duration-150"
        >
          <IconMenu />
        </button>
      </header>

      {mobileOpen && (
        <div className="md:hidden fixed inset-0 z-40 flex">
          <div className="absolute inset-0 bg-black/40" onClick={() => setMobileOpen(false)} />
          <nav className="relative w-72 max-w-[85vw] h-full bg-white dark:bg-black border-r border-neutral-200 dark:border-neutral-800 p-4 flex flex-col gap-1 overflow-y-auto">
            <button
              onClick={() => setMobileOpen(false)}
              aria-label="Cerrar menú"
              className="absolute top-4 right-4 w-8 h-8 flex items-center justify-center rounded-lg text-neutral-500 dark:text-neutral-400 hover:bg-neutral-100 dark:hover:bg-neutral-900 transition-colors duration-150"
            >
              <IconClose />
            </button>
            <SidebarContent onNavigate={() => setMobileOpen(false)} />
          </nav>
        </div>
      )}

      <nav className="hidden md:flex md:w-60 md:h-screen md:sticky md:top-0 bg-neutral-50/60 dark:bg-neutral-950/60 border-r border-neutral-200 dark:border-neutral-800 p-4 md:flex-col gap-1">
        <SidebarContent />
      </nav>

      <main className="flex-1 p-4 md:p-8 bg-white dark:bg-black">
        <Outlet />
      </main>
      <QuickAddFab />
    </div>
  );
}
