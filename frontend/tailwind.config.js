/** @type {import('tailwindcss').Config} */
export default {
  darkMode: "class",
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        // The single neon accent, reserved for a handful of specific
        // highlights (PRO badge, upgrade CTA, focus rings, the notes hub).
        // Everything else stays strictly neutral grayscale.
        accent: {
          DEFAULT: "#22c55e", // solid fills (buttons, badges) — same in both themes
          strong: "#15803d", // text/icon/border accent on light (white) backgrounds
          soft: "#4ade80", // text/icon/border accent on dark (black) backgrounds — glows
        },
      },
      boxShadow: {
        // A soft neon glow, used sparingly on the accent CTA/PRO badge — the
        // one "specific detail" allowed to feel like light instead of ink.
        glow: "0 0 0 1px rgba(34,197,94,0.25), 0 0 20px -2px rgba(34,197,94,0.45)",
        "glow-sm": "0 0 0 1px rgba(34,197,94,0.2), 0 0 10px -3px rgba(34,197,94,0.4)",
        // A crisper "lit sign" version: a sharp colored border ring plus a
        // wider halo, for the components that should read as electric/neon
        // (accent buttons, badges, active nav item, hovered cards).
        neon: "0 0 0 1.5px rgba(34,197,94,0.9), 0 0 1px 1px rgba(34,197,94,0.75), 0 0 16px 2px rgba(34,197,94,0.55), 0 0 34px -4px rgba(34,197,94,0.65)",
        "neon-sm": "0 0 0 1.5px rgba(34,197,94,0.85), 0 0 10px 1px rgba(34,197,94,0.5), 0 0 20px -6px rgba(34,197,94,0.55)",
      },
    },
  },
  plugins: [],
}
