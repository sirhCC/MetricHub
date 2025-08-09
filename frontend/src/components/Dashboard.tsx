import { useQuery } from '@tanstack/react-query';
import { Activity, TrendingUp, Clock, AlertTriangle, CheckCircle, XCircle, Settings, Zap, Target, BarChart3, Puzzle, Shield } from 'lucide-react';
import { useState } from 'react';
import apiService from '../services/api';
import MetricsChart from './MetricsChart';
import PluginManagement from './PluginManagement';
import { Card } from './ui/Card';
import { Stat } from './ui/Stat';
import { DarkModeToggle } from './ui/Toggle';
import { Skeleton } from './ui/Skeleton';
import IncidentWidget from './IncidentWidget';

function Dashboard() {
  const [activeTab, setActiveTab] = useState<'overview' | 'charts' | 'plugins'>('overview');
  const { data: doraMetrics, isLoading, error, refetch } = useQuery({
    queryKey: ['dora-metrics'],
    queryFn: () => apiService.getDoraMetrics(),
    refetchInterval: 30000,
  });
  const { data: healthData } = useQuery({
    queryKey: ['health'],
    queryFn: () => apiService.getHealth(),
    refetchInterval: 15000,
  });

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-50 to-indigo-100 dark:from-gray-950 dark:to-gray-900">
        <Card className="max-w-md w-full p-8 text-center">
          <AlertTriangle className="h-12 w-12 text-red-500 mx-auto mb-4" />
          <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-2">Connection Error</h2>
          <p className="text-gray-600 dark:text-gray-400 mb-4 text-sm">Unable to reach MetricHub API. Ensure backend is running on http://localhost:8080</p>
          <button onClick={() => refetch()} className="px-4 py-2 bg-blue-600 text-white rounded-lg text-sm font-medium hover:bg-blue-700 transition-colors">Retry</button>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-100 dark:from-gray-950 dark:via-gray-900 dark:to-gray-800 transition-colors">
      <header className="bg-white/80 dark:bg-gray-900/70 backdrop-blur-sm shadow-lg border-b border-white/20 dark:border-gray-800 sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center space-x-4">
              <div className="flex items-center space-x-3">
                <div className="relative">
                  <div className="p-2 bg-gradient-to-br from-blue-600 to-indigo-700 rounded-xl shadow-lg">
                    <Activity className="h-6 w-6 text-white" />
                  </div>
                  <div className="absolute -top-1 -right-1 h-3 w-3 bg-emerald-500 rounded-full border-2 border-white dark:border-gray-900 animate-pulse" />
                </div>
                <div>
                  <h1 className="text-2xl font-bold bg-gradient-to-r from-gray-900 to-blue-800 dark:from-gray-100 dark:to-indigo-300 bg-clip-text text-transparent">MetricHub</h1>
                  <p className="text-sm text-gray-500 dark:text-gray-400">Universal DevOps Metrics</p>
                </div>
              </div>
              <nav className="flex items-center space-x-1 bg-white/60 dark:bg-gray-800/60 backdrop-blur-sm rounded-xl p-1 border border-white/40 dark:border-gray-700 shadow-lg">
                {([
                  { key: 'overview', label: 'Overview', icon: Activity },
                  { key: 'charts', label: 'Analytics', icon: BarChart3 },
                  { key: 'plugins', label: 'Plugins', icon: Puzzle },
                ] as const).map(tab => (
                  <button
                    key={tab.key}
                    onClick={() => setActiveTab(tab.key)}
                    className={`flex items-center space-x-2 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 ${activeTab === tab.key ? 'bg-white dark:bg-gray-900 text-blue-600 dark:text-indigo-300 shadow-md' : 'text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-gray-100 hover:bg-white/50 dark:hover:bg-gray-900/30'}`}
                  >
                    <tab.icon className="h-4 w-4" />
                    <span>{tab.label}</span>
                  </button>
                ))}
              </nav>
            </div>
            <div className="flex items-center space-x-3">
              <div className="flex items-center space-x-2 bg-white/60 dark:bg-gray-800/60 backdrop-blur-sm rounded-xl px-4 py-2 border border-white/40 dark:border-gray-700 shadow-lg">
                <div className="relative">
                  <div className={`w-3 h-3 rounded-full ${healthData?.status === 'healthy' ? 'bg-emerald-500' : 'bg-red-500'}`} />
                  {healthData?.status === 'healthy' && <div className="absolute inset-0 w-3 h-3 bg-emerald-500 rounded-full animate-ping opacity-75" />}
                </div>
                <span className="text-sm font-medium text-gray-700 dark:text-gray-300">{healthData?.status === 'healthy' ? 'System Healthy' : 'System Error'}</span>
              </div>
              <select className="text-sm border border-white/40 dark:border-gray-700 rounded-xl px-4 py-2 bg-white/60 dark:bg-gray-800/60 backdrop-blur-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all duration-200 shadow-lg text-gray-700 dark:text-gray-200">
                <option value="last-30-days">Last 30 days</option>
                <option value="last-7-days">Last 7 days</option>
                <option value="last-90-days">Last 90 days</option>
              </select>
              <DarkModeToggle />
              <button className="p-2 text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-gray-100 hover:bg-white/60 dark:hover:bg-gray-800/60 rounded-xl transition-colors duration-200 shadow-lg" title="Settings">
                <Settings className="h-5 w-5" />
              </button>
            </div>
          </div>
        </div>
      </header>
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {activeTab === 'overview' && (
          <>
            <div className="text-center mb-12">
              <div className="relative">
                <h2 className="text-5xl font-bold bg-gradient-to-r from-blue-600 via-purple-600 to-indigo-600 bg-clip-text text-transparent mb-4">DevOps Performance Dashboard</h2>
                <div className="absolute -inset-1 bg-gradient-to-r from-blue-600 via-purple-600 to-indigo-600 rounded-lg blur opacity-10" />
              </div>
              <p className="text-xl text-gray-600 dark:text-gray-400 max-w-3xl mx-auto leading-relaxed">Track your team's DORA metrics and accelerate software delivery with data-driven insights</p>
            </div>
              <div className="mb-16">
              <div className="flex items-center justify-between mb-8">
                <div>
                  <h3 className="text-3xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 dark:from-gray-100 dark:to-gray-300 bg-clip-text text-transparent mb-2">DORA Metrics</h3>
                  <p className="text-gray-600 dark:text-gray-400">Four key metrics that indicate software delivery performance</p>
                </div>
                <div className="flex items-center space-x-2 bg-white/60 dark:bg-gray-800/60 backdrop-blur-sm rounded-xl px-4 py-2 shadow-lg border border-white/40 dark:border-gray-700">
                  <div className="w-2 h-2 bg-emerald-500 rounded-full animate-pulse" />
                  <span className="text-sm font-medium text-gray-700 dark:text-gray-300">Live Data</span>
                </div>
              </div>
                {doraMetrics?.data?.overall_performance && (
                  <div className="flex flex-wrap items-center gap-2 mb-6">
                    <span className="inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold bg-gradient-to-r from-indigo-600 to-blue-600 text-white shadow">
                      <Shield className="h-3.5 w-3.5 mr-1" /> Overall {capitalize(doraMetrics.data.overall_performance)} Performance
                    </span>
                    {doraMetrics.data.classification && Object.entries(doraMetrics.data.classification).filter(([k])=>k!=='overall').map(([k,v]) => (
                      <span key={k} className={`inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium ${badgeColor(v as string)}`}>{labelize(k)}: {capitalize(v as string)}</span>
                    ))}
                  </div>
                )}
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                {isLoading ? (
                  Array.from({ length: 4 }).map((_, i) => (
                    <Card key={i} className="p-5">
                      <Skeleton className="h-5 w-20 mb-4" />
                      <Skeleton className="h-10 w-28 mb-2" />
                      <Skeleton className="h-3 w-full" />
                    </Card>
                  ))
                ) : (
                  <>
                    <Card className="p-5" interactive>
                      <Stat label="Deployment Frequency" value={doraMetrics?.data.deployment_frequency || 0} delta="+3%" deltaType="up" icon={<TrendingUp className="h-5 w-5" />} />
                      <p className="mt-3 text-sm text-gray-600 dark:text-gray-400">Deployments per day</p>
                    </Card>
                    <Card className="p-5" interactive>
                      <Stat label="Lead Time" value={doraMetrics?.data.lead_time || 'N/A'} delta="-5%" deltaType="up" icon={<Clock className="h-5 w-5" />} />
                      <p className="mt-3 text-sm text-gray-600 dark:text-gray-400">Time from commit to production</p>
                    </Card>
                    <Card className="p-5" interactive>
                      <Stat label="MTTR" value={doraMetrics?.data.mttr || 'N/A'} delta="stable" deltaType="neutral" icon={<Zap className="h-5 w-5" />} />
                      <p className="mt-3 text-sm text-gray-600 dark:text-gray-400">Mean time to recovery</p>
                    </Card>
                    <Card className="p-5" interactive>
                      <Stat label="Change Failure Rate" value={doraMetrics?.data.change_failure_rate ? `${(doraMetrics.data.change_failure_rate * 100).toFixed(1)}%` : 'N/A'} delta="-1.2%" deltaType="down" icon={<Target className="h-5 w-5" />} />
                      <p className="mt-3 text-sm text-gray-600 dark:text-gray-400">Percentage of failed changes</p>
                    </Card>
                  </>
                )}
              </div>
            </div>
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mb-16">
              <Card interactive className="p-8">
                <div className="flex items-center space-x-3 mb-6">
                  <div className="p-3 bg-gradient-to-br from-blue-500 to-purple-600 rounded-xl"><Activity className="h-6 w-6 text-white" /></div>
                  <h4 className="text-2xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">Getting Started</h4>
                </div>
                <div className="space-y-4">
                  {[
                    { step: 1, text: 'Configure your CI/CD integrations', color: 'from-blue-500 to-blue-600' },
                    { step: 2, text: 'Set up incident management webhooks', color: 'from-purple-500 to-purple-600' },
                    { step: 3, text: 'Start collecting DORA metrics', color: 'from-emerald-500 to-emerald-600' },
                  ].map(item => (
                    <div key={item.step} className="flex items-center space-x-4 group/item hover:bg-blue-50 dark:hover:bg-blue-950/40 rounded-xl p-3 transition-colors">
                      <div className={`w-10 h-10 bg-gradient-to-br ${item.color} rounded-full flex items-center justify-center group-hover/item:scale-110 transition-transform`}>
                        <span className="text-sm font-bold text-white">{item.step}</span>
                      </div>
                      <span className="text-gray-700 dark:text-gray-300 font-medium">{item.text}</span>
                    </div>
                  ))}
                </div>
              </Card>
              <Card interactive className="p-8 lg:col-span-1">
                <div className="flex items-center space-x-3 mb-6">
                  <div className="p-3 bg-gradient-to-br from-emerald-500 to-blue-600 rounded-xl"><CheckCircle className="h-6 w-6 text-white" /></div>
                  <h4 className="text-2xl font-bold bg-gradient-to-r from-emerald-600 to-blue-600 bg-clip-text text-transparent">System Status</h4>
                </div>
                <div className="space-y-2">
                  {[
                    { label: 'API Server', status: 'Healthy', icon: CheckCircle, color: 'emerald', ping: true },
                    { label: 'Database', status: 'Connected', icon: CheckCircle, color: 'emerald' },
                    { label: 'Cache Layer', status: 'Active', icon: CheckCircle, color: 'emerald' },
                    { label: 'Plugin Health', status: '1 of 3 active', icon: XCircle, color: 'yellow' },
                  ].map((s, i) => (
                    <div key={i} className={`flex items-center justify-between p-3 rounded-xl transition-colors group/status ${s.color === 'yellow' ? 'hover:bg-yellow-50 dark:hover:bg-yellow-950/30' : 'hover:bg-emerald-50 dark:hover:bg-emerald-950/40'}`}>
                      <span className="text-gray-700 dark:text-gray-300 font-medium">{s.label}</span>
                      <div className="flex items-center space-x-2">
                        <div className="relative">
                          <s.icon className={`h-5 w-5 ${s.color === 'yellow' ? 'text-yellow-500' : 'text-emerald-500'}`} />
                          {s.ping && <div className="absolute inset-0 rounded-full bg-emerald-500 opacity-20 animate-ping" />}
                        </div>
                        <span className={`font-semibold ${s.color === 'yellow' ? 'text-yellow-600' : 'text-emerald-600'}`}>{s.status}</span>
                      </div>
                    </div>
                  ))}
                </div>
              </Card>
              <IncidentWidget className="lg:col-span-1" />
            </div>
            <div className="text-center">
              <div className="relative group inline-block">
                <div className="absolute -inset-1 bg-gradient-to-r from-blue-600 via-purple-600 to-indigo-600 rounded-2xl blur opacity-20 group-hover:opacity-30 transition duration-1000" />
                <div className="relative bg-white/80 dark:bg-gray-900/70 backdrop-blur-sm rounded-2xl px-8 py-4 border border-white/40 dark:border-gray-700 shadow-xl">
                  <div className="flex items-center space-x-6">
                    <div className="flex items-center space-x-3">
                      <div className="p-2 bg-gradient-to-br from-blue-600 to-indigo-700 rounded-lg"><Activity className="h-5 w-5 text-white" /></div>
                      <span className="text-lg font-bold bg-gradient-to-r from-gray-900 to-blue-800 dark:from-gray-100 dark:to-indigo-300 bg-clip-text text-transparent">MetricHub v0.1.0</span>
                    </div>
                    <div className="w-px h-6 bg-gradient-to-b from-gray-300 to-gray-400 dark:from-gray-700 dark:to-gray-600" />
                    <span className="text-gray-600 dark:text-gray-400 font-medium">Universal DevOps Metrics Collector</span>
                  </div>
                </div>
              </div>
              <p className="text-gray-500 dark:text-gray-400 mt-6 font-medium">Built with <span className="text-red-500">❤️</span> for high-performing development teams</p>
            </div>
          </>
        )}
        {activeTab === 'charts' && (
          <div className="space-y-8">
            <div className="text-center mb-12">
              <h2 className="text-4xl font-bold bg-gradient-to-r from-blue-600 via-purple-600 to-indigo-600 bg-clip-text text-transparent mb-4">DORA Metrics Analytics</h2>
              <p className="text-xl text-gray-600 dark:text-gray-400 max-w-3xl mx-auto">Visualize trends and patterns in your software delivery performance</p>
            </div>
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
              <MetricsChart metric="deploymentFreq" title="Deployment Frequency Trend" color="#3B82F6" gradientId="deploymentGradient" />
              <MetricsChart metric="leadTime" title="Lead Time Trend" color="#8B5CF6" gradientId="leadTimeGradient" />
              <MetricsChart metric="mttr" title="Mean Time to Recovery" color="#10B981" gradientId="mttrGradient" />
              <MetricsChart metric="changeFailureRate" title="Change Failure Rate" color="#F59E0B" gradientId="failureGradient" />
            </div>
          </div>
        )}
        {activeTab === 'plugins' && (
          <div className="space-y-8">
            <div className="text-center mb-12">
              <h2 className="text-4xl font-bold bg-gradient-to-r from-blue-600 via-purple-600 to-indigo-600 bg-clip-text text-transparent mb-4">Plugin Management</h2>
              <p className="text-xl text-gray-600 dark:text-gray-400 max-w-3xl mx-auto">Configure and manage your data collection integrations</p>
            </div>
            <PluginManagement />
          </div>
        )}
      </main>
    </div>
  );
}

export default Dashboard;

// Helper presentation utilities
function capitalize(v: string) { return v.charAt(0).toUpperCase() + v.slice(1); }
function labelize(key: string) { return key.replace(/_/g,' ').replace(/\b\w/g, c => c.toUpperCase()); }
function badgeColor(level: string) {
  switch(level) {
    case 'elite': return 'bg-emerald-600 text-white';
    case 'high': return 'bg-blue-600 text-white';
    case 'medium': return 'bg-amber-500 text-gray-900';
    case 'low': return 'bg-red-600 text-white';
    default: return 'bg-gray-400 text-gray-900';
  }
}
