import { NavLink, Outlet, Link } from "react-router-dom";
import ThemeToggle from "./ThemeToggle";
import { IconRoadmap, IconFlag, IconUsersGear, IconTests, IconArrowLeft } from "./icons";

const navItems = [
  { to: "/admin", label: "Roadmap", end: true, icon: IconRoadmap },
  { to: "/admin/features", label: "Feature flags", icon: IconFlag },
  { to: "/admin/users", label: "Usuarios", icon: IconUsersGear },
  { to: "/admin/tests", label: "Tests", icon: IconTests },
];

export default function AdminLayout() {
  return (
    <div className="min-h-screen flex flex-col md:flex-row bg-white dark:bg-black">
      <nav className="md:w-60 md:h-screen md:sticky md:top-0 bg-neutral-50/60 dark:bg-neutral-950/60 border-b md:border-b-0 md:border-r border-neutral-200 dark:border-neutral-800 p-4 flex md:flex-col gap-1 overflow-x-auto">
        <div className="mb-5 hidden md:flex items-center gap-2.5">
          <div className="w-8 h-8 rounded-lg bg-accent flex items-center justify-center flex-shrink-0 shadow-glow-sm">
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
        <div className="flex md:flex-col gap-1">
          {navItems.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              end={item.end}
              className={({ isActive }) =>
                `flex items-center gap-2.5 whitespace-nowrap px-3 py-2 rounded-lg text-sm font-medium transition-colors duration-150 ${
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
        <div className="mt-auto pt-2 md:border-t border-neutral-200 dark:border-neutral-800">
          <ThemeToggle />
        </div>
      </nav>
      <main className="flex-1 p-6 md:p-8 bg-white dark:bg-black">
        <Outlet />
      </main>
    </div>
  );
}
