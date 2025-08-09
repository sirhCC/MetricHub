import type { ReactNode } from 'react';

interface CardProps {
  children: ReactNode;
  className?: string;
  interactive?: boolean;
}

export function Card({ children, className = '', interactive }: CardProps) {
  return (
    <div className={`relative rounded-2xl border border-gray-200/60 dark:border-gray-700/60 bg-white/90 dark:bg-gray-900/70 backdrop-blur-sm shadow-sm ${interactive ? 'hover:shadow-lg hover:-translate-y-0.5 transition-all duration-300' : ''} ${className}`}>
      {children}
    </div>
  );
}

export function CardHeader({ children, className = '' }: { children: ReactNode; className?: string }) {
  return <div className={`px-5 pt-5 pb-3 ${className}`}>{children}</div>;
}

export function CardTitle({ children, className = '' }: { children: ReactNode; className?: string }) {
  return <h3 className={`text-sm font-semibold tracking-wide uppercase text-gray-500 dark:text-gray-400 ${className}`}>{children}</h3>;
}

export function CardContent({ children, className = '' }: { children: ReactNode; className?: string }) {
  return <div className={`px-5 pb-5 ${className}`}>{children}</div>;
}

export function Divider() {
  return <div className="h-px mx-5 bg-gradient-to-r from-transparent via-gray-200 dark:via-gray-700 to-transparent"/>;
}
