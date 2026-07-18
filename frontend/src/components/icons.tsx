// A small hand-drawn set of line icons (24x24, stroke=currentColor) used only
// for nav/section identity — no icon library dependency, matching the app's
// "simple, self-contained" convention.
import type { SVGProps } from "react";

function base(props: SVGProps<SVGSVGElement>) {
  return {
    width: 18,
    height: 18,
    viewBox: "0 0 24 24",
    fill: "none",
    stroke: "currentColor",
    strokeWidth: 1.75,
    strokeLinecap: "round" as const,
    strokeLinejoin: "round" as const,
    ...props,
  };
}

export function IconTasks(props: SVGProps<SVGSVGElement>) {
  return (
    <svg {...base(props)}>
      <path d="M9 11l2.5 2.5L16 9" />
      <rect x="3.5" y="3.5" width="17" height="17" rx="4" />
    </svg>
  );
}

export function IconPlanner(props: SVGProps<SVGSVGElement>) {
  return (
    <svg {...base(props)}>
      <rect x="3.5" y="4.5" width="17" height="16" rx="3" />
      <path d="M3.5 9.5h17M8 3v3M16 3v3M8 13.5h3M8 17h6" />
    </svg>
  );
}

export function IconInbox(props: SVGProps<SVGSVGElement>) {
  return (
    <svg {...base(props)}>
      <path d="M3.5 12h5l1.7 3h3.6l1.7-3h5" />
      <path d="M6.2 5.5h11.6l2.2 6.5v7a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2v-7z" />
    </svg>
  );
}

export function IconHabits(props: SVGProps<SVGSVGElement>) {
  return (
    <svg {...base(props)}>
      <path d="M12 21c-3.5-2.4-7-5.7-7-9.8C5 7.6 7.4 5 10.3 5c.9 0 1.9.4 2.7 1.5C13.8 5.4 14.8 5 15.7 5 18.6 5 21 7.6 21 11.2c0 4.1-3.5 7.4-7 9.8-.7.4-1.3.4-2 0z" />
      <path d="M9.5 12l1.7 1.7L14.5 10" />
    </svg>
  );
}

export function IconPanorama(props: SVGProps<SVGSVGElement>) {
  return (
    <svg {...base(props)}>
      <path d="M4 20V10M10 20V4M16 20v-7M21 20H3" />
    </svg>
  );
}

export function IconWorkouts(props: SVGProps<SVGSVGElement>) {
  return (
    <svg {...base(props)}>
      <path d="M6.5 7v10M17.5 7v10M2.5 10v4M21.5 10v4M6.5 12h11" />
    </svg>
  );
}

export function IconStudies(props: SVGProps<SVGSVGElement>) {
  return (
    <svg {...base(props)}>
      <path d="M4 5.5A2.5 2.5 0 0 1 6.5 3H12v18H6.5A2.5 2.5 0 0 1 4 18.5z" />
      <path d="M20 5.5A2.5 2.5 0 0 0 17.5 3H12v18h5.5a2.5 2.5 0 0 0 2.5-2.5z" />
    </svg>
  );
}

export function IconFinance(props: SVGProps<SVGSVGElement>) {
  return (
    <svg {...base(props)}>
      <path d="M3.5 7.5A2.5 2.5 0 0 1 6 5h11a2.5 2.5 0 0 1 2.5 2.5v9A2.5 2.5 0 0 1 17 19H6a2.5 2.5 0 0 1-2.5-2.5z" />
      <path d="M15.5 12.5h2.5v3h-2.5a1.5 1.5 0 0 1 0-3z" />
      <path d="M3.5 9h16" />
    </svg>
  );
}

export function IconNotes(props: SVGProps<SVGSVGElement>) {
  return (
    <svg {...base(props)}>
      <path d="M7 3.5h7l4 4V19a1.5 1.5 0 0 1-1.5 1.5H7A1.5 1.5 0 0 1 5.5 19V5A1.5 1.5 0 0 1 7 3.5z" />
      <path d="M14 3.5V8h4M8.5 12h7M8.5 15.5h5" />
    </svg>
  );
}

export function IconFamily(props: SVGProps<SVGSVGElement>) {
  return (
    <svg {...base(props)}>
      <circle cx="9" cy="8" r="3" />
      <path d="M3.5 20c0-3 2.5-5 5.5-5s5.5 2 5.5 5" />
      <circle cx="17" cy="8.5" r="2.2" />
      <path d="M15.8 12.2c2 .3 3.7 1.9 3.7 4.3" />
    </svg>
  );
}

export function IconCoaching(props: SVGProps<SVGSVGElement>) {
  return (
    <svg {...base(props)}>
      <circle cx="12" cy="8" r="3.2" />
      <path d="M5.5 20c0-3.6 2.9-6 6.5-6s6.5 2.4 6.5 6" />
      <path d="M9 8l2 2 3-3.5" />
    </svg>
  );
}

export function IconAdmin(props: SVGProps<SVGSVGElement>) {
  return (
    <svg {...base(props)}>
      <path d="M12 3l7 3v5.5c0 4.6-3 8.4-7 9.5-4-1.1-7-4.9-7-9.5V6z" />
      <path d="M9 12l2 2 4-4" />
    </svg>
  );
}

export function IconUpgrade(props: SVGProps<SVGSVGElement>) {
  return (
    <svg {...base(props)}>
      <path d="M12 3l1.8 4.9L19 9l-4.9 1.8L12 16l-1.8-5.2L5 9l5.2-1.1z" />
      <path d="M19 15l.8 2.2L22 18l-2.2.8L19 21l-.8-2.2L16 18l2.2-.8z" />
    </svg>
  );
}

export function IconSun(props: SVGProps<SVGSVGElement>) {
  return (
    <svg {...base(props)}>
      <circle cx="12" cy="12" r="4" />
      <path d="M12 2.5v2M12 19.5v2M4.2 4.2l1.4 1.4M18.4 18.4l1.4 1.4M2.5 12h2M19.5 12h2M4.2 19.8l1.4-1.4M18.4 5.6l1.4-1.4" />
    </svg>
  );
}

export function IconMoon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg {...base(props)}>
      <path d="M20 14.5A8.5 8.5 0 1 1 9.5 4a7 7 0 0 0 10.5 10.5z" />
    </svg>
  );
}

export function IconRoadmap(props: SVGProps<SVGSVGElement>) {
  return (
    <svg {...base(props)}>
      <path d="M4 5h11l3 3-3 3H4z" />
      <path d="M4 5v15" />
    </svg>
  );
}

export function IconFlag(props: SVGProps<SVGSVGElement>) {
  return (
    <svg {...base(props)}>
      <rect x="4" y="3.5" width="6" height="6" rx="1.2" />
      <rect x="14" y="3.5" width="6" height="6" rx="1.2" />
      <rect x="4" y="14.5" width="6" height="6" rx="1.2" />
      <rect x="14" y="14.5" width="6" height="6" rx="1.2" />
    </svg>
  );
}

export function IconUsersGear(props: SVGProps<SVGSVGElement>) {
  return (
    <svg {...base(props)}>
      <circle cx="9" cy="8" r="3.2" />
      <path d="M3.5 20c0-3.6 2.5-6 5.5-6s5.5 2.4 5.5 6" />
      <circle cx="18" cy="17" r="1.6" />
      <path d="M18 14v-.8M18 20.8V20M20.6 15.4l-.7.4M15.4 18.6l-.7.4M20.6 18.6l-.7-.4M15.4 15.4l-.7.4" />
    </svg>
  );
}

export function IconTests(props: SVGProps<SVGSVGElement>) {
  return (
    <svg {...base(props)}>
      <path d="M9 3.5h6M10 3.5v5.3L5.5 17a2 2 0 0 0 1.8 2.9h9.4A2 2 0 0 0 18.5 17L14 8.8V3.5" />
      <path d="M7.5 15h9" />
    </svg>
  );
}

export function IconArrowLeft(props: SVGProps<SVGSVGElement>) {
  return (
    <svg {...base(props)}>
      <path d="M19 12H5M11 6l-6 6 6 6" />
    </svg>
  );
}

export function IconLogout(props: SVGProps<SVGSVGElement>) {
  return (
    <svg {...base(props)}>
      <path d="M9 4H6a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h3M15 16l4-4-4-4M19 12H9" />
    </svg>
  );
}
