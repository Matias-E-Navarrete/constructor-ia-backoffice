import { useState, type FormEvent } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import ThemeToggle from "../components/ThemeToggle";

export default function Login() {
  const { login } = useAuth();
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
      await login(email, password);
      navigate("/app/habits");
    } catch {
      setError("Email o contraseña incorrectos");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-white dark:bg-black relative">
      <div className="absolute top-4 right-4 w-40">
        <ThemeToggle />
      </div>
      <form onSubmit={handleSubmit} className="w-full max-w-sm bg-white dark:bg-neutral-900 p-8 rounded-xl shadow-sm border border-neutral-200 dark:border-neutral-800">
        <h1 className="text-2xl font-semibold mb-6 text-neutral-900 dark:text-neutral-50">Iniciar sesión</h1>
        {error && <p className="mb-4 text-sm text-red-600 dark:text-red-400">{error}</p>}
        <label htmlFor="email" className="block text-sm font-medium text-neutral-700 dark:text-neutral-300 mb-1">Email</label>
        <input
          id="email"
          type="email"
          required
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="w-full mb-4 rounded-md border border-neutral-300 dark:border-neutral-700 bg-white dark:bg-neutral-900 text-neutral-900 dark:text-neutral-100 px-3 py-2"
        />
        <label htmlFor="password" className="block text-sm font-medium text-neutral-700 dark:text-neutral-300 mb-1">Contraseña</label>
        <input
          id="password"
          type="password"
          required
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          className="w-full mb-6 rounded-md border border-neutral-300 dark:border-neutral-700 bg-white dark:bg-neutral-900 text-neutral-900 dark:text-neutral-100 px-3 py-2"
        />
        <button
          type="submit"
          disabled={submitting}
          className="w-full bg-neutral-900 text-white dark:bg-white dark:text-neutral-900 rounded-md py-2 font-medium disabled:opacity-50"
        >
          {submitting ? "Ingresando..." : "Ingresar"}
        </button>
        <p className="mt-4 text-sm text-neutral-600 dark:text-neutral-400 text-center">
          ¿No tenés cuenta? <Link to="/register" className="text-neutral-900 dark:text-neutral-100 font-medium">Registrate</Link>
        </p>
      </form>
    </div>
  );
}
