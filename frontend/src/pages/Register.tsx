import { useState, type FormEvent } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import ThemeToggle from "../components/ThemeToggle";

export default function Register() {
  const { register } = useAuth();
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);
    setSubmitting(true);
    try {
      await register(email, password);
      navigate("/app/tasks");
    } catch {
      setError("No pudimos crear tu cuenta. Probá con otro email o una contraseña más larga.");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-white dark:bg-black relative">
      <div className="absolute top-4 right-4 w-40">
        <ThemeToggle />
      </div>
      <form onSubmit={handleSubmit} className="w-full max-w-sm surface-card p-8">
        <div className="w-10 h-10 rounded-lg bg-neutral-900 dark:bg-white flex items-center justify-center shadow-glow-sm mb-5">
          <span className="text-white dark:text-neutral-900 font-bold text-base">r</span>
        </div>
        <h1 className="text-2xl font-semibold mb-6 text-neutral-900 dark:text-neutral-50">Crear cuenta</h1>
        {error && <p className="mb-4 text-sm text-red-600 dark:text-red-400">{error}</p>}
        <label htmlFor="email" className="block text-sm font-medium text-neutral-700 dark:text-neutral-300 mb-1">Email</label>
        <input
          id="email"
          type="email"
          required
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="input-field w-full mb-4"
        />
        <label htmlFor="password" className="block text-sm font-medium text-neutral-700 dark:text-neutral-300 mb-1">Contraseña</label>
        <input
          id="password"
          type="password"
          required
          minLength={8}
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          className="input-field w-full mb-6"
        />
        <button
          type="submit"
          disabled={submitting}
          className="btn-primary w-full py-2"
        >
          {submitting ? "Creando..." : "Crear cuenta"}
        </button>
        <p className="mt-4 text-sm text-neutral-600 dark:text-neutral-400 text-center">
          ¿Ya tenés cuenta?{" "}
          <Link to="/login" className="text-neutral-900 dark:text-neutral-100 font-medium hover:opacity-80 transition-opacity duration-150">
            Iniciá sesión
          </Link>
        </p>
      </form>
    </div>
  );
}
