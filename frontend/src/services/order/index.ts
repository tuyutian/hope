import { BaseApiService } from "../base";
import type { ApiResponse } from "@/types/api.ts";
import type { OrderListParams, OrderListResponse } from "@/types/order";
import type { DashboardResponse } from "@/types/home";

export class OrderService extends BaseApiService {
  constructor() {
    super("api/v1/order/");
  }

  // 获取仪表盘数据 - 返回原始的 ApiResponse
  getDashboard(days: string): Promise<ApiResponse<DashboardResponse>> {
    return this.get<DashboardResponse>(`dashboard?days=${days}`);
  }

  // 获取订单列表 - 返回原始的 ApiResponse
  getList(params: OrderListParams): Promise<ApiResponse<OrderListResponse>> {
    return this.post<OrderListResponse>("list", params);
  }
}

// 导出实例
export const orderService = new OrderService();
