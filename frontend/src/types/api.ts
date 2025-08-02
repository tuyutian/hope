// 通用API响应类型
export interface ApiResponse<T = any> {
  code: number;
  data?: T;
  message: string;
}

// API端点配置
export interface ApiEndpoint {
  url: string;
  method: "GET" | "POST" | "PUT" | "DELETE" | "PATCH";
}


// API服务配置（支持所有 axios 配置选项）
export interface ApiServiceConfig {
  baseURL?: string;
  timeout?: number;
  headers?: Record<string, string>;
  [key: string]: any; // 支持所有 axios 配置选项
}
