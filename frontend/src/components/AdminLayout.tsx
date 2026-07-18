import { NavLink, Outlet, Link } from "react-router-dom";
import ThemeToggle from "./ThemeToggle";

const navItems = [
  { to: "/admin", label: "Roadmap", end: true },
  { to: "/admin/features", label: "Feature flags" },
  { to: "/admin/users", label: "Usuarios" },
  { to: "/admin/tests", label: "Tests" },
];

export default function AdminLayout() {
  return (
    <div className="min-h-screen flex flex-col md:flex-row bg-white dark:bg-black">
      <nav className="md:w-56 bg-white dark:bg-black border-b md:border-b-0 md:border-r border-slate-200 dark:border-neutral-800 p-4 flex md:flex-col gap-1 overflow-x-auto">
        <div className="mb-4 hidden md:block">
          <p className="font-semibold text-slate-900 dark:text-slate-50">rimu admin</p>
          <Link to="/app/habits" className="text-xs text-slate-500 dark:text-slate-400 underline">
            volver a la app
          </Link>
        </div>
        {navItems.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            end={item.end}
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
        <div className="mt-auto">
          <ThemeToggle />
        </div>
      </nav>
      <main className="flex-1 p-6 bg-white dark:bg-black">
        <Outlet />
      </main>
    </div>
  );
}
