import { BaseApiService } from "../base";
import type { ApiResponse } from "@/types/api.ts";

export interface LoginParams {
  username: string;
  password: string;
}

export interface SessionData {
  user: any;
  token: string;
}

export class AuthService extends BaseApiService {
  constructor() {
    super("api/v1/auth/");
  }

  // 获取会话数据
  getSessionData(): Promise<ApiResponse> {
    return this.get("session");
  }
}

// 导出实例
export const authService = new AuthService();
