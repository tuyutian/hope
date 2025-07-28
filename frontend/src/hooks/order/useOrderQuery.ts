import { useQuery, UseQueryOptions, keepPreviousData } from "@tanstack/react-query";
import { formatTimestampToUSDate } from "@/utils/tools.ts";
import type {
  OrderItem,
  OrderAPIResponse,
  UseOrderQueryData,
  OrderListParams,
  OrderListResponse,
} from "@/types/order.ts";
import { type ApiResponse, orderService } from "@/api";

// 查询键生成函数
export const orderQueryKeys = {
  all: ["orders"] as const,
  lists: () => [...orderQueryKeys.all, "list"] as const,
  list: (params: OrderListParams) => [...orderQueryKeys.lists(), params] as const,
};

// 订单数据转换函数
const transformOrderData = (apiData: OrderAPIResponse[]): OrderItem[] => {
  return apiData.map(item => ({
    id: item.id,
    order_name: item.order_name,
    date: item.order_completion_at > 0 ? formatTimestampToUSDate(item.order_completion_at) : "-",
    item: item.sku_num,
    paymentStatus: item.financial_status,
    fulfillmentStatus: item.fulfillment_status, // 添加缺失的字段
    total: `${item.total_price_amount} ${item.currency}`,
    protectionFee: `${item.protectify_amount} ${item.currency}`,
    paymentDate: item.payment_date > 0 ? formatTimestampToUSDate(item.payment_date) : "-", // 添加缺失的字段
    protectionFeeAmount: parseFloat(item.protectify_amount) || 0, // 添加缺失的字段，转换为数字
  }));
};

// 订单查询 hook
export const useOrderQuery = (
  params: OrderListParams,
  options?: Omit<UseQueryOptions<ApiResponse<OrderListResponse>, Error, UseOrderQueryData>, "queryKey" | "queryFn">
) => {
  return useQuery<ApiResponse<OrderListResponse>, Error, UseOrderQueryData>({
    queryKey: orderQueryKeys.list(params),
    queryFn: () => orderService.getList(params),
    select: (data): UseOrderQueryData => ({
      list: transformOrderData(data?.data?.list || []),
      total: data?.data?.total || 0,
    }),
    placeholderData: keepPreviousData, // 替换 keepPreviousData: true
    staleTime: 1000 * 60 * 5, // 5 分钟内数据视为新鲜
    gcTime: 1000 * 60 * 10, // 10 分钟后清除缓存
    retry: 3,
    retryDelay: attemptIndex => Math.min(1000 * 2 ** attemptIndex, 30000),
    ...options,
  });
};
