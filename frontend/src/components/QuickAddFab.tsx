import { useState } from "react";
import { useNavigate } from "react-router-dom";

const options = [
  { label: "Tarea", to: "/app/tasks" },
  { label: "Hábito", to: "/app/habits" },
  { label: "Entrenamiento", to: "/app/workouts" },
  { label: "Transacción", to: "/app/finance" },
  { label: "Nota", to: "/app/notes" },
  { label: "Captura rápida (Inbox)", to: "/app/inbox" },
];

export default function QuickAddFab() {
  const [open, setOpen] = useState(false);
  const navigate = useNavigate();

  return (
    <div className="fixed bottom-6 right-6 z-30">
      {open && (
        <div className="mb-2 surface-card shadow-lg overflow-hidden w-52">
          {options.map((o) => (
            <button
              key={o.to}
              onClick={() => {
                setOpen(false);
                navigate(o.to);
              }}
              className="w-full text-left px-4 py-2.5 text-sm text-neutral-700 dark:text-neutral-300 hover:bg-neutral-100 dark:hover:bg-neutral-800 transition-colors duration-150"
            >
              {o.label}
            </button>
          ))}
        </div>
      )}
      <button
        onClick={() => setOpen((o) => !o)}
        aria-label="Crear nuevo"
        className="w-14 h-14 rounded-full bg-accent text-white text-2xl leading-none flex items-center justify-center shadow-neon transition-all duration-150 hover:scale-105 active:scale-95"
      >
        {open ? "×" : "+"}
      </button>
    </div>
  );
}
