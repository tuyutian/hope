import { useCallback, useMemo } from "react";
import { useOrderQuery } from "@/hooks/order/useOrderQuery.ts";
import { useOrderActions } from "@/hooks/order/useOrderActions.ts";
import { useOrderFilters } from "@/hooks/order/useOrderFilters.ts";
import { useOrderState } from "@/hooks/order/useOrderState.ts";

const ITEMS_PER_PAGE = 20;

export const useOrderPageLogic = () => {
  const filters = useOrderFilters(ITEMS_PER_PAGE);
  const state = useOrderState();
  const { prefetchOrderPage } = useOrderActions();

  // 获取订单数据
  const orderQuery = useOrderQuery(filters.queryParams, {
    placeholderData: previousData => previousData, // 使用 placeholderData 替代 keepPreviousData
    enabled: true,
  });

  // 计算分页信息
  const totalPages = Math.ceil((orderQuery.data?.total || 0) / ITEMS_PER_PAGE);

  // 预加载下一页
  const prefetchNextPage = useCallback(() => {
    if (filters.currentPage < totalPages) {
      const nextPageParams = {
        ...filters.queryParams,
        page: filters.currentPage + 1,
      };
      prefetchOrderPage(nextPageParams);
    }
  }, [filters.currentPage, filters.queryParams, totalPages, prefetchOrderPage]);

  // 处理标签切换
  const handleTabSelect = useCallback(
    (index: number) => {
      state.startTabLoading();
      filters.setSelectedTab(index);
      state.stopTabLoading();
    },
    [filters.setSelectedTab, state.startTabLoading, state.stopTabLoading]
  );

  // 分页处理
  const handlePreviousPage = useCallback(() => {
    filters.setCurrentPage(Math.max(1, filters.currentPage - 1));
  }, [filters.currentPage, filters.setCurrentPage]);

  const handleNextPage = useCallback(() => {
    filters.setCurrentPage(Math.min(totalPages, filters.currentPage + 1));
  }, [filters.currentPage, filters.setCurrentPage, totalPages]);

  // 计算加载状态
  const loadingStates = useMemo(
    () => ({
      isInitialLoading: orderQuery.isLoading && !orderQuery.data,
      isTableLoading: orderQuery.isLoading || state.isTabLoading,
      isRefreshing: orderQuery.isFetching && !orderQuery.isPlaceholderData, // 使用 isPlaceholderData 替代 isPreviousData
      isPageChanging: orderQuery.isFetching && orderQuery.isPlaceholderData,
    }),
    [
      orderQuery.isLoading,
      orderQuery.isFetching,
      orderQuery.isPlaceholderData, // 使用 isPlaceholderData 替代 isPreviousData
      orderQuery.data,
      state.isTabLoading,
    ]
  );

  return {
    // 数据
    orders: orderQuery.data?.list || [],
    total: orderQuery.data?.total || 0,
    totalPages,

    // 状态
    ...loadingStates,
    isError: orderQuery.isError,
    error: orderQuery.error,

    // 过滤器
    filters,

    // 操作
    handleTabSelect,
    handlePreviousPage,
    handleNextPage,
    prefetchNextPage,

    // 常量
    ITEMS_PER_PAGE,
  };
};
