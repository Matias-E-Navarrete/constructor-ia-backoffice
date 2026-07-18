import { NavLink, Outlet, Link } from "react-router-dom";

const navItems = [
  { to: "/admin", label: "Roadmap", end: true },
  { to: "/admin/features", label: "Feature flags" },
  { to: "/admin/users", label: "Usuarios" },
  { to: "/admin/tests", label: "Tests" },
];

export default function AdminLayout() {
  return (
    <div className="min-h-screen flex flex-col md:flex-row bg-slate-50">
      <nav className="md:w-56 bg-white border-b md:border-b-0 md:border-r border-slate-200 p-4 flex md:flex-col gap-1 overflow-x-auto">
        <div className="mb-4 hidden md:block">
          <p className="font-semibold text-slate-900">rimu admin</p>
          <Link to="/app/habits" className="text-xs text-slate-500 underline">
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
                isActive ? "bg-slate-900 text-white" : "text-slate-700 hover:bg-slate-100"
              }`
            }
          >
            {item.label}
          </NavLink>
        ))}
      </nav>
      <main className="flex-1 p-6">
        <Outlet />
      </main>
    </div>
  );
}
