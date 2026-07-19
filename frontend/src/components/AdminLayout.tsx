import { useState } from "react";
import { NavLink, Outlet, Link } from "react-router-dom";
import ThemeToggle from "./ThemeToggle";
import { IconRoadmap, IconFlag, IconUsersGear, IconTests, IconArrowLeft, IconMenu, IconClose } from "./icons";

const navItems = [
  { to: "/admin", label: "Roadmap", end: true, icon: IconRoadmap },
  { to: "/admin/features", label: "Feature flags", icon: IconFlag },
  { to: "/admin/users", label: "Usuarios", icon: IconUsersGear },
  { to: "/admin/tests", label: "Tests", icon: IconTests },
];

function SidebarContent({ onNavigate }: { onNavigate?: () => void }) {
  return (
    <>
      <div className="mb-5 flex items-center gap-2.5">
        <div className="w-8 h-8 rounded-lg bg-accent flex items-center justify-center flex-shrink-0 shadow-neon-sm">
          <span className="text-white font-bold text-sm">r</span>
        </div>
        <div className="min-w-0">
          <p className="font-semibold text-neutral-900 dark:text-neutral-50 leading-tight">admin</p>
          <Link
            to="/app/tasks"
            className="flex items-center gap-1 text-xs text-neutral-500 dark:text-neutral-400 hover:text-neutral-700 dark:hover:text-neutral-300 transition-colors duration-150"
          >
            <IconArrowLeft width={12} height={12} />
            volver a la app
          </Link>
        </div>
      </div>
      <div className="flex flex-col gap-1">
        {navItems.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            end={item.end}
            onClick={onNavigate}
            className={({ isActive }) =>
              `flex items-center gap-2.5 whitespace-nowrap px-3 py-2 rounded-lg text-sm font-medium transition-all duration-150 ${
                isActive
                  ? "bg-neutral-900 text-white dark:bg-white dark:text-neutral-900 shadow-neon-sm"
                  : "text-neutral-600 dark:text-neutral-400 hover:bg-neutral-200/60 dark:hover:bg-neutral-800/60 hover:text-neutral-900 dark:hover:text-neutral-100"
              }`
            }
          >
            <item.icon className="flex-shrink-0" />
            {item.label}
          </NavLink>
        ))}
      </div>
      <div className="mt-auto pt-2 border-t border-neutral-200 dark:border-neutral-800">
        <ThemeToggle />
      </div>
    </>
  );
}

export default function AdminLayout() {
  const [mobileOpen, setMobileOpen] = useState(false);

  return (
    <div className="min-h-screen flex flex-col md:flex-row bg-white dark:bg-black">
      <header className="md:hidden sticky top-0 z-30 flex items-center justify-between px-4 py-3 bg-white/95 dark:bg-black/95 backdrop-blur border-b border-neutral-200 dark:border-neutral-800">
        <div className="flex items-center gap-2.5">
          <div className="w-7 h-7 rounded-lg bg-accent flex items-center justify-center flex-shrink-0 shadow-neon-sm">
            <span className="text-white font-bold text-xs">r</span>
          </div>
          <span className="font-semibold text-neutral-900 dark:text-neutral-50">admin</span>
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
    </div>
  );
}
