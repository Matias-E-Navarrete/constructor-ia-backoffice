import { useEffect, useState } from "react";
import * as habitsApi from "../../api/habits";
import { UpgradeRequiredError, FeatureDisabledError } from "../../api/client";
import UpgradeWall from "../../components/UpgradeWall";
import FeatureDisabledBanner from "../../components/FeatureDisabledBanner";
import { categoryColor } from "../../lib/categoryColor";

const RADAR_SIZE = 240;
const RADAR_CENTER = RADAR_SIZE / 2;
const RADAR_RADIUS = 90;
const RADAR_LEVELS = [0.25, 0.5, 0.75, 1];

function radarPoint(axisIndex: number, axisCount: number, value: number) {
  const angle = -Math.PI / 2 + (axisIndex * 2 * Math.PI) / axisCount;
  const r = value * RADAR_RADIUS;
  return [RADAR_CENTER + r * Math.cos(angle), RADAR_CENTER + r * Math.sin(angle)];
}

function polygonPoints(values: number[]) {
  return values.map((v, i) => radarPoint(i, values.length, v).join(",")).join(" ");
}

function StreakRing({ habit }: { habit: habitsApi.PanoramaHabit }) {
  const radius = 30;
  const circumference = 2 * Math.PI * radius;
  const target = Math.max(habit.best_streak, 7);
  const progress = Math.min(habit.current_streak / target, 1);
  const dot = categoryColor(habit.name);

  return (
    <div className="flex flex-col items-center gap-2" data-testid={`streak-ring-${habit.habit_id}`}>
      <svg width={80} height={80} viewBox="0 0 80 80">
        <circle cx={40} cy={40} r={radius} fill="none" stroke="currentColor" className="text-neutral-200 dark:text-neutral-800" strokeWidth={8} />
        <circle
          cx={40}
          cy={40}
          r={radius}
          fill="none"
          stroke="#d946ef"
          strokeWidth={8}
          strokeLinecap="round"
          strokeDasharray={circumference}
          strokeDashoffset={circumference * (1 - progress)}
          transform="rotate(-90 40 40)"
        />
        <text x={40} y={44} textAnchor="middle" className="fill-neutral-900 dark:fill-neutral-50 text-lg font-semibold">
          {habit.current_streak}
        </text>
      </svg>
      <div className="flex items-center gap-1.5 text-xs text-neutral-600 dark:text-neutral-300">
        <span className="w-2 h-2 rounded-full flex-shrink-0" style={{ background: dot }} />
        {habit.name}
      </div>
      <p className="text-[11px] text-neutral-400 dark:text-neutral-500">Mejor racha: {habit.best_streak}</p>
    </div>
  );
}

function Heatmap({ habit }: { habit: habitsApi.PanoramaHabit }) {
  const [hovered, setHovered] = useState<habitsApi.HeatmapDay | null>(null);
  const weeks: habitsApi.HeatmapDay[][] = [];
  for (let i = 0; i < habit.heatmap.length; i += 7) {
    weeks.push(habit.heatmap.slice(i, i + 7));
  }

  return (
    <div className="mb-4" data-testid={`heatmap-${habit.habit_id}`}>
      <p className="text-xs font-medium text-neutral-600 dark:text-neutral-300 mb-1.5">{habit.name}</p>
      <div className="flex gap-[3px] relative">
        {weeks.map((week, wi) => (
          <div key={wi} className="flex flex-col gap-[3px]">
            {week.map((day) => (
              <div
                key={day.date}
                onMouseEnter={() => setHovered(day)}
                onMouseLeave={() => setHovered(null)}
                title={`${day.date}: ${day.completed ? "Completado" : "Sin registro"}`}
                className={`w-[10px] h-[10px] rounded-sm ${day.completed ? "bg-accent" : "bg-neutral-100 dark:bg-neutral-800"}`}
              />
            ))}
          </div>
        ))}
        {hovered && (
          <div className="absolute -top-6 left-0 text-[11px] bg-neutral-900 dark:bg-white text-white dark:text-neutral-900 px-2 py-0.5 rounded whitespace-nowrap">
            {hovered.date} · {hovered.completed ? "Completado" : "Sin registro"}
          </div>
        )}
      </div>
    </div>
  );
}

export default function Panorama() {
  const [panorama, setPanorama] = useState<habitsApi.PanoramaHabit[] | null>(null);
  const [locked, setLocked] = useState(false);
  const [disabled, setDisabled] = useState(false);

  async function load() {
    setLocked(false);
    try {
      setPanorama(await habitsApi.getPanorama());
    } catch (err) {
      if (err instanceof UpgradeRequiredError) setLocked(true);
      else if (err instanceof FeatureDisabledError) setDisabled(true);
    }
  }

  useEffect(() => {
    load();
  }, []);

  if (disabled) return <FeatureDisabledBanner />;

  return (
    <div className="max-w-3xl space-y-8">
      <h1 className="text-xl font-semibold text-neutral-900 dark:text-neutral-50">Panorama de hábitos</h1>

      {locked && <UpgradeWall feature="Panorama de hábitos" />}

      {panorama && panorama.length === 0 && (
        <p className="text-neutral-500 dark:text-neutral-400 text-sm">Todavía no tenés hábitos activos para mostrar en el panorama.</p>
      )}

      {panorama && panorama.length > 0 && (
        <>
          <div>
            <h2 className="text-sm font-medium text-neutral-500 dark:text-neutral-400 mb-3">Rachas</h2>
            <div className="flex flex-wrap gap-6">
              {panorama.map((h) => (
                <StreakRing key={h.habit_id} habit={h} />
              ))}
            </div>
          </div>

          <div>
            <h2 className="text-sm font-medium text-neutral-500 dark:text-neutral-400 mb-3">Calendario (últimos 90 días)</h2>
            {panorama.map((h) => (
              <Heatmap key={h.habit_id} habit={h} />
            ))}
          </div>

          <div>
            <h2 className="text-sm font-medium text-neutral-500 dark:text-neutral-400 mb-3">Consistencia semanal</h2>
            <div className="flex flex-col md:flex-row items-center gap-6">
              <svg width={RADAR_SIZE} height={RADAR_SIZE} viewBox={`0 0 ${RADAR_SIZE} ${RADAR_SIZE}`}>
                {RADAR_LEVELS.map((level) => (
                  <polygon
                    key={level}
                    points={polygonPoints(new Array(8).fill(level))}
                    fill="none"
                    stroke="currentColor"
                    className="text-neutral-200 dark:text-neutral-800"
                    strokeWidth={1}
                  />
                ))}
                {new Array(8).fill(0).map((_, i) => {
                  const [x, y] = radarPoint(i, 8, 1);
                  return (
                    <line key={i} x1={RADAR_CENTER} y1={RADAR_CENTER} x2={x} y2={y} stroke="currentColor" className="text-neutral-200 dark:text-neutral-800" strokeWidth={1} />
                  );
                })}
                {panorama.map((h) => (
                  <polygon
                    key={h.habit_id}
                    points={polygonPoints(h.weekly_rates)}
                    fill={categoryColor(h.name)}
                    fillOpacity={0.15}
                    stroke={categoryColor(h.name)}
                    strokeWidth={2}
                  />
                ))}
              </svg>
              <ul className="text-xs text-neutral-600 dark:text-neutral-300 space-y-1.5">
                {panorama.map((h) => (
                  <li key={h.habit_id} className="flex items-center gap-1.5">
                    <span className="w-2 h-2 rounded-full flex-shrink-0" style={{ background: categoryColor(h.name) }} />
                    {h.name}
                  </li>
                ))}
              </ul>
            </div>
          </div>

          <div>
            <h2 className="text-sm font-medium text-neutral-500 dark:text-neutral-400 mb-3">Detalle</h2>
            <table className="w-full text-sm text-left">
              <thead className="text-xs text-neutral-400 dark:text-neutral-500">
                <tr>
                  <th className="font-medium pb-2">Hábito</th>
                  <th className="font-medium pb-2">Racha actual</th>
                  <th className="font-medium pb-2">Mejor racha</th>
                  <th className="font-medium pb-2">Esta semana</th>
                </tr>
              </thead>
              <tbody className="text-neutral-800 dark:text-neutral-200">
                {panorama.map((h) => (
                  <tr key={h.habit_id} className="border-t border-neutral-100 dark:border-neutral-800">
                    <td className="py-2">{h.name}</td>
                    <td className="py-2">{h.current_streak}</td>
                    <td className="py-2">{h.best_streak}</td>
                    <td className="py-2">{Math.round((h.weekly_rates[h.weekly_rates.length - 1] ?? 0) * 100)}%</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </>
      )}
    </div>
  );
}
