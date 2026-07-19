import { apiFetch, ApiError } from "./client";

export type AssistantMessage = { id: string; role: "user" | "assistant"; content: string; created_at: string };

export class AssistantNotConfiguredError extends ApiError {}

export function listMessages() {
  return apiFetch<AssistantMessage[]>("/api/assistant/messages");
}

export async function sendMessage(content: string) {
  try {
    return await apiFetch<AssistantMessage>("/api/assistant/messages", { method: "POST", body: { content } });
  } catch (err) {
    if (err instanceof ApiError && err.code === "assistant_not_configured") {
      throw new AssistantNotConfiguredError(err.status, err.code, err.message);
    }
    throw err;
  }
}

export function clearConversation() {
  return apiFetch<void>("/api/assistant/messages", { method: "DELETE" });
}
