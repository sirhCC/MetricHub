import type { DORAMetrics, MetricResponse, MetricData, PluginsResponse, HealthCheck, APIError, SystemState, Incident, Deployment } from '../types';

class APIService {
  private baseUrl: string;

  constructor() {
    this.baseUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';
    console.log(`🌐 API Base URL: ${this.baseUrl}`);
  }

  private async request<T>(endpoint: string, options?: RequestInit): Promise<T> {
    const url = `${this.baseUrl}/api/v1${endpoint}`;
    console.log(`🔗 Making API request to: ${url}`);
    
    try {
      const response = await fetch(url, {
        headers: {
          'Content-Type': 'application/json',
          ...options?.headers,
        },
        ...options,
      });

      console.log(`📡 Response status: ${response.status} for ${endpoint}`);

      if (!response.ok) {
        const errorData: APIError = await response.json().catch(() => ({
          error: `HTTP ${response.status}: ${response.statusText}`,
        }));
        throw new Error(errorData.error || `Request failed with status ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error(`❌ API request failed for ${endpoint}:`, error);
      throw error;
    }
  }

  // Health Check Methods
  async getHealth(): Promise<HealthCheck> {
    return this.request<HealthCheck>('/health');
  }

  async getDatabaseHealth(): Promise<HealthCheck> {
    return this.request<HealthCheck>('/health/database');
  }

  async getRedisHealth(): Promise<HealthCheck> {
    return this.request<HealthCheck>('/health/redis');
  }

  // DORA Metrics Methods
  async getDoraMetrics(): Promise<MetricResponse<DORAMetrics>> {
    return this.request<MetricResponse<DORAMetrics>>('/metrics/dora');
  }

  async getDeploymentFrequency(): Promise<MetricResponse<MetricData>> {
    return this.request<MetricResponse<MetricData>>('/metrics/dora/deployment-frequency');
  }

  async getLeadTime(): Promise<MetricResponse<MetricData>> {
    return this.request<MetricResponse<MetricData>>('/metrics/dora/lead-time');
  }

  async getMTTR(): Promise<MetricResponse<MetricData>> {
    return this.request<MetricResponse<MetricData>>('/metrics/dora/mttr');
  }

  async getChangeFailureRate(): Promise<MetricResponse<MetricData>> {
    return this.request<MetricResponse<MetricData>>('/metrics/dora/change-failure-rate');
  }

  // Plugin Methods
  async getPlugins(): Promise<PluginsResponse> {
    return this.request<PluginsResponse>('/plugins');
  }

  async getPluginHealth(pluginName: string): Promise<HealthCheck> {
    return this.request<HealthCheck>(`/plugins/${pluginName}/health`);
  }

  // Webhook Methods
  async sendWebhook(pluginName: string, payload: Record<string, unknown>): Promise<{ message: string; plugin: string; processed_at: string }> {
    return this.request(`/webhook/${pluginName}`, {
      method: 'POST',
      body: JSON.stringify(payload),
    });
  }

  // State snapshot (deployments + incidents)
  async getState(): Promise<SystemState> {
    return this.request<SystemState>('/state');
  }

  // Deployment simulation
  async createDeployment(partial?: Partial<Deployment>): Promise<{ deployment: Deployment }> {
    const payload: Partial<Deployment> = {
      id: `dep-${Date.now()}`,
      service: partial?.service || 'api',
      environment: partial?.environment || 'prod',
      status: partial?.status || 'success',
      start_time: new Date().toISOString(),
      commit_sha: Math.random().toString(16).substring(2, 9),
      ...partial,
    };
    return this.request('/deployments', { method: 'POST', body: JSON.stringify(payload) });
  }

  // Incident simulation
  async createIncident(partial?: Partial<Incident>): Promise<{ incident: Incident }> {
    const payload: Partial<Incident> = {
      id: `inc-${Date.now()}`,
      title: partial?.title || 'Synthetic Incident',
      description: partial?.description || 'Simulated incident for demo purposes',
      service: partial?.service || 'api',
      environment: partial?.environment || 'prod',
      severity: partial?.severity || 'medium',
      start_time: new Date().toISOString(),
      ...partial,
    };
    return this.request('/incidents', { method: 'POST', body: JSON.stringify(payload) });
  }

  async resolveIncident(id: string): Promise<{ resolved_at: string }> {
    return this.request(`/incidents/${id}/resolve`, { method: 'POST' });
  }
}

// Export singleton instance
export const apiService = new APIService();
export default apiService;
