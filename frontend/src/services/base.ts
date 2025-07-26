import request from "~/utils/request";
import type { ApiResponse, ApiEndpoint, ApiServiceConfig } from "@/types/api.ts";

// 基础API服务类
export class BaseApiService {
  private config: ApiServiceConfig;

  constructor(
    private baseUrl: string = "",
    config: ApiServiceConfig = {}
  ) {
    this.config = {
      timeout: import.meta.env.VITE_REQUEST_TIME_OUT,
      baseURL: `${import.meta.env.VITE_API_BASE_URL}/insurance/api`,
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json;charset=utf-8",
      },
      ...config, // 合并传入的配置
    };
  }

  // 通用请求方法
  private async request<T>(endpoint: ApiEndpoint, params?: any): Promise<ApiResponse<T>> {
    const url = this.baseUrl + endpoint.url;
    // 应用全局配置
    const requestConfig = {
      timeout: this.config.timeout,
      headers: this.config.headers,
      baseURL: this.config.baseURL,
    };

    switch (endpoint.method) {
      case "GET":
        return request.get(url, { params, ...requestConfig });
      case "POST":
        return request.post(url, params, requestConfig);
      case "PUT":
        return request.put(url, params, requestConfig);
      case "DELETE":
        return request.delete(url, { params, ...requestConfig });
      case "PATCH":
        return request.patch(url, params, requestConfig);
      default:
        throw new Error(`Unsupported method: ${endpoint.method}`);
    }
  }

  // GET 请求
  protected get<T>(url: string, params?: any): Promise<ApiResponse<T>> {
    return this.request<T>({ url, method: "GET" }, params);
  }

  // POST 请求
  protected post<T>(url: string, params?: any): Promise<ApiResponse<T>> {
    return this.request<T>({ url, method: "POST" }, params);
  }

  // PUT 请求
  protected put<T>(url: string, params?: any): Promise<ApiResponse<T>> {
    return this.request<T>({ url, method: "PUT" }, params);
  }

  // DELETE 请求
  protected delete<T>(url: string, params?: any): Promise<ApiResponse<T>> {
    return this.request<T>({ url, method: "DELETE" }, params);
  }

  // PATCH 请求
  protected patch<T>(url: string, params?: any): Promise<ApiResponse<T>> {
    return this.request<T>({ url, method: "PATCH" }, params);
  }
}
