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
      <h1 className="text-xl font-semibold text-slate-900">Coaching</h1>

      <div className="flex gap-2">
        <input
          value={groupName}
          onChange={(e) => setGroupName(e.target.value)}
          placeholder="Nombre del grupo de coaching"
          className="flex-1 rounded-md border border-slate-300 px-3 py-2"
        />
        <button onClick={handleCreate} className="bg-slate-900 text-white px-4 py-2 rounded-md text-sm font-medium">
          Crear
        </button>
      </div>
      {locked && <UpgradeWall feature="Crear un grupo de coaching" />}

      <ul className="space-y-2">
        {memberships.map((m) => (
          <li key={m.group.id} className="bg-white border border-slate-200 rounded-lg p-4">
            <button onClick={() => openGroup(m.group.id)} className="font-medium text-slate-800">
              {m.group.name}
            </button>
            <span className="ml-2 text-xs text-slate-400">({m.role})</span>
          </li>
        ))}
        {memberships.length === 0 && <p className="text-slate-500 text-sm">Todavía no tenés grupos de coaching.</p>}
      </ul>

      {selectedGroup && (
        <div className="bg-white border border-slate-200 rounded-lg p-4 space-y-4">
          <h2 className="font-medium text-slate-800">{selectedGroup.group.name}</h2>

          <div className="flex gap-2">
            <input
              value={inviteEmail}
              onChange={(e) => setInviteEmail(e.target.value)}
              placeholder="Email del cliente"
              className="flex-1 rounded-md border border-slate-300 px-3 py-2"
            />
            <button onClick={handleInvite} className="bg-slate-100 text-slate-700 px-4 py-2 rounded-md text-sm font-medium">
              Invitar
            </button>
          </div>

          <ul className="text-sm text-slate-600">
            {selectedGroup.members
              .filter((m) => m.role === "client")
              .map((m) => (
                <li key={m.user_id} className="flex items-center justify-between">
                  <span>{m.user_id}</span>
                  <button onClick={() => viewClient(m.user_id)} className="text-xs bg-slate-100 px-2 py-1 rounded-md">
                    Ver progreso
                  </button>
                </li>
              ))}
          </ul>

          {clientHabits && (
            <div>
              <p className="text-sm font-medium text-slate-700">Hábitos del cliente</p>
              <ul className="text-sm text-slate-600">
                {clientHabits.map((h) => (
                  <li key={h.id}>{h.name}</li>
                ))}
              </ul>
            </div>
          )}
          {clientSessions && (
            <div>
              <p className="text-sm font-medium text-slate-700">Entrenamientos del cliente</p>
              <ul className="text-sm text-slate-600">
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
