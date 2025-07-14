import { useCallback, useEffect, useMemo, useState } from "react";
import { useSearchParams } from "react-router";
import debounce from "lodash.debounce";
import { getDateRange } from "@/utils/dateUtils";
import type { FulfillmentStatus, PaymentStatus, SortDirection, SortOption, TimeRange } from "@/types/order";
import { SORT_DIRECTION, SORT_OPTIONS, TIME_RANGES } from "@/constants/orderFilters.ts";

export const ORDER_FILTER_TYPES = [0, 1, 2] as const;
export type OrderFilterType = (typeof ORDER_FILTER_TYPES)[number];

interface UseOrderFiltersReturn {
  // 基础状态
  selectedTab: number;
  currentPage: number;
  queryValue: string;
  debouncedQuery: string;

  // 高级筛选状态
  timeRange: TimeRange;
  customStartDate?: string;
  customEndDate?: string;
  paymentStatus?: PaymentStatus;
  fulfillmentStatus?: FulfillmentStatus;
  sortBy: SortOption;
  sortDirection: SortDirection;

  // 操作函数
  setSelectedTab: (index: number) => void;
  setCurrentPage: (page: number) => void;
  setQueryValue: (value: string) => void;
  setTimeRange: (range: TimeRange, startDate?: string, endDate?: string) => void;
  setPaymentStatus: (status?: PaymentStatus) => void;
  setFulfillmentStatus: (status?: FulfillmentStatus) => void;
  setSortOptions: (sortBy: SortOption, sortDirection: SortDirection) => void;
  clearQuery: () => void;
  clearAllFilters: () => void;

  // 计算值
  filterType: OrderFilterType;
  queryParams: {
    page: number;
    size: number;
    type: OrderFilterType;
    query: string;
    time_range: TimeRange;
    start_date?: string;
    end_date?: string;
    payment_status?: PaymentStatus;
    fulfillment_status?: FulfillmentStatus;
    sort_by?: SortOption;
    sort_direction?: SortDirection;
  };
}

export const useOrderFilters = (itemsPerPage: number = 20): UseOrderFiltersReturn => {
  const [searchParams, setSearchParams] = useSearchParams();

  // 从URL初始化基础状态
  const [selectedTab, setSelectedTabState] = useState(() => {
    const tabFromUrl = searchParams.get("tab");
    return tabFromUrl ? parseInt(tabFromUrl) : 0;
  });

  const [currentPage, setCurrentPageState] = useState(() => {
    const pageFromUrl = searchParams.get("page");
    return pageFromUrl ? parseInt(pageFromUrl) : 1;
  });

  const [queryValue, setQueryValueState] = useState(() => {
    return searchParams.get("query") || "";
  });

  // 高级筛选状态
  const [timeRange, setTimeRangeState] = useState<TimeRange>(() => {
    return (searchParams.get("timeRange") as TimeRange) || TIME_RANGES.LAST_30_DAYS;
  });

  const [customStartDate, setCustomStartDateState] = useState(() => {
    return searchParams.get("startDate") || undefined;
  });

  const [customEndDate, setCustomEndDateState] = useState(() => {
    return searchParams.get("endDate") || undefined;
  });

  const [paymentStatus, setPaymentStatusState] = useState<PaymentStatus | undefined>(() => {
    return (searchParams.get("paymentStatus") as PaymentStatus) || undefined;
  });

  const [fulfillmentStatus, setFulfillmentStatusState] = useState<FulfillmentStatus | undefined>(() => {
    return (searchParams.get("fulfillmentStatus") as FulfillmentStatus) || undefined;
  });

  const [sortBy, setSortByState] = useState<SortOption>(() => {
    return (searchParams.get("sortBy") as SortOption) || SORT_OPTIONS.ORDER_NUMBER;
  });

  const [sortDirection, setSortDirectionState] = useState<SortDirection>(() => {
    return (searchParams.get("sortDirection") as SortDirection) || SORT_DIRECTION.DESC;
  });

  const [debouncedQuery, setDebouncedQuery] = useState(queryValue);

  // 防抖处理搜索
  const debouncedQueryUpdate = useMemo(() => debounce((val: string) => setDebouncedQuery(val), 500), []);

  // 同步状态到URL
  useEffect(() => {
    const params = new URLSearchParams();

    if (selectedTab > 0) params.set("tab", selectedTab.toString());
    if (currentPage > 1) params.set("page", currentPage.toString());
    if (queryValue) params.set("query", queryValue);
    if (timeRange && timeRange !== TIME_RANGES.LAST_30_DAYS) params.set("timeRange", timeRange);
    if (customStartDate) params.set("startDate", customStartDate);
    if (customEndDate) params.set("endDate", customEndDate);
    if (paymentStatus) params.set("paymentStatus", paymentStatus);
    if (fulfillmentStatus) params.set("fulfillmentStatus", fulfillmentStatus);
    if (sortBy) params.set("sortBy", sortBy);
    if (sortDirection) params.set("sortDirection", sortDirection);

    setSearchParams(params, { replace: true });
  }, [
    selectedTab,
    currentPage,
    queryValue,
    timeRange,
    customStartDate,
    customEndDate,
    paymentStatus,
    fulfillmentStatus,
    sortBy,
    sortDirection,
    setSearchParams,
  ]);

  // 处理查询值变化
  useEffect(() => {
    debouncedQueryUpdate(queryValue);
  }, [queryValue, debouncedQueryUpdate]);

  // 清理防抖
  useEffect(() => {
    return () => {
      debouncedQueryUpdate.cancel();
    };
  }, [debouncedQueryUpdate]);

  // 操作函数
  const setSelectedTab = useCallback((index: number) => {
    setSelectedTabState(index);
    setCurrentPageState(1);
  }, []);

  const setCurrentPage = useCallback((page: number) => {
    setCurrentPageState(page);
  }, []);

  const setQueryValue = useCallback((value: string) => {
    setQueryValueState(value);
    setCurrentPageState(1);
  }, []);

  const setTimeRange = useCallback((range: TimeRange, startDate?: string, endDate?: string) => {
    setTimeRangeState(range);
    setCustomStartDateState(startDate);
    setCustomEndDateState(endDate);
    setCurrentPageState(1);
  }, []);

  const setPaymentStatus = useCallback((status?: PaymentStatus) => {
    setPaymentStatusState(status);
    setCurrentPageState(1);
  }, []);

  const setFulfillmentStatus = useCallback((status?: FulfillmentStatus) => {
    setFulfillmentStatusState(status);
    setCurrentPageState(1);
  }, []);

  const setSortOptions = useCallback((sortBy: SortOption, sortDirection: SortDirection) => {
    setSortByState(sortBy);
    setSortDirectionState(sortDirection);
    setCurrentPageState(1);
  }, []);

  const clearQuery = useCallback(() => {
    setQueryValueState("");
    setCurrentPageState(1);
  }, []);

  const clearAllFilters = useCallback(() => {
    setSelectedTabState(0);
    setCurrentPageState(1);
    setQueryValueState("");
    setTimeRangeState(TIME_RANGES.LAST_30_DAYS);
    setCustomStartDateState(undefined);
    setCustomEndDateState(undefined);
    setPaymentStatusState(undefined);
    setFulfillmentStatusState(undefined);
    setSortByState(SORT_OPTIONS.ORDER_NUMBER);
    setSortDirectionState(SORT_DIRECTION.DESC);
  }, []);

  // 计算值
  const filterType = ORDER_FILTER_TYPES[selectedTab];

  const queryParams = useMemo(() => {
    const dateRange =
      timeRange === TIME_RANGES.CUSTOM
        ? { startDate: customStartDate, endDate: customEndDate }
        : getDateRange(timeRange);

    return {
      page: currentPage,
      size: itemsPerPage,
      type: filterType,
      query: debouncedQuery,
      time_range: timeRange,
      start_date: dateRange.startDate,
      end_date: dateRange.endDate,
      payment_status: paymentStatus,
      fulfillment_status: fulfillmentStatus,
      sort_by: sortBy,
      sort_direction: sortDirection,
    };
  }, [
    currentPage,
    itemsPerPage,
    filterType,
    debouncedQuery,
    timeRange,
    customStartDate,
    customEndDate,
    paymentStatus,
    fulfillmentStatus,
    sortBy,
    sortDirection,
  ]);

  return {
    selectedTab,
    currentPage,
    queryValue,
    debouncedQuery,
    timeRange,
    customStartDate,
    customEndDate,
    paymentStatus,
    fulfillmentStatus,
    sortBy,
    sortDirection,
    setSelectedTab,
    setCurrentPage,
    setQueryValue,
    setTimeRange,
    setPaymentStatus,
    setFulfillmentStatus,
    setSortOptions,
    clearQuery,
    clearAllFilters,
    filterType,
    queryParams,
  };
};
