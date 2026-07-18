import { useEffect, useRef, useState } from "react";
import * as notesApi from "../../api/notes";
import { UpgradeRequiredError } from "../../api/client";
import UpgradeWall from "../../components/UpgradeWall";
import { useTheme } from "../../context/ThemeContext";

type PositionedNode = notesApi.GraphNode & { x: number; y: number; vx: number; vy: number };

const WIDTH = 640;
const HEIGHT = 420;

// Notes are the hub of the vault, so they get the one neon accent; every
// other entity type stays a plain neutral dot to keep the graph mono-first.
const kindColorLight: Record<string, string> = {
  note: "#a21caf",
  habit: "#a3a3a3",
  workout: "#a3a3a3",
  finance: "#a3a3a3",
};

const kindColorDark: Record<string, string> = {
  note: "#e879f9",
  habit: "#737373",
  workout: "#737373",
  finance: "#737373",
};

// A small hand-rolled force simulation (repulsion between all nodes, spring
// along edges, mild pull to center) so the graph settles without pulling in
// an external physics library for a handful of nodes.
function useForceLayout(graph: notesApi.Graph | null) {
  const [nodes, setNodes] = useState<PositionedNode[]>([]);

  useEffect(() => {
    if (!graph) return;
    const g = graph;
    let positioned: PositionedNode[] = g.nodes.map((n, i) => ({
      ...n,
      x: WIDTH / 2 + Math.cos(i) * 100,
      y: HEIGHT / 2 + Math.sin(i) * 100,
      vx: 0,
      vy: 0,
    }));

    let frame = 0;
    let raf: number;
    function tick() {
      const byId = new Map(positioned.map((n) => [n.id, n]));

      for (const a of positioned) {
        let fx = 0;
        let fy = 0;
        for (const b of positioned) {
          if (a === b) continue;
          const dx = a.x - b.x;
          const dy = a.y - b.y;
          const distSq = Math.max(dx * dx + dy * dy, 1);
          const force = 1200 / distSq;
          fx += (dx / Math.sqrt(distSq)) * force;
          fy += (dy / Math.sqrt(distSq)) * force;
        }
        fx += (WIDTH / 2 - a.x) * 0.01;
        fy += (HEIGHT / 2 - a.y) * 0.01;
        a.vx = (a.vx + fx) * 0.6;
        a.vy = (a.vy + fy) * 0.6;
      }

      for (const e of g.edges) {
        const from = byId.get(e.from);
        const to = byId.get(e.to);
        if (!from || !to) continue;
        const dx = to.x - from.x;
        const dy = to.y - from.y;
        const dist = Math.sqrt(dx * dx + dy * dy) || 1;
        const stretch = (dist - 120) * 0.02;
        const ux = dx / dist;
        const uy = dy / dist;
        from.vx += ux * stretch;
        from.vy += uy * stretch;
        to.vx -= ux * stretch;
        to.vy -= uy * stretch;
      }

      positioned = positioned.map((n) => ({
        ...n,
        x: Math.min(WIDTH - 20, Math.max(20, n.x + n.vx)),
        y: Math.min(HEIGHT - 20, Math.max(20, n.y + n.vy)),
      }));

      frame++;
      setNodes(positioned);
      if (frame < 120) raf = requestAnimationFrame(tick);
    }
    raf = requestAnimationFrame(tick);
    return () => cancelAnimationFrame(raf);
  }, [graph]);

  return nodes;
}

export default function NotesGraph() {
  const [graph, setGraph] = useState<notesApi.Graph | null>(null);
  const [locked, setLocked] = useState(false);
  const loaded = useRef(false);
  const nodes = useForceLayout(graph);
  const { theme } = useTheme();
  const kindColor = theme === "dark" ? kindColorDark : kindColorLight;
  const edgeColor = theme === "dark" ? "#404040" : "#cbd5e1";
  const labelColor = theme === "dark" ? "#e2e8f0" : "#334155";

  useEffect(() => {
    if (loaded.current) return;
    loaded.current = true;
    notesApi
      .getGraph()
      .then(setGraph)
      .catch((err) => {
        if (err instanceof UpgradeRequiredError) setLocked(true);
      });
  }, []);

  if (locked) return <UpgradeWall feature="Vista de grafo de notas" />;
  if (!graph) return <p className="text-neutral-500 dark:text-neutral-400">Cargando grafo...</p>;

  const byId = new Map(nodes.map((n) => [n.id, n]));

  return (
    <div>
      <h1 className="text-xl font-semibold text-neutral-900 dark:text-neutral-50 mb-4">Grafo de notas</h1>
      <svg width={WIDTH} height={HEIGHT} className="bg-white dark:bg-neutral-900 border border-neutral-200 dark:border-neutral-800 rounded-lg">
        {graph.edges.map((e, i) => {
          const from = byId.get(e.from);
          const to = byId.get(e.to);
          if (!from || !to) return null;
          return <line key={i} x1={from.x} y1={from.y} x2={to.x} y2={to.y} stroke={edgeColor} strokeWidth={1} />;
        })}
        {nodes.map((n) => (
          <g key={n.id}>
            <circle cx={n.x} cy={n.y} r={8} fill={kindColor[n.kind] ?? "#737373"} />
            <text x={n.x + 12} y={n.y + 4} fontSize={11} fill={labelColor}>
              {n.title}
            </text>
          </g>
        ))}
      </svg>
    </div>
  );
}
