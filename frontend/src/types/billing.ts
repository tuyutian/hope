
export interface FilterParams {
  sort: string;
  page: number;
  size: number;
  minTime: number;
  maxTime: number;
}

export interface CurrentPeriodResponse {
  period_start: number;
  period_end: number;
  amount: number;
}