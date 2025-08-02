import { TIME_RANGES } from '@/constants/orderFilters.ts';
import {TimeRange} from "@/types/order.ts";
import dayjs from 'dayjs';

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

/**
 * 格式化Unix时间戳（秒）为可读的日期格式
 * @param timestamp Unix时间戳（秒）
 * @param format 日期格式，默认为 "MMM DD, YYYY"
 * @returns 格式化后的日期字符串
 */
export const formatUnixTimestamp = (timestamp: number, format: string = "MMM DD, YYYY"): string => {
  if (!timestamp || timestamp === 0) {
    return '-';
  }
  
  // 将秒时间戳转换为毫秒
  return dayjs(timestamp * 1000).format(format);
};

/**
 * 格式化Unix时间戳（秒）为日期时间格式
 * @param timestamp Unix时间戳（秒）
 * @returns 格式化后的日期时间字符串
 */
export const formatUnixTimestampWithTime = (timestamp: number): string => {
  return formatUnixTimestamp(timestamp, "MMM DD, YYYY HH:mm:ss");
};

/**
 * 格式化两个Unix时间戳为日期范围
 * @param startTimestamp 开始时间戳（秒）
 * @param endTimestamp 结束时间戳（秒）
 * @param format 日期格式，默认为 "MMM DD, YYYY"
 * @returns 格式化后的日期范围字符串
 */
export const formatUnixTimestampRange = (
  startTimestamp: number, 
  endTimestamp: number, 
  format: string = "MMM DD, YYYY"
): string => {
  if (!startTimestamp || !endTimestamp || startTimestamp === 0 || endTimestamp === 0) {
    return '-';
  }
  
  const startDate = formatUnixTimestamp(startTimestamp, format);
  const endDate = formatUnixTimestamp(endTimestamp, format);
  
  return `${startDate} - ${endDate}`;
};
