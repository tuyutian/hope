import { BaseApiService } from "../base";
import type { ApiResponse } from "@/types/api.ts";

export interface UserConfig {
  money_symbol: string;
  has_subscribe: boolean;
}

export interface UpdateUserParams {
  name: string;
  value: string;
}

export interface UpdateGuideParams {
  name: string;
  open: boolean;
}

export class UserService extends BaseApiService {
  constructor() {
    super("v1/user/");
  }

  // 获取用户配置
  getConfig(): Promise<ApiResponse<UserConfig>> {
    return this.get("conf");
  }

  // 更新用户设置
  updateSetting(params: UpdateUserParams): Promise<ApiResponse> {
    return this.post("setting", params);
  }

  // 更新引导步骤
  updateGuide(params: UpdateGuideParams): Promise<ApiResponse> {
    return this.post("step", params);
  }

  // 开始订阅
  startSubscription(): Promise<ApiResponse> {
    return this.get("subscribe");
  }

  // 获取会话数据
  getSessionData(): Promise<ApiResponse> {
    return this.get("session");
  }
}

// 导出实例
export const userService = new UserService();
