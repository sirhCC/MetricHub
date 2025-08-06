import { useQuery } from '@tanstack/react-query';
import { Activity, TrendingUp, Clock, AlertTriangle, CheckCircle, XCircle, Settings, Zap, Target, BarChart3, Puzzle } from 'lucide-react';
import { useState } from 'react';
import apiService from '../services/api';
import MetricsChart from './MetricsChart';
import PluginManagement from './PluginManagement';

interface MetricCardProps {
  title: string;
  value: string | number;
  unit?: string;
  trend?: 'up' | 'down' | 'stable';
  icon: React.ReactNode;
  description: string;
  color: 'blue' | 'purple' | 'green' | 'orange';
}

function MetricCard({ title, value, unit, trend, icon, description, color }: MetricCardProps) {
  const getTrendColor = () => {
    if (!trend) return 'text-gray-400';
    switch (trend) {
      case 'up': return 'text-emerald-500';
      case 'down': return 'text-red-500';
      case 'stable': return 'text-blue-500';
    }
  };

  const getTrendSymbol = () => {
    if (!trend) return '';
    switch (trend) {
      case 'up': return '↗';
      case 'down': return '↘';
      case 'stable': return '→';
    }
  };

  const getCardStyles = () => {
    switch (color) {
      case 'blue':
        return {
          gradient: 'from-blue-500 to-blue-600',
          iconBg: 'bg-blue-100',
          iconColor: 'text-blue-600',
          shadow: 'shadow-blue-100'
        };
      case 'purple':
        return {
          gradient: 'from-purple-500 to-purple-600',
          iconBg: 'bg-purple-100',
          iconColor: 'text-purple-600',
          shadow: 'shadow-purple-100'
        };
      case 'green':
        return {
          gradient: 'from-emerald-500 to-emerald-600',
          iconBg: 'bg-emerald-100',
          iconColor: 'text-emerald-600',
          shadow: 'shadow-emerald-100'
        };
      case 'orange':
        return {
          gradient: 'from-orange-500 to-orange-600',
          iconBg: 'bg-orange-100',
          iconColor: 'text-orange-600',
          shadow: 'shadow-orange-100'
        };
    }
  };

  const styles = getCardStyles();

  return (
    <div className={`relative group bg-white rounded-2xl shadow-lg border-0 p-6 hover:shadow-xl transition-all duration-300 hover:-translate-y-1 ${styles.shadow}`}>
      {/* Gradient overlay */}
      <div className={`absolute inset-0 bg-gradient-to-br ${styles.gradient} opacity-0 group-hover:opacity-5 rounded-2xl transition-opacity duration-300`} />
      
      <div className="relative">
        <div className="flex items-start justify-between mb-4">
          <div className={`p-3 ${styles.iconBg} rounded-xl ${styles.iconColor} group-hover:scale-110 transition-transform duration-300`}>
            {icon}
          </div>
          {trend && (
            <div className={`flex items-center space-x-1 ${getTrendColor()}`}>
              <span className="text-lg font-bold">{getTrendSymbol()}</span>
            </div>
          )}
        </div>
        
        <div className="space-y-2">
          <p className="text-sm font-medium text-gray-600 uppercase tracking-wide">{title}</p>
          <div className="flex items-baseline space-x-2">
            <p className="text-3xl font-bold text-gray-900">
              {value}
            </p>
            {unit && (
              <span className="text-sm font-medium text-gray-500">{unit}</span>
            )}
          </div>
          <p className="text-sm text-gray-600 leading-relaxed">{description}</p>
        </div>
      </div>
    </div>
  );
}

function Dashboard() {
  const [activeTab, setActiveTab] = useState<'overview' | 'charts' | 'plugins'>('overview');
  
  const { data: doraMetrics, isLoading, error } = useQuery({
    queryKey: ['dora-metrics'],
    queryFn: () => apiService.getDoraMetrics(),
    refetchInterval: 30000,
  });

  const { data: healthData } = useQuery({
    queryKey: ['health'],
    queryFn: () => apiService.getHealth(),
    refetchInterval: 10000,
  });

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading MetricHub Dashboard...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center bg-white p-8 rounded-lg shadow-sm border border-red-200">
          <AlertTriangle className="h-12 w-12 text-red-500 mx-auto mb-4" />
          <h2 className="text-xl font-bold text-gray-900 mb-2">Connection Error</h2>
          <p className="text-gray-600 mb-4">Unable to connect to MetricHub API</p>
          <p className="text-sm text-gray-500">
            Make sure the backend server is running on <code>http://localhost:8080</code>
          </p>
          <button 
            onClick={() => window.location.reload()} 
            className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  const metrics = doraMetrics?.data;

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-100">
      <header className="bg-white/80 backdrop-blur-sm shadow-lg border-b border-white/20 sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center space-x-4">
              <div className="flex items-center space-x-3">
                <div className="relative">
                  <div className="p-2 bg-gradient-to-br from-blue-600 to-indigo-700 rounded-xl shadow-lg">
                    <Activity className="h-6 w-6 text-white" />
                  </div>
                  <div className="absolute -top-1 -right-1 h-3 w-3 bg-emerald-500 rounded-full border-2 border-white animate-pulse"></div>
                </div>
                <div>
                  <h1 className="text-2xl font-bold bg-gradient-to-r from-gray-900 to-blue-800 bg-clip-text text-transparent">
                    MetricHub
                  </h1>
                  <p className="text-sm text-gray-500">Universal DevOps Metrics</p>
                </div>
              </div>

              {/* Navigation Tabs */}
              <nav className="flex items-center space-x-1 bg-white/60 backdrop-blur-sm rounded-xl p-1 border border-white/40 shadow-lg">
                <button
                  onClick={() => setActiveTab('overview')}
                  className={`flex items-center space-x-2 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 ${
                    activeTab === 'overview'
                      ? 'bg-white text-blue-600 shadow-md'
                      : 'text-gray-600 hover:text-gray-900 hover:bg-white/50'
                  }`}
                >
                  <Activity className="h-4 w-4" />
                  <span>Overview</span>
                </button>
                <button
                  onClick={() => setActiveTab('charts')}
                  className={`flex items-center space-x-2 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 ${
                    activeTab === 'charts'
                      ? 'bg-white text-blue-600 shadow-md'
                      : 'text-gray-600 hover:text-gray-900 hover:bg-white/50'
                  }`}
                >
                  <BarChart3 className="h-4 w-4" />
                  <span>Analytics</span>
                </button>
                <button
                  onClick={() => setActiveTab('plugins')}
                  className={`flex items-center space-x-2 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 ${
                    activeTab === 'plugins'
                      ? 'bg-white text-blue-600 shadow-md'
                      : 'text-gray-600 hover:text-gray-900 hover:bg-white/50'
                  }`}
                >
                  <Puzzle className="h-4 w-4" />
                  <span>Plugins</span>
                </button>
              </nav>
            </div>
            
            <div className="flex items-center space-x-4">
              <div className="flex items-center space-x-2 bg-white/60 backdrop-blur-sm rounded-xl px-4 py-2 border border-white/40 shadow-lg">
                <div className="relative">
                  <div 
                    className={`w-3 h-3 rounded-full ${
                      healthData?.status === 'healthy' ? 'bg-emerald-500' : 'bg-red-500'
                    }`}
                  />
                  {healthData?.status === 'healthy' && (
                    <div className="absolute inset-0 w-3 h-3 bg-emerald-500 rounded-full animate-ping opacity-75"></div>
                  )}
                </div>
                <span className="text-sm font-medium text-gray-700">
                  {healthData?.status === 'healthy' ? 'System Healthy' : 'System Error'}
                </span>
              </div>
              
              <select className="text-sm border border-white/40 rounded-xl px-4 py-2 bg-white/60 backdrop-blur-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all duration-200 shadow-lg">
                <option value="last-30-days">Last 30 days</option>
                <option value="last-7-days">Last 7 days</option>
                <option value="last-90-days">Last 90 days</option>
              </select>

              <button className="p-2 text-gray-600 hover:text-gray-900 hover:bg-white/60 rounded-xl transition-colors duration-200 shadow-lg">
                <Settings className="h-5 w-5" />
              </button>
            </div>
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {activeTab === 'overview' && (
          <>
            {/* Hero Section */}
            <div className="text-center mb-12">
              <div className="relative">
                <h2 className="text-5xl font-bold bg-gradient-to-r from-blue-600 via-purple-600 to-indigo-600 bg-clip-text text-transparent mb-4">
                  DevOps Performance Dashboard
                </h2>
                <div className="absolute -inset-1 bg-gradient-to-r from-blue-600 via-purple-600 to-indigo-600 rounded-lg blur opacity-10"></div>
              </div>
              <p className="text-xl text-gray-600 max-w-3xl mx-auto leading-relaxed">
                Track your team's DORA metrics and accelerate software delivery with data-driven insights
              </p>
            </div>

            {/* DORA Metrics Grid */}
            <div className="mb-16">
              <div className="flex items-center justify-between mb-8">
                <div>
                  <h3 className="text-3xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent mb-2">
                    DORA Metrics
                  </h3>
                  <p className="text-gray-600">Four key metrics that indicate software delivery performance</p>
                </div>
                <div className="flex items-center space-x-2 bg-white/60 backdrop-blur-sm rounded-xl px-4 py-2 shadow-lg">
                  <div className="w-2 h-2 bg-emerald-500 rounded-full animate-pulse"></div>
                  <span className="text-sm font-medium text-gray-700">Live Data</span>
                </div>
              </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
            <MetricCard
              title="Deployment Frequency"
              value={metrics?.deployment_frequency || 0}
              unit="per day"
              trend="up"
              icon={<TrendingUp className="h-6 w-6" />}
              description="How often you deploy to production"
              color="blue"
            />
            
            <MetricCard
              title="Lead Time for Changes"
              value={metrics?.lead_time || "N/A"}
              trend="down"
              icon={<Clock className="h-6 w-6" />}
              description="Time from commit to production"
              color="purple"
            />
            
            <MetricCard
              title="Mean Time to Recovery"
              value={metrics?.mttr || "N/A"}
              trend="stable"
              icon={<Zap className="h-6 w-6" />}
              description="Average time to recover from incidents"
              color="green"
            />
            
            <MetricCard
              title="Change Failure Rate"
              value={metrics?.change_failure_rate ? `${(metrics.change_failure_rate * 100).toFixed(1)}%` : "N/A"}
              trend="down"
              icon={<Target className="h-6 w-6" />}
              description="Percentage of deployments causing failures"
              color="orange"
            />
          </div>
        </div>

        {/* Information Cards */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-16">
          {/* Getting Started */}
          <div className="relative group">
            <div className="absolute -inset-1 bg-gradient-to-r from-blue-500 to-purple-600 rounded-2xl blur opacity-20 group-hover:opacity-30 transition duration-1000"></div>
            <div className="relative bg-white rounded-2xl shadow-xl border-0 p-8">
              <div className="flex items-center space-x-3 mb-6">
                <div className="p-3 bg-gradient-to-br from-blue-500 to-purple-600 rounded-xl">
                  <Activity className="h-6 w-6 text-white" />
                </div>
                <h4 className="text-2xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
                  Getting Started
                </h4>
              </div>
              <div className="space-y-4">
                <div className="flex items-center space-x-4 group/item hover:bg-blue-50 rounded-xl p-3 transition-colors">
                  <div className="w-10 h-10 bg-gradient-to-br from-blue-500 to-blue-600 rounded-full flex items-center justify-center group-hover/item:scale-110 transition-transform">
                    <span className="text-sm font-bold text-white">1</span>
                  </div>
                  <span className="text-gray-700 font-medium">Configure your CI/CD integrations</span>
                </div>
                <div className="flex items-center space-x-4 group/item hover:bg-purple-50 rounded-xl p-3 transition-colors">
                  <div className="w-10 h-10 bg-gradient-to-br from-purple-500 to-purple-600 rounded-full flex items-center justify-center group-hover/item:scale-110 transition-transform">
                    <span className="text-sm font-bold text-white">2</span>
                  </div>
                  <span className="text-gray-700 font-medium">Set up incident management webhooks</span>
                </div>
                <div className="flex items-center space-x-4 group/item hover:bg-emerald-50 rounded-xl p-3 transition-colors">
                  <div className="w-10 h-10 bg-gradient-to-br from-emerald-500 to-emerald-600 rounded-full flex items-center justify-center group-hover/item:scale-110 transition-transform">
                    <span className="text-sm font-bold text-white">3</span>
                  </div>
                  <span className="text-gray-700 font-medium">Start collecting DORA metrics</span>
                </div>
              </div>
            </div>
          </div>

          {/* System Status */}
          <div className="relative group">
            <div className="absolute -inset-1 bg-gradient-to-r from-emerald-500 to-blue-600 rounded-2xl blur opacity-20 group-hover:opacity-30 transition duration-1000"></div>
            <div className="relative bg-white rounded-2xl shadow-xl border-0 p-8">
              <div className="flex items-center space-x-3 mb-6">
                <div className="p-3 bg-gradient-to-br from-emerald-500 to-blue-600 rounded-xl">
                  <CheckCircle className="h-6 w-6 text-white" />
                </div>
                <h4 className="text-2xl font-bold bg-gradient-to-r from-emerald-600 to-blue-600 bg-clip-text text-transparent">
                  System Status
                </h4>
              </div>
              <div className="space-y-4">
                <div className="flex items-center justify-between p-3 hover:bg-emerald-50 rounded-xl transition-colors group/status">
                  <span className="text-gray-700 font-medium">API Server</span>
                  <div className="flex items-center space-x-2">
                    <div className="relative">
                      <CheckCircle className="h-5 w-5 text-emerald-500" />
                      <div className="absolute inset-0 rounded-full bg-emerald-500 opacity-20 animate-ping"></div>
                    </div>
                    <span className="text-emerald-600 font-semibold">Healthy</span>
                  </div>
                </div>
                <div className="flex items-center justify-between p-3 hover:bg-emerald-50 rounded-xl transition-colors group/status">
                  <span className="text-gray-700 font-medium">Database</span>
                  <div className="flex items-center space-x-2">
                    <CheckCircle className="h-5 w-5 text-emerald-500" />
                    <span className="text-emerald-600 font-semibold">Connected</span>
                  </div>
                </div>
                <div className="flex items-center justify-between p-3 hover:bg-emerald-50 rounded-xl transition-colors group/status">
                  <span className="text-gray-700 font-medium">Cache Layer</span>
                  <div className="flex items-center space-x-2">
                    <CheckCircle className="h-5 w-5 text-emerald-500" />
                    <span className="text-emerald-600 font-semibold">Active</span>
                  </div>
                </div>
                <div className="flex items-center justify-between p-3 hover:bg-yellow-50 rounded-xl transition-colors group/status">
                  <span className="text-gray-700 font-medium">Plugin Health</span>
                  <div className="flex items-center space-x-2">
                    <XCircle className="h-5 w-5 text-yellow-500" />
                    <span className="text-yellow-600 font-semibold">1 of 3 active</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="text-center">
          <div className="relative group inline-block">
            <div className="absolute -inset-1 bg-gradient-to-r from-blue-600 via-purple-600 to-indigo-600 rounded-2xl blur opacity-20 group-hover:opacity-30 transition duration-1000"></div>
            <div className="relative bg-white/80 backdrop-blur-sm rounded-2xl px-8 py-4 border border-white/40 shadow-xl">
              <div className="flex items-center space-x-6">
                <div className="flex items-center space-x-3">
                  <div className="p-2 bg-gradient-to-br from-blue-600 to-indigo-700 rounded-lg">
                    <Activity className="h-5 w-5 text-white" />
                  </div>
                  <span className="text-lg font-bold bg-gradient-to-r from-gray-900 to-blue-800 bg-clip-text text-transparent">
                    MetricHub v0.1.0
                  </span>
                </div>
                <div className="w-px h-6 bg-gradient-to-b from-gray-300 to-gray-400"></div>
                <span className="text-gray-600 font-medium">Universal DevOps Metrics Collector</span>
              </div>
            </div>
          </div>
          <p className="text-gray-500 mt-6 font-medium">
            Built with <span className="text-red-500">❤️</span> for high-performing development teams
          </p>
        </div>
        </>
        )}

        {activeTab === 'charts' && (
          <div className="space-y-8">
            <div className="text-center mb-12">
              <h2 className="text-4xl font-bold bg-gradient-to-r from-blue-600 via-purple-600 to-indigo-600 bg-clip-text text-transparent mb-4">
                DORA Metrics Analytics
              </h2>
              <p className="text-xl text-gray-600 max-w-3xl mx-auto">
                Visualize trends and patterns in your software delivery performance
              </p>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
              <MetricsChart
                metric="deploymentFreq"
                title="Deployment Frequency Trend"
                color="#3B82F6"
                gradientId="deploymentGradient"
              />
              <MetricsChart
                metric="leadTime"
                title="Lead Time Trend"
                color="#8B5CF6"
                gradientId="leadTimeGradient"
              />
              <MetricsChart
                metric="mttr"
                title="Mean Time to Recovery"
                color="#10B981"
                gradientId="mttrGradient"
              />
              <MetricsChart
                metric="changeFailureRate"
                title="Change Failure Rate"
                color="#F59E0B"
                gradientId="failureGradient"
              />
            </div>
          </div>
        )}

        {activeTab === 'plugins' && (
          <div className="space-y-8">
            <div className="text-center mb-12">
              <h2 className="text-4xl font-bold bg-gradient-to-r from-blue-600 via-purple-600 to-indigo-600 bg-clip-text text-transparent mb-4">
                Plugin Management
              </h2>
              <p className="text-xl text-gray-600 max-w-3xl mx-auto">
                Configure and manage your data collection integrations
              </p>
            </div>

            <PluginManagement />
          </div>
        )}
      </main>
    </div>
  );
}

export default Dashboard;
