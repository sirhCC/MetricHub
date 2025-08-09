import { useEffect, useState } from 'react';

export function DarkModeToggle() {
  const [enabled, setEnabled] = useState(false);

  useEffect(() => {
    const root = document.documentElement;
    if (enabled) root.classList.add('dark'); else root.classList.remove('dark');
  }, [enabled]);

  return (
    <button
      onClick={() => setEnabled(e => !e)}
      className="group relative inline-flex items-center gap-2 rounded-xl border border-gray-200 dark:border-gray-700 bg-white/70 dark:bg-gray-900/60 px-3 py-2 text-xs font-medium text-gray-600 dark:text-gray-300 shadow-sm hover:shadow transition-all"
      title="Toggle dark mode"
    >
      <span className={`w-2.5 h-2.5 rounded-full ${enabled ? 'bg-indigo-500' : 'bg-gray-300 dark:bg-gray-600'}`}></span>
      {enabled ? 'Dark' : 'Light'}
    </button>
  );
}
