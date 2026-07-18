import { useEffect, useState } from "react";
import * as adminApi from "../../api/admin";

export default function AdminUsers() {
  const [users, setUsers] = useState<adminApi.AdminUser[]>([]);

  async function load() {
    setUsers(await adminApi.listUsers());
  }

  useEffect(() => {
    load();
  }, []);

  async function handlePlan(id: string, plan: "free" | "pro") {
    await adminApi.setUserPlan(id, plan);
    load();
  }

  async function handleRole(id: string, role: "user" | "admin") {
    await adminApi.setUserRole(id, role);
    load();
  }

  return (
    <div>
      <h1 className="text-xl font-semibold text-neutral-900 dark:text-neutral-50 mb-4">Usuarios</h1>
      <table className="w-full bg-white dark:bg-neutral-900 border border-neutral-200 dark:border-neutral-800 rounded-lg overflow-hidden text-sm">
        <thead className="bg-neutral-50 dark:bg-neutral-800 text-neutral-500 dark:text-neutral-400 text-left">
          <tr>
            <th className="p-3">Email</th>
            <th className="p-3">Plan</th>
            <th className="p-3">Rol</th>
          </tr>
        </thead>
        <tbody>
          {users.map((u) => (
            <tr key={u.id} className="border-t border-neutral-100 dark:border-neutral-800 text-neutral-800 dark:text-neutral-200">
              <td className="p-3">{u.email}</td>
              <td className="p-3">
                <select value={u.plan} onChange={(e) => handlePlan(u.id, e.target.value as "free" | "pro")} className="rounded border border-neutral-300 dark:border-neutral-700 bg-white dark:bg-neutral-900 text-neutral-900 dark:text-neutral-100 px-2 py-1">
                  <option value="free">free</option>
                  <option value="pro">pro</option>
                </select>
              </td>
              <td className="p-3">
                <select value={u.role} onChange={(e) => handleRole(u.id, e.target.value as "user" | "admin")} className="rounded border border-neutral-300 dark:border-neutral-700 bg-white dark:bg-neutral-900 text-neutral-900 dark:text-neutral-100 px-2 py-1">
                  <option value="user">user</option>
                  <option value="admin">admin</option>
                </select>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
