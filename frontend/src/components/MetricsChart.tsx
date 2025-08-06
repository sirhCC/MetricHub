import { XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, AreaChart, Area } from 'recharts';

// Sample data - in real app this would come from API
const generateSampleData = () => {
  const data = [];
  const now = new Date();
  
  for (let i = 29; i >= 0; i--) {
    const date = new Date(now);
    date.setDate(date.getDate() - i);
    
    data.push({
      date: date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' }),
      deploymentFreq: Math.floor(Math.random() * 10) + 1,
      leadTime: Math.floor(Math.random() * 24) + 2,
      mttr: Math.floor(Math.random() * 4) + 1,
      changeFailureRate: Math.floor(Math.random() * 15) + 2,
    });
  }
  
  return data;
};

interface MetricsChartProps {
  metric: 'deploymentFreq' | 'leadTime' | 'mttr' | 'changeFailureRate';
  title: string;
  color: string;
  gradientId: string;
}

const MetricsChart = ({ metric, title, color, gradientId }: MetricsChartProps) => {
  const data = generateSampleData();

  const getMetricLabel = () => {
    switch (metric) {
      case 'deploymentFreq': return 'Deployments per Day';
      case 'leadTime': return 'Hours';
      case 'mttr': return 'Hours';
      case 'changeFailureRate': return 'Failure %';
      default: return '';
    }
  };

  return (
    <div className="bg-white rounded-2xl shadow-xl border-0 p-6 hover:shadow-2xl transition-all duration-300">
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-xl font-bold text-gray-900">{title}</h3>
        <div className="px-3 py-1 bg-gradient-to-r from-blue-100 to-purple-100 rounded-full">
          <span className="text-sm font-medium text-gray-700">Last 30 days</span>
        </div>
      </div>
      
      <ResponsiveContainer width="100%" height={200}>
        <AreaChart data={data}>
          <defs>
            <linearGradient id={gradientId} x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor={color} stopOpacity={0.3}/>
              <stop offset="95%" stopColor={color} stopOpacity={0.05}/>
            </linearGradient>
          </defs>
          <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
          <XAxis 
            dataKey="date" 
            axisLine={false}
            tickLine={false}
            tick={{ fontSize: 12, fill: '#6b7280' }}
          />
          <YAxis 
            axisLine={false}
            tickLine={false}
            tick={{ fontSize: 12, fill: '#6b7280' }}
            label={{ value: getMetricLabel(), angle: -90, position: 'insideLeft', style: { textAnchor: 'middle', fill: '#6b7280', fontSize: 12 } }}
          />
          <Tooltip 
            contentStyle={{
              backgroundColor: 'white',
              border: 'none',
              borderRadius: '12px',
              boxShadow: '0 10px 25px rgba(0, 0, 0, 0.1)',
              fontSize: '14px'
            }}
            labelStyle={{ color: '#374151', fontWeight: 'bold' }}
          />
          <Area 
            type="monotone" 
            dataKey={metric} 
            stroke={color} 
            strokeWidth={3}
            fill={`url(#${gradientId})`}
            dot={{ fill: color, strokeWidth: 2, r: 4 }}
            activeDot={{ r: 6, stroke: color, strokeWidth: 2, fill: 'white' }}
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
};

export default MetricsChart;
