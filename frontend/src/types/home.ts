export interface OrderStatistic {
  refund: number;
  orders: number;
  total: number;
  sales: number;
}

export interface OrderStatisticsTableItem {
  date: string;
  sales: number;
  refund: number;
}

export interface DashboardResponse {
  code: number;
  message: string;
  data: {
    order_statistics: OrderStatistic;
    order_statistics_table: OrderStatisticsTableItem[];
  };
}

export interface DashboardData {
  orderStaticsTable: OrderStatisticsTableItem[];
  orderStatics: OrderStatistic;
}