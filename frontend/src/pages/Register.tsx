import { useState, type FormEvent } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

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
      navigate("/app/habits");
    } catch {
      setError("No pudimos crear tu cuenta. Probá con otro email o una contraseña más larga.");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-50">
      <form onSubmit={handleSubmit} className="w-full max-w-sm bg-white p-8 rounded-xl shadow-sm border border-slate-200">
        <h1 className="text-2xl font-semibold mb-6 text-slate-900">Crear cuenta</h1>
        {error && <p className="mb-4 text-sm text-red-600">{error}</p>}
        <label htmlFor="email" className="block text-sm font-medium text-slate-700 mb-1">Email</label>
        <input
          id="email"
          type="email"
          required
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="w-full mb-4 rounded-md border border-slate-300 px-3 py-2"
        />
        <label htmlFor="password" className="block text-sm font-medium text-slate-700 mb-1">Contraseña</label>
        <input
          id="password"
          type="password"
          required
          minLength={8}
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          className="w-full mb-6 rounded-md border border-slate-300 px-3 py-2"
        />
        <button
          type="submit"
          disabled={submitting}
          className="w-full bg-slate-900 text-white rounded-md py-2 font-medium disabled:opacity-50"
        >
          {submitting ? "Creando..." : "Crear cuenta"}
        </button>
        <p className="mt-4 text-sm text-slate-600 text-center">
          ¿Ya tenés cuenta? <Link to="/login" className="text-slate-900 font-medium">Iniciá sesión</Link>
        </p>
      </form>
    </div>
  );
}
