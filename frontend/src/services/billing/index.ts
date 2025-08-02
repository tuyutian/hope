import { BaseApiService } from "../base";
import type { ApiResponse } from "@/types/api.ts";
import type { FilterParams,CurrentPeriodResponse } from "@/types/billing";

export class BillingService extends BaseApiService {
  constructor() {
    super("v1/billing/");
  }

  // 获取账单数据
  getData(params: FilterParams): Promise<ApiResponse> {
    return this.post("list", params);
  }

  // 获取账单详情
  getDetailData(params: FilterParams): Promise<ApiResponse> {
    return this.post("details", params);
  }

  // 获取当前周期
  getCurrentPeriod(): Promise<ApiResponse<CurrentPeriodResponse>> {
    return this.get("current");
  }
}

// 导出实例
export const billingService = new BillingService();
