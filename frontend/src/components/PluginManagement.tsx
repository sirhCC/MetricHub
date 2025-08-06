import { useState } from 'react';
import { Settings, CheckCircle, XCircle, AlertTriangle, Plus } from 'lucide-react';

interface Plugin {
  id: string;
  name: string;
  description: string;
  status: 'active' | 'inactive' | 'error';
  version: string;
  lastSync: string;
  config?: Record<string, unknown>;
}

const mockPlugins: Plugin[] = [
  {
    id: 'github',
    name: 'GitHub Integration',
    description: 'Collect deployment and commit data from GitHub repositories',
    status: 'active',
    version: '1.2.0',
    lastSync: '2 minutes ago',
  },
  {
    id: 'jenkins',
    name: 'Jenkins CI/CD',
    description: 'Track build and deployment metrics from Jenkins',
    status: 'inactive',
    version: '1.1.5',
    lastSync: '1 hour ago',
  },
  {
    id: 'datadog',
    name: 'Datadog Monitoring',
    description: 'Import incident and performance data from Datadog',
    status: 'error',
    version: '1.0.3',
    lastSync: '5 hours ago',
  },
  {
    id: 'jira',
    name: 'Jira Service Management',
    description: 'Track incident resolution and change management',
    status: 'active',
    version: '1.3.1',
    lastSync: '10 minutes ago',
  },
];

const PluginManagement = () => {
  const [plugins, setPlugins] = useState<Plugin[]>(mockPlugins);

  const getStatusIcon = (status: Plugin['status']) => {
    switch (status) {
      case 'active':
        return <CheckCircle className="h-5 w-5 text-emerald-500" />;
      case 'inactive':
        return <XCircle className="h-5 w-5 text-gray-400" />;
      case 'error':
        return <AlertTriangle className="h-5 w-5 text-red-500" />;
    }
  };

  const getStatusBadge = (status: Plugin['status']) => {
    const baseClasses = "px-3 py-1 rounded-full text-xs font-medium";
    switch (status) {
      case 'active':
        return `${baseClasses} bg-emerald-100 text-emerald-700`;
      case 'inactive':
        return `${baseClasses} bg-gray-100 text-gray-700`;
      case 'error':
        return `${baseClasses} bg-red-100 text-red-700`;
    }
  };

  const togglePlugin = (id: string) => {
    setPlugins(prev => prev.map(plugin => 
      plugin.id === id 
        ? { ...plugin, status: plugin.status === 'active' ? 'inactive' : 'active' }
        : plugin
    ));
  };

  return (
    <div className="bg-white rounded-2xl shadow-xl border-0 p-8">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h3 className="text-2xl font-bold bg-gradient-to-r from-gray-900 to-blue-800 bg-clip-text text-transparent">
            Plugin Management
          </h3>
          <p className="text-gray-600 mt-2">Manage your data collection integrations</p>
        </div>
        <button className="flex items-center space-x-2 bg-gradient-to-r from-blue-600 to-purple-600 text-white px-6 py-3 rounded-xl hover:shadow-lg transition-all duration-200 hover:-translate-y-0.5">
          <Plus className="h-5 w-5" />
          <span className="font-medium">Add Plugin</span>
        </button>
      </div>

      <div className="space-y-4">
        {plugins.map((plugin) => (
          <div
            key={plugin.id}
            className="border border-gray-200 rounded-xl p-6 hover:shadow-md transition-all duration-200 hover:border-blue-300"
          >
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-4">
                <div className="p-3 bg-gradient-to-br from-blue-500 to-purple-600 rounded-xl">
                  {getStatusIcon(plugin.status)}
                </div>
                <div>
                  <h4 className="text-lg font-semibold text-gray-900">{plugin.name}</h4>
                  <p className="text-gray-600 text-sm">{plugin.description}</p>
                  <div className="flex items-center space-x-4 mt-2">
                    <span className="text-xs text-gray-500">Version {plugin.version}</span>
                    <span className="text-xs text-gray-500">•</span>
                    <span className="text-xs text-gray-500">Last sync: {plugin.lastSync}</span>
                  </div>
                </div>
              </div>

              <div className="flex items-center space-x-4">
                <span className={getStatusBadge(plugin.status)}>
                  {plugin.status.charAt(0).toUpperCase() + plugin.status.slice(1)}
                </span>
                
                <div className="flex items-center space-x-2">
                  <button
                    onClick={() => togglePlugin(plugin.id)}
                    className={`px-4 py-2 rounded-lg font-medium transition-colors duration-200 ${
                      plugin.status === 'active'
                        ? 'bg-red-100 text-red-700 hover:bg-red-200'
                        : 'bg-emerald-100 text-emerald-700 hover:bg-emerald-200'
                    }`}
                  >
                    {plugin.status === 'active' ? 'Disable' : 'Enable'}
                  </button>
                  
                  <button className="p-2 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-colors duration-200">
                    <Settings className="h-5 w-5" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="mt-8 p-6 bg-gradient-to-r from-blue-50 to-purple-50 rounded-xl border border-blue-200">
        <div className="flex items-start space-x-3">
          <div className="p-2 bg-blue-600 rounded-lg">
            <Settings className="h-5 w-5 text-white" />
          </div>
          <div>
            <h4 className="font-semibold text-gray-900 mb-2">Need a custom integration?</h4>
            <p className="text-gray-600 text-sm mb-3">
              MetricHub supports custom plugins. Build your own integration or request one from our community.
            </p>
            <button className="text-blue-600 hover:text-blue-700 font-medium text-sm">
              Learn about plugin development →
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default PluginManagement;
