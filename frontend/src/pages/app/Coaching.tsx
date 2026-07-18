import { useEffect, useState } from "react";
import * as groupsApi from "../../api/groups";
import * as habitsApi from "../../api/habits";
import * as workoutsApi from "../../api/workouts";
import { UpgradeRequiredError } from "../../api/client";
import UpgradeWall from "../../components/UpgradeWall";

export default function Coaching() {
  const [memberships, setMemberships] = useState<groupsApi.Membership[]>([]);
  const [groupName, setGroupName] = useState("");
  const [locked, setLocked] = useState(false);
  const [selectedGroup, setSelectedGroup] = useState<{ group: groupsApi.Group; members: groupsApi.Member[] } | null>(null);
  const [inviteEmail, setInviteEmail] = useState("");
  const [clientHabits, setClientHabits] = useState<habitsApi.Habit[] | null>(null);
  const [clientSessions, setClientSessions] = useState<workoutsApi.Session[] | null>(null);

  async function load() {
    const all = await groupsApi.listMyGroups();
    setMemberships(all.filter((m) => m.group.kind === "coaching"));
  }

  useEffect(() => {
    load();
  }, []);

  async function handleCreate() {
    if (!groupName.trim()) return;
    setLocked(false);
    try {
      await groupsApi.createGroup(groupName.trim(), "coaching");
      setGroupName("");
      load();
    } catch (err) {
      if (err instanceof UpgradeRequiredError) setLocked(true);
    }
  }

  async function openGroup(id: string) {
    setSelectedGroup(await groupsApi.getGroup(id));
    setClientHabits(null);
    setClientSessions(null);
  }

  async function handleInvite() {
    if (!selectedGroup || !inviteEmail.trim()) return;
    await groupsApi.inviteMember(selectedGroup.group.id, inviteEmail.trim());
    setInviteEmail("");
    openGroup(selectedGroup.group.id);
  }

  async function viewClient(userId: string) {
    setClientHabits(await habitsApi.listHabits(userId));
    setClientSessions(await workoutsApi.listSessions(userId));
  }

  return (
    <div className="max-w-2xl space-y-6">
      <h1 className="text-xl font-semibold text-neutral-900 dark:text-neutral-50">Coaching</h1>

      <div className="flex gap-2">
        <input
          value={groupName}
          onChange={(e) => setGroupName(e.target.value)}
          placeholder="Nombre del grupo de coaching"
          className="input-field flex-1"
        />
        <button onClick={handleCreate} className="btn-primary">
          Crear
        </button>
      </div>
      {locked && <UpgradeWall feature="Crear un grupo de coaching" />}

      <ul className="space-y-2">
        {memberships.map((m) => (
          <li key={m.group.id} className="surface-card-hover p-4">
            <button onClick={() => openGroup(m.group.id)} className="font-medium text-neutral-800 dark:text-neutral-200 hover:text-accent-strong dark:hover:text-accent-soft transition-colors duration-150">
              {m.group.name}
            </button>
            <span className="ml-2 text-xs text-neutral-400 dark:text-neutral-500">({m.role})</span>
          </li>
        ))}
        {memberships.length === 0 && <p className="text-neutral-500 dark:text-neutral-400 text-sm">Todavía no tenés grupos de coaching.</p>}
      </ul>

      {selectedGroup && (
        <div className="surface-card p-4 space-y-4">
          <h2 className="font-medium text-neutral-800 dark:text-neutral-200">{selectedGroup.group.name}</h2>

          <div className="flex gap-2">
            <input
              value={inviteEmail}
              onChange={(e) => setInviteEmail(e.target.value)}
              placeholder="Email del cliente"
              className="input-field flex-1"
            />
            <button onClick={handleInvite} className="btn-secondary">
              Invitar
            </button>
          </div>

          <ul className="text-sm text-neutral-600 dark:text-neutral-400">
            {selectedGroup.members
              .filter((m) => m.role === "client")
              .map((m) => (
                <li key={m.user_id} className="flex items-center justify-between">
                  <span>{m.user_id}</span>
                  <button onClick={() => viewClient(m.user_id)} className="btn-secondary text-xs px-2 py-1">
                    Ver progreso
                  </button>
                </li>
              ))}
          </ul>

          {clientHabits && (
            <div>
              <p className="text-sm font-medium text-neutral-700 dark:text-neutral-300">Hábitos del cliente</p>
              <ul className="text-sm text-neutral-600 dark:text-neutral-400">
                {clientHabits.map((h) => (
                  <li key={h.id}>{h.name}</li>
                ))}
              </ul>
            </div>
          )}
          {clientSessions && (
            <div>
              <p className="text-sm font-medium text-neutral-700 dark:text-neutral-300">Entrenamientos del cliente</p>
              <ul className="text-sm text-neutral-600 dark:text-neutral-400">
                {clientSessions.map((s) => (
                  <li key={s.id}>{s.session_date}</li>
                ))}
              </ul>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
