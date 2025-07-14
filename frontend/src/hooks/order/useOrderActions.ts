import { useQueryClient } from '@tanstack/react-query';
import { orderQueryKeys } from './useOrderQuery.ts';
import { useCallback } from 'react';
import {OrderList} from "@/api";

export const useOrderActions = () => {
  const queryClient = useQueryClient();

  // 刷新所有订单查询
  const refreshOrders = useCallback(() => {
    queryClient.invalidateQueries({ queryKey: orderQueryKeys.all });
  }, [queryClient]);

  // 清除订单缓存
  const clearOrderCache = useCallback(() => {
    queryClient.removeQueries({ queryKey: orderQueryKeys.all });
  }, [queryClient]);

  // 预取特定页面的数据
  const prefetchOrderPage = useCallback(async (params: any) => {
    await queryClient.prefetchQuery({
      queryKey: orderQueryKeys.list(params),
      queryFn: () => OrderList(params),
      staleTime: 1000 * 60 * 5,
    });
  }, [queryClient]);

  // 乐观更新订单状态（如果需要的话）
  const updateOrderOptimistically = useCallback((orderId: string, updates: Partial<any>) => {
    queryClient.setQueriesData(
      { queryKey: orderQueryKeys.lists() },
      (oldData: any) => {
        if (!oldData) return oldData;
        
        return {
          ...oldData,
          orders: oldData.orders.map((order: any) =>
            order.id === orderId ? { ...order, ...updates } : order
          ),
        };
      }
    );
  }, [queryClient]);

  return {
    refreshOrders,
    clearOrderCache,
    prefetchOrderPage,
    updateOrderOptimistically,
  };
};
