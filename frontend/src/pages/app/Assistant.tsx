import { useEffect, useRef, useState } from "react";
import * as assistantApi from "../../api/assistant";
import { UpgradeRequiredError, FeatureDisabledError } from "../../api/client";
import UpgradeWall from "../../components/UpgradeWall";
import FeatureDisabledBanner from "../../components/FeatureDisabledBanner";
import { IconAssistant } from "../../components/icons";

export default function Assistant() {
  const [messages, setMessages] = useState<assistantApi.AssistantMessage[]>([]);
  const [disabled, setDisabled] = useState(false);
  const [locked, setLocked] = useState(false);
  const [notConfigured, setNotConfigured] = useState(false);
  const [sendError, setSendError] = useState(false);
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
    const pendingId = `pending-${Date.now()}`;
    setDraft("");
    setNotConfigured(false);
    setSendError(false);
    setMessages((m) => [...m, { id: pendingId, role: "user", content, created_at: new Date().toISOString() }]);
    setSending(true);
    try {
      const reply = await assistantApi.sendMessage(content);
      setMessages((m) => [...m, reply]);
    } catch (err) {
      if (err instanceof assistantApi.AssistantNotConfiguredError) {
        setNotConfigured(true);
      } else if (err instanceof UpgradeRequiredError) {
        setLocked(true);
      } else {
        setMessages((m) => m.filter((msg) => msg.id !== pendingId));
        setDraft(content);
        setSendError(true);
      }
    } finally {
      setSending(false);
    }
  }

  async function handleClear() {
    await assistantApi.clearConversation();
    setMessages([]);
  }

  if (disabled) return <FeatureDisabledBanner />;

  const inputBar = (
    <div className="flex gap-2">
      <input
        value={draft}
        onChange={(e) => setDraft(e.target.value)}
        onKeyDown={(e) => e.key === "Enter" && handleSend()}
        placeholder="Escribí tu mensaje..."
        className="input-field flex-1 py-3"
        disabled={sending}
      />
      <button onClick={handleSend} disabled={sending} className="btn-accent px-5">
        Enviar
      </button>
    </div>
  );

  return (
    <div className="max-w-3xl mx-auto w-full h-[calc(100vh-6rem)] flex flex-col">
      {locked && <UpgradeWall feature="El asistente de IA" />}

      {!locked && messages.length === 0 && (
        <div className="flex-1 flex flex-col items-center justify-center gap-6 px-4">
          <div className="w-12 h-12 rounded-xl bg-accent flex items-center justify-center shadow-neon">
            <IconAssistant width={22} height={22} className="text-white" />
          </div>
          <div className="text-center">
            <h1 className="text-2xl font-semibold text-neutral-900 dark:text-neutral-50 mb-1.5">¿En qué te ayudo hoy?</h1>
            <p className="text-sm text-neutral-500 dark:text-neutral-400">Preguntame sobre tus tareas, hábitos, entrenamientos, estudio o finanzas.</p>
          </div>
          <div className="w-full max-w-xl">{inputBar}</div>
          {notConfigured && (
            <p className="text-xs text-neutral-500 dark:text-neutral-400 text-center max-w-xl">
              El asistente todavía no está configurado en este entorno (falta la API key). Tu mensaje quedó guardado.
            </p>
          )}
          {sendError && (
            <p className="text-xs text-red-500 dark:text-red-400 text-center max-w-xl">
              No se pudo enviar el mensaje. Intentá de nuevo.
            </p>
          )}
        </div>
      )}

      {!locked && messages.length > 0 && (
        <>
          <div className="flex items-center justify-between mb-3 px-1">
            <span className="text-sm font-medium text-neutral-500 dark:text-neutral-400">Asistente</span>
            <button onClick={handleClear} className="text-xs text-neutral-400 dark:text-neutral-500 hover:text-red-500 dark:hover:text-red-400 transition-colors duration-150">
              Borrar conversación
            </button>
          </div>

          <div className="flex-1 overflow-y-auto px-1 space-y-6">
            {messages.map((m) =>
              m.role === "user" ? (
                <div key={m.id} className="flex justify-end">
                  <div className="max-w-[75%] rounded-2xl bg-neutral-100 dark:bg-neutral-800 px-4 py-2.5 text-sm text-neutral-900 dark:text-neutral-100 whitespace-pre-wrap">
                    {m.content}
                  </div>
                </div>
              ) : (
                <div key={m.id} className="flex gap-3">
                  <div className="w-6 h-6 rounded-md bg-accent flex items-center justify-center flex-shrink-0 mt-0.5 shadow-neon-sm">
                    <IconAssistant width={13} height={13} className="text-white" />
                  </div>
                  <p className="flex-1 pt-0.5 text-sm text-neutral-800 dark:text-neutral-200 leading-relaxed whitespace-pre-wrap">{m.content}</p>
                </div>
              )
            )}
            {sending && (
              <div className="flex gap-3">
                <div className="w-6 h-6 rounded-md bg-accent flex items-center justify-center flex-shrink-0 mt-0.5">
                  <IconAssistant width={13} height={13} className="text-white" />
                </div>
                <p className="pt-0.5 text-sm text-neutral-400 dark:text-neutral-500">Pensando...</p>
              </div>
            )}
            <div ref={bottomRef} />
          </div>

          {notConfigured && (
            <p className="text-xs text-neutral-500 dark:text-neutral-400 mb-2 px-1">
              El asistente todavía no está configurado en este entorno (falta la API key). Tu mensaje quedó guardado.
            </p>
          )}
          {sendError && (
            <p className="text-xs text-red-500 dark:text-red-400 mb-2 px-1">No se pudo enviar el mensaje. Intentá de nuevo.</p>
          )}

          <div className="pt-3 border-t border-neutral-200 dark:border-neutral-800">{inputBar}</div>
        </>
      )}
    </div>
  );
}
