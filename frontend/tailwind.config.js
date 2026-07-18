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
    },
  },
  plugins: [],
}
