import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { AlertTriangle, Plus, RefreshCw, CheckCircle2 } from 'lucide-react';
import apiService from '../services/api';
import type { Incident, SystemState } from '../types';
import { Card } from './ui/Card';
import { Skeleton } from './ui/Skeleton';
import { useState } from 'react';

interface IncidentWidgetProps {
  className?: string;
}

export function IncidentWidget({ className }: IncidentWidgetProps) {
  const qc = useQueryClient();
  const [creating, setCreating] = useState(false);

  const { data: state, isLoading } = useQuery<SystemState>({
    queryKey: ['system-state'],
    queryFn: () => apiService.getState(),
    refetchInterval: 20000,
  });

  const incidents: Incident[] = (state?.incidents || []).sort((a,b) => a.start_time.localeCompare(b.start_time));
  const active = incidents.filter(i => !i.resolved_time);
  const resolved = incidents.filter(i => i.resolved_time);

  const createIncidentMutation = useMutation({
    mutationFn: () => {
      setCreating(true);
      return apiService.createIncident({ severity: randomSeverity(), service: randomService() });
    },
    onSettled: () => setCreating(false),
    onSuccess: () => { qc.invalidateQueries({ queryKey: ['system-state'] }); qc.invalidateQueries({ queryKey: ['dora-metrics'] }); },
  });

  const resolveMutation = useMutation({
    mutationFn: (id: string) => apiService.resolveIncident(id),
    onSuccess: () => { qc.invalidateQueries({ queryKey: ['system-state'] }); qc.invalidateQueries({ queryKey: ['dora-metrics'] }); },
  });

  const simulateDeployment = useMutation({
    mutationFn: () => apiService.createDeployment({ status: Math.random() < 0.15 ? 'failed' : 'success', service: randomService() }),
    onSuccess: () => { qc.invalidateQueries({ queryKey: ['system-state'] }); qc.invalidateQueries({ queryKey: ['dora-metrics'] }); },
  });

  return (
    <Card className={`p-6 ${className || ''}`} interactive>
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center space-x-2">
          <div className="p-2 rounded-lg bg-gradient-to-br from-red-500 to-orange-600"><AlertTriangle className="h-5 w-5 text-white" /></div>
          <h3 className="text-xl font-bold bg-gradient-to-r from-red-600 to-orange-600 bg-clip-text text-transparent">Incidents</h3>
        </div>
        <div className="flex items-center space-x-2">
          <button disabled={creating} onClick={() => createIncidentMutation.mutate()} className="flex items-center space-x-1 px-3 py-2 rounded-lg text-sm font-medium bg-red-600 hover:bg-red-700 text-white shadow disabled:opacity-60">
            <Plus className="h-4 w-4" /> <span>Simulate Incident</span>
          </button>
          <button onClick={() => simulateDeployment.mutate()} className="flex items-center space-x-1 px-3 py-2 rounded-lg text-sm font-medium bg-blue-600 hover:bg-blue-700 text-white shadow">
            <RefreshCw className="h-4 w-4" /> <span>Simulate Deployment</span>
          </button>
        </div>
      </div>
      {isLoading ? (
        <div className="space-y-3">
          {Array.from({ length: 3 }).map((_,i) => <Skeleton key={i} className="h-14 w-full" />)}
        </div>
      ) : (
        <div className="space-y-6">
          <section>
            <div className="flex items-center justify-between mb-2">
              <h4 className="text-sm font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">Active ({active.length})</h4>
              <span className="text-xs text-gray-400">Resolved {resolved.length}</span>
            </div>
            {active.length === 0 && (
              <div className="p-4 rounded-lg bg-emerald-50 dark:bg-emerald-900/20 text-emerald-700 dark:text-emerald-300 text-sm flex items-center space-x-2">
                <CheckCircle2 className="h-4 w-4" />
                <span>No active incidents</span>
              </div>
            )}
            <ul className="space-y-3">
              {active.map(inc => (
                <li key={inc.id} className="p-3 rounded-lg bg-white/70 dark:bg-gray-800/70 backdrop-blur border border-white/40 dark:border-gray-700 flex items-start justify-between">
                  <div className="pr-4">
                    <div className="flex items-center space-x-2 mb-1">
                      <span className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-semibold ${severityColors(inc.severity)}`}>{inc.severity}</span>
                      <span className="text-sm font-medium text-gray-800 dark:text-gray-200">{inc.title}</span>
                    </div>
                    <div className="text-xs text-gray-500 dark:text-gray-400 space-x-2">
                      <span>{inc.service}</span>
                      <span className="opacity-40">•</span>
                      <span>{inc.environment}</span>
                      <span className="opacity-40">•</span>
                      <span>Since {timeAgo(inc.start_time)}</span>
                    </div>
                  </div>
                  <div className="flex items-center space-x-2">
                    <button onClick={() => resolveMutation.mutate(inc.id)} className="text-xs px-2 py-1 rounded-md bg-emerald-600 hover:bg-emerald-700 text-white font-medium">Resolve</button>
                  </div>
                </li>
              ))}
            </ul>
          </section>
          {resolved.length > 0 && (
            <section>
              <h4 className="text-sm font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400 mb-2">Recently Resolved</h4>
              <ul className="space-y-2 max-h-48 overflow-auto pr-1">
                {resolved.slice(-5).reverse().map(inc => (
                  <li key={inc.id} className="p-2 rounded-md bg-gray-50 dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700 text-xs flex items-center justify-between">
                    <div className="truncate">
                      <span className={`inline-block w-2 h-2 rounded-full mr-2 ${severityDot(inc.severity)}`} />
                      <span className="font-medium text-gray-700 dark:text-gray-300 mr-2">{inc.title}</span>
                      <span className="text-gray-400">MTTR {mttr(inc)}</span>
                    </div>
                  </li>
                ))}
              </ul>
            </section>
          )}
        </div>
      )}
    </Card>
  );
}

function timeAgo(iso: string): string {
  const start = new Date(iso).getTime();
  const now = Date.now();
  const diffMs = now - start;
  const diffMin = Math.floor(diffMs / 60000);
  if (diffMin < 1) return 'just now';
  if (diffMin < 60) return `${diffMin}m ago`;
  const diffH = Math.floor(diffMin / 60);
  if (diffH < 24) return `${diffH}h ago`;
  const diffD = Math.floor(diffH / 24);
  return `${diffD}d ago`;
}

function mttr(inc: Incident): string {
  if (!inc.resolved_time) return '-';
  const start = new Date(inc.start_time).getTime();
  const end = new Date(inc.resolved_time).getTime();
  const diffMin = Math.floor((end - start)/60000);
  if (diffMin < 60) return `${diffMin}m`;
  const h = Math.floor(diffMin/60); const m = diffMin % 60;
  return `${h}h ${m}m`;
}

function severityColors(sev: string): string {
  switch (sev) {
    case 'critical': return 'bg-red-600 text-white';
    case 'high': return 'bg-orange-500 text-white';
    case 'medium': return 'bg-amber-400 text-gray-900';
    default: return 'bg-gray-300 text-gray-800';
  }
}
function severityDot(sev: string): string {
  switch (sev) {
    case 'critical': return 'bg-red-500';
    case 'high': return 'bg-orange-500';
    case 'medium': return 'bg-amber-400';
    default: return 'bg-gray-400';
  }
}

function randomSeverity(): string {
  const r = Math.random();
  if (r < 0.1) return 'critical';
  if (r < 0.3) return 'high';
  if (r < 0.65) return 'medium';
  return 'low';
}
function randomService(): string {
  const services = ['api', 'auth', 'billing', 'frontend'];
  return services[Math.floor(Math.random()*services.length)];
}

export default IncidentWidget;
