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
          DEFAULT: "#d946ef", // solid fills (buttons, badges) — same in both themes
          strong: "#a21caf", // text/icon/border accent on light (white) backgrounds
          soft: "#e879f9", // text/icon/border accent on dark (black) backgrounds — glows
        },
      },
      boxShadow: {
        // A soft neon glow, used sparingly on the accent CTA/PRO badge — the
        // one "specific detail" allowed to feel like light instead of ink.
        glow: "0 0 0 1px rgba(217,70,239,0.25), 0 0 20px -2px rgba(217,70,239,0.45)",
        "glow-sm": "0 0 0 1px rgba(217,70,239,0.2), 0 0 10px -3px rgba(217,70,239,0.4)",
      },
    },
  },
  plugins: [],
}
