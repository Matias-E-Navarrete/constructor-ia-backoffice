import { useEffect, useState } from "react";
import * as groupsApi from "../../api/groups";

export default function Family() {
  const [memberships, setMemberships] = useState<groupsApi.Membership[]>([]);
  const [groupName, setGroupName] = useState("");
  const [selected, setSelected] = useState<{ group: groupsApi.Group; members: groupsApi.Member[] } | null>(null);
  const [inviteEmail, setInviteEmail] = useState("");

  async function load() {
    const all = await groupsApi.listMyGroups();
    setMemberships(all.filter((m) => m.group.kind === "family"));
  }

  useEffect(() => {
    load();
  }, []);

  async function handleCreate() {
    if (!groupName.trim()) return;
    await groupsApi.createGroup(groupName.trim(), "family");
    setGroupName("");
    load();
  }

  async function openGroup(id: string) {
    setSelected(await groupsApi.getGroup(id));
  }

  async function handleInvite() {
    if (!selected || !inviteEmail.trim()) return;
    await groupsApi.inviteMember(selected.group.id, inviteEmail.trim());
    setInviteEmail("");
    openGroup(selected.group.id);
  }

  return (
    <div className="max-w-2xl space-y-6">
      <h1 className="text-xl font-semibold text-neutral-900 dark:text-neutral-50">Familia</h1>

      <div className="flex gap-2">
        <input
          value={groupName}
          onChange={(e) => setGroupName(e.target.value)}
          placeholder="Nombre del grupo familiar"
          className="flex-1 rounded-md border border-neutral-300 dark:border-neutral-700 bg-white dark:bg-neutral-900 text-neutral-900 dark:text-neutral-100 px-3 py-2"
        />
        <button onClick={handleCreate} className="bg-neutral-900 text-white dark:bg-white dark:text-neutral-900 px-4 py-2 rounded-md text-sm font-medium">
          Crear
        </button>
      </div>

      <ul className="space-y-2">
        {memberships.map((m) => (
          <li key={m.group.id} className="bg-white dark:bg-neutral-900 border border-neutral-200 dark:border-neutral-800 rounded-lg p-4">
            <button onClick={() => openGroup(m.group.id)} className="font-medium text-neutral-800 dark:text-neutral-200">
              {m.group.name}
            </button>
            <span className="ml-2 text-xs text-neutral-400 dark:text-neutral-500">({m.role})</span>
          </li>
        ))}
        {memberships.length === 0 && <p className="text-neutral-500 dark:text-neutral-400 text-sm">Todavía no formás parte de un grupo familiar.</p>}
      </ul>

      {selected && (
        <div className="bg-white dark:bg-neutral-900 border border-neutral-200 dark:border-neutral-800 rounded-lg p-4">
          <h2 className="font-medium text-neutral-800 dark:text-neutral-200 mb-2">{selected.group.name}</h2>
          <ul className="text-sm text-neutral-600 dark:text-neutral-400 mb-4">
            {selected.members.map((m) => (
              <li key={m.user_id}>
                {m.user_id} — {m.role}
              </li>
            ))}
          </ul>
          <div className="flex gap-2">
            <input
              value={inviteEmail}
              onChange={(e) => setInviteEmail(e.target.value)}
              placeholder="Email a invitar"
              className="flex-1 rounded-md border border-neutral-300 dark:border-neutral-700 bg-white dark:bg-neutral-900 text-neutral-900 dark:text-neutral-100 px-3 py-2"
            />
            <button onClick={handleInvite} className="bg-neutral-100 dark:bg-neutral-800 text-neutral-700 dark:text-neutral-300 px-4 py-2 rounded-md text-sm font-medium">
              Invitar
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
