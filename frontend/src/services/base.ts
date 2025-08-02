import request from "~/utils/request";
import type { ApiResponse, ApiEndpoint, ApiServiceConfig } from "@/types/api.ts";
import type { AxiosRequestConfig } from "axios";

// 基础API服务类
export class BaseApiService {
  private config: ApiServiceConfig;

  constructor(
    private baseUrl: string = "",
    config: ApiServiceConfig = {}
  ) {
    this.config = {
      timeout: import.meta.env.VITE_REQUEST_TIME_OUT,
      baseURL: `${import.meta.env.VITE_API_BASE_URL}/protectify/api`,
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json;charset=utf-8",
      },
      ...config, // 合并传入的配置
    };
  }

  // 通用请求方法
  private async request<T>(endpoint: ApiEndpoint, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    const url = this.baseUrl + endpoint.url;
    
    // 合并配置，处理headers的深度合并
    const requestConfig = {
      ...this.config,
      ...config,
      headers: {
        ...this.config.headers,
        ...config?.headers,
      },
    };

    switch (endpoint.method) {
      case "GET":
        return request.get(url, { params: data, ...requestConfig });
      case "POST":
        return request.post(url, data, requestConfig);
      case "PUT":
        return request.put(url, data, requestConfig);
      case "DELETE":
        return request.delete(url, { params: data, ...requestConfig });
      case "PATCH":
        return request.patch(url, data, requestConfig);
      default:
        throw new Error(`Unsupported method: ${endpoint.method}`);
    }
  }

  // GET 请求
  protected get<T = any>(url: string, params?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    return this.request<T>({ url, method: "GET" }, params, config);
  }

  // POST 请求
  protected post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    return this.request<T>({ url, method: "POST" }, data, config);
  }

  // PUT 请求
  protected put<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    return this.request<T>({ url, method: "PUT" }, data, config);
  }

  // DELETE 请求
  protected delete<T = any>(url: string, params?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    return this.request<T>({ url, method: "DELETE" }, params, config);
  }

  // PATCH 请求
  protected patch<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    return this.request<T>({ url, method: "PATCH" }, data, config);
  }
}
