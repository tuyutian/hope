import { TIME_RANGES } from '@/constants/orderFilters.ts';
import {TimeRange} from "@/types/order.ts";

export const getDateRange = (timeRange: TimeRange): { startDate: string; endDate: string } => {
  const now = new Date();
  const endDate = now.toISOString().split('T')[0];
  
  let startDate: string;
  
  switch (timeRange) {
    case TIME_RANGES.TODAY:
      startDate = endDate;
      break;
    case TIME_RANGES.LAST_7_DAYS:
      startDate = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000).toISOString().split('T')[0];
      break;
    case TIME_RANGES.LAST_14_DAYS:
      startDate = new Date(now.getTime() - 14 * 24 * 60 * 60 * 1000).toISOString().split('T')[0];
      break;
    case TIME_RANGES.LAST_30_DAYS:
      startDate = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0];
      break;
    case TIME_RANGES.LAST_60_DAYS:
      startDate = new Date(now.getTime() - 60 * 24 * 60 * 60 * 1000).toISOString().split('T')[0];
      break;
    case TIME_RANGES.LAST_90_DAYS:
      startDate = new Date(now.getTime() - 90 * 24 * 60 * 60 * 1000).toISOString().split('T')[0];
      break;
    case TIME_RANGES.LAST_180_DAYS:
      startDate = new Date(now.getTime() - 180 * 24 * 60 * 60 * 1000).toISOString().split('T')[0];
      break;
    case TIME_RANGES.LAST_360_DAYS:
      startDate = new Date(now.getTime() - 360 * 24 * 60 * 60 * 1000).toISOString().split('T')[0];
      break;
    default:
      startDate = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0];
  }
  
  return { startDate, endDate };
};

export const getTimeRangeLabel = (timeRange: TimeRange): string => {
  const labels = {
    [TIME_RANGES.TODAY]: 'Today',
    [TIME_RANGES.LAST_7_DAYS]: 'Last 7 days',
    [TIME_RANGES.LAST_14_DAYS]: 'Last 14 days',
    [TIME_RANGES.LAST_30_DAYS]: 'Last 30 days',
    [TIME_RANGES.LAST_60_DAYS]: 'Last 60 days',
    [TIME_RANGES.LAST_90_DAYS]: 'Last 90 days',
    [TIME_RANGES.LAST_180_DAYS]: 'Last 180 days',
    [TIME_RANGES.LAST_360_DAYS]: 'Last 360 days',
    [TIME_RANGES.CUSTOM]: 'Custom range',
  };
  
  return labels[timeRange] || 'Last 30 days';
};
