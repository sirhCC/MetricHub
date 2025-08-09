import type { ReactNode } from 'react';

interface StatProps {
  label: string;
  value: string | number;
  icon?: ReactNode;
  delta?: string;
  deltaType?: 'up' | 'down' | 'neutral';
  tooltip?: string;
}

const deltaColors = {
  up: 'text-emerald-600 dark:text-emerald-400',
  down: 'text-red-600 dark:text-red-400',
  neutral: 'text-gray-500 dark:text-gray-400'
};

export function Stat({ label, value, icon, delta, deltaType = 'neutral', tooltip }: StatProps) {
  return (
    <div className="flex flex-col gap-2">
      <div className="flex items-center gap-2">
        {icon && <div className="p-2 rounded-lg bg-gray-50 dark:bg-gray-800 text-gray-600 dark:text-gray-300">{icon}</div>}
        <span className="text-xs font-medium tracking-wide uppercase text-gray-500 dark:text-gray-400" title={tooltip}>{label}</span>
      </div>
      <div className="flex items-baseline gap-2">
        <span className="text-3xl font-semibold text-gray-900 dark:text-gray-100 tabular-nums">{value}</span>
        {delta && <span className={`text-xs font-medium ${deltaColors[deltaType]}`}>{delta}</span>}
      </div>
    </div>
  );
}
