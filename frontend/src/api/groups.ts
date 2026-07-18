import { apiFetch } from "./client";

export type Group = { id: string; name: string; kind: "family" | "coaching"; owner_user_id: string };
export type Membership = { group: Group; role: "owner" | "member" | "coach" | "client" };
export type Member = { user_id: string; role: string };

export function listMyGroups() {
  return apiFetch<Membership[]>("/api/groups");
}

export function createGroup(name: string, kind: "family" | "coaching") {
  return apiFetch<Group>("/api/groups", { method: "POST", body: { name, kind } });
}

export function getGroup(id: string) {
  return apiFetch<{ group: Group; members: Member[] }>(`/api/groups/${id}`);
}

export function inviteMember(groupId: string, email: string) {
  return apiFetch<Member>(`/api/groups/${groupId}/members`, { method: "POST", body: { email } });
}

export function removeMember(groupId: string, userId: string) {
  return apiFetch<void>(`/api/groups/${groupId}/members/${userId}`, { method: "DELETE" });
}
