import { computed } from "vue";
import { useFetch } from "#app";

type Holiday = { summary: string };
export type HolidaysData = Record<string, Holiday>;

type ApiResponse = {
  status: "success" | "error";
  message: string;
  data: HolidaysData;
};

function pad2(n: number) {
  return String(n).padStart(2, "0");
}
function dateKey(y: number, m0: number, d: number) {
  return `${y}-${pad2(m0 + 1)}-${pad2(d)}`;
}

export function useHolidays() {
  const { token } = useAuth();

  // hit your Fiber API which now returns { status, message, data: { 'YYYY-MM-DD': { summary } } }
  const { data, error, pending, refresh } = useFetch<ApiResponse>(
    "http://localhost:8000/api/dashboard",
    {
      headers: token.value
        ? {
            Authorization: `Bearer ${token.value}`,
          }
        : undefined,
      server: false,
    },
  );

  const holidaysMap = computed<HolidaysData | null>(() => {
    if (!data.value || data.value.status !== "success") return null;
    return data.value.data;
  });

  // helpers against the map (provide currentMonth/year as args so the caller controls state)
  const isHoliday = (
    map: HolidaysData | null,
    y: number,
    m0: number,
    d: number,
  ) => {
    if (!map) return false;
    const h = map[dateKey(y, m0, d)];
    return !!(h && typeof h === "object" && "summary" in h);
  };

  const getHolidayInfo = (
    map: HolidaysData | null,
    y: number,
    m0: number,
    d: number,
  ) => {
    if (!map) return null;
    const h = map[dateKey(y, m0, d)];
    return h && typeof h === "object" && "summary" in h ? h : null;
  };

  const monthHolidays = (map: HolidaysData | null, y: number, m0: number) => {
    if (!map) return [] as { date: number; summary: string; dateStr: string }[];
    const out: { date: number; summary: string; dateStr: string }[] = [];
    for (const [k, v] of Object.entries(map)) {
      if (!/^\d{4}-\d{2}-\d{2}$/.test(k)) continue;
      const dt = new Date(k);
      if (dt.getFullYear() === y && dt.getMonth() === m0) {
        out.push({ date: dt.getDate(), summary: v.summary, dateStr: k });
      }
    }
    return out.sort((a, b) => a.date - b.date);
  };

  return {
    data, // full API payload
    error,
    pending,
    refresh,
    holidaysMap, // computed map from data.data (or null)
    isHoliday,
    getHolidayInfo,
    monthHolidays,
  };
}
