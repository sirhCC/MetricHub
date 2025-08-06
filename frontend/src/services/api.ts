import type { DORAMetrics, MetricResponse, MetricData, PluginsResponse, HealthCheck, APIError } from '../types';

class APIService {
  private baseUrl: string;

  constructor() {
    this.baseUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';
  }

  private async request<T>(endpoint: string, options?: RequestInit): Promise<T> {
    const url = `${this.baseUrl}/api/v1${endpoint}`;
    
    try {
      const response = await fetch(url, {
        headers: {
          'Content-Type': 'application/json',
          ...options?.headers,
        },
        ...options,
      });

      if (!response.ok) {
        const errorData: APIError = await response.json().catch(() => ({
          error: `HTTP ${response.status}: ${response.statusText}`,
        }));
        throw new Error(errorData.error || `Request failed with status ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error(`API request failed for ${endpoint}:`, error);
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
  async sendWebhook(pluginName: string, payload: any): Promise<{ message: string; plugin: string; processed_at: string }> {
    return this.request(`/webhook/${pluginName}`, {
      method: 'POST',
      body: JSON.stringify(payload),
    });
  }
}

// Export singleton instance
export const apiService = new APIService();
export default apiService;
