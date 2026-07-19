import { useEffect, useRef, useState } from "react";
import * as assistantApi from "../../api/assistant";
import { UpgradeRequiredError, FeatureDisabledError } from "../../api/client";
import UpgradeWall from "../../components/UpgradeWall";
import FeatureDisabledBanner from "../../components/FeatureDisabledBanner";

export default function Assistant() {
  const [messages, setMessages] = useState<assistantApi.AssistantMessage[]>([]);
  const [disabled, setDisabled] = useState(false);
  const [locked, setLocked] = useState(false);
  const [notConfigured, setNotConfigured] = useState(false);
  const [draft, setDraft] = useState("");
  const [sending, setSending] = useState(false);
  const bottomRef = useRef<HTMLDivElement>(null);

  async function load() {
    try {
      setMessages(await assistantApi.listMessages());
    } catch (err) {
      if (err instanceof UpgradeRequiredError) setLocked(true);
      else if (err instanceof FeatureDisabledError) setDisabled(true);
    }
  }

  useEffect(() => {
    load();
  }, []);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  async function handleSend() {
    if (!draft.trim() || sending) return;
    const content = draft.trim();
    setDraft("");
    setNotConfigured(false);
    setMessages((m) => [...m, { id: `pending-${Date.now()}`, role: "user", content, created_at: new Date().toISOString() }]);
    setSending(true);
    try {
      const reply = await assistantApi.sendMessage(content);
      setMessages((m) => [...m, reply]);
    } catch (err) {
      if (err instanceof assistantApi.AssistantNotConfiguredError) setNotConfigured(true);
      else if (err instanceof UpgradeRequiredError) setLocked(true);
    } finally {
      setSending(false);
    }
  }

  async function handleClear() {
    await assistantApi.clearConversation();
    setMessages([]);
  }

  if (disabled) return <FeatureDisabledBanner />;

  return (
    <div className="max-w-2xl h-full flex flex-col">
      <div className="flex items-center justify-between mb-4">
        <h1 className="text-xl font-semibold text-neutral-900 dark:text-neutral-50">Asistente</h1>
        {messages.length > 0 && (
          <button onClick={handleClear} className="text-xs text-neutral-400 dark:text-neutral-500 hover:text-red-500 dark:hover:text-red-400 transition-colors duration-150">
            Borrar conversación
          </button>
        )}
      </div>

      {locked && <UpgradeWall feature="El asistente de IA" />}

      {!locked && (
        <>
          <div className="surface-card flex-1 p-4 mb-3 min-h-[360px] max-h-[60vh] overflow-y-auto flex flex-col gap-3">
            {messages.length === 0 && (
              <p className="text-neutral-400 dark:text-neutral-500 text-sm m-auto text-center">
                Preguntale lo que quieras sobre tus tareas, hábitos, entrenamientos o finanzas.
              </p>
            )}
            {messages.map((m) => (
              <div key={m.id} className={`flex ${m.role === "user" ? "justify-end" : "justify-start"}`}>
                <div
                  className={`max-w-[80%] rounded-xl px-3 py-2 text-sm whitespace-pre-wrap ${
                    m.role === "user"
                      ? "bg-neutral-900 text-white dark:bg-white dark:text-neutral-900"
                      : "bg-neutral-100 dark:bg-neutral-800 text-neutral-800 dark:text-neutral-200"
                  }`}
                >
                  {m.content}
                </div>
              </div>
            ))}
            {sending && (
              <div className="flex justify-start">
                <div className="bg-neutral-100 dark:bg-neutral-800 text-neutral-400 dark:text-neutral-500 rounded-xl px-3 py-2 text-sm">Pensando...</div>
              </div>
            )}
            <div ref={bottomRef} />
          </div>

          {notConfigured && (
            <p className="text-xs text-neutral-500 dark:text-neutral-400 mb-2">
              El asistente todavía no está configurado en este entorno (falta la API key). Tu mensaje quedó guardado.
            </p>
          )}

          <div className="flex gap-2">
            <input
              value={draft}
              onChange={(e) => setDraft(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && handleSend()}
              placeholder="Escribí tu mensaje..."
              className="input-field flex-1"
              disabled={sending}
            />
            <button onClick={handleSend} disabled={sending} className="btn-accent">
              Enviar
            </button>
          </div>
        </>
      )}
    </div>
  );
}
