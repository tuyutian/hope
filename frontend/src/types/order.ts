import {
  FULFILLMENT_STATUS,
  PAYMENT_STATUS,
  SORT_DIRECTION,
  SORT_OPTIONS,
  TIME_RANGES
} from "@/constants/orderFilters.ts";
export type FulfillmentStatus = typeof FULFILLMENT_STATUS[keyof typeof FULFILLMENT_STATUS];

export type SortOption = typeof SORT_OPTIONS[keyof typeof SORT_OPTIONS];
export type SortDirection = typeof SORT_DIRECTION[keyof typeof SORT_DIRECTION];
export type TimeRange = typeof TIME_RANGES[keyof typeof TIME_RANGES];
export type PaymentStatus = typeof PAYMENT_STATUS[keyof typeof PAYMENT_STATUS];

export interface OrderAPIResponse {
  id: string;
  order_name: string;
  order_completion_at: number;
  sku_num: number;
  financial_status: string;
  fulfillment_status: string;
  total_price_amount: string;
  currency: string;
  insurance_amount: string;
  payment_date: number;
}



export interface OrderListParams {
  page: number;
  size: number;
  type: number;
  query?: string;
  time_range?: TimeRange;
  start_date?: string;
  end_date?: string;
  payment_status?: PaymentStatus;
  fulfillment_status?: FulfillmentStatus;
  sort_by?: SortOption;
  sort_direction?: SortDirection;
}

export interface OrderListResponse {
  code: number;
  data: {
    list: OrderAPIResponse[];
    total: number;
  };
  message: string;
}

export interface UseOrderQueryData {
  list: OrderItem[];
  total: number;
}

export interface OrderItem {
  id: string;
  order_name: string;
  date: string;
  item: number;
  paymentStatus: string;
  fulfillmentStatus: string;
  total: string;
  protectionFee: string;
  paymentDate: string;
  protectionFeeAmount: number;
  [key: string]: unknown; // 添加索引签名以支持 useIndexResourceState
}

export interface OrderAPIResponse {
  id: string;
  order_name: string;
  order_completion_at: number;
  sku_num: number;
  financial_status: string;
  fulfillment_status: string;
  total_price_amount: string;
  currency: string;
  insurance_amount: string;
  payment_date: number;
}


// 筛选参数接口
export interface OrderFilters {
  timeRange: TimeRange;
  customStartDate?: string;
  customEndDate?: string;
  searchQuery: string;
  paymentStatus?: PaymentStatus;
  fulfillmentStatus?: FulfillmentStatus;
  sortBy?: SortOption;
  sortDirection?: SortDirection;
}

