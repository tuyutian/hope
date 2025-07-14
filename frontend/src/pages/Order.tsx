import { Page, Card } from "@shopify/polaris";
import { useEffect } from 'react';
import {
  OrderFilters,
  OrderTable,
  OrderPagination,
  OrderErrorBoundary,
  OrderLoadingState,
} from "@/pages/order/components";
import { useOrderPageLogic } from "@/hooks/order/useOrderPageLogic";

export default function Order() {
  const {
    // 数据
    orders,
    total,
    totalPages,
    
    // 状态
    isInitialLoading,
    isTableLoading,
    isRefreshing,
    isPageChanging,
    isError,
    error,
    
    // 过滤器
    filters,
    
    // 操作
    handleTabSelect,
    handlePrimaryAction,
    handlePreviousPage,
    handleNextPage,
    prefetchNextPage,
    
    // 常量
    ITEMS_PER_PAGE,
  } = useOrderPageLogic();

  // 预加载下一页
  useEffect(() => {
    prefetchNextPage();
  }, [prefetchNextPage]);

  // 初始加载状态
  if (isInitialLoading) {
    return <OrderLoadingState itemsPerPage={ITEMS_PER_PAGE} />;
  }

  // 错误状态
  if (isError && error) {
    return <OrderErrorBoundary error={error} />;
  }

  return (
    <Page title="Protection Orders">
      <Card padding="0">
        <OrderFilters
          // 基础筛选
          queryValue={filters.queryValue}
          onQueryChange={filters.setQueryValue}
          onQueryClear={filters.clearQuery}
          onClearAll={filters.clearAllFilters}
          onCancel={filters.clearQuery}
          selected={filters.selectedTab}
          onSelect={handleTabSelect}
          isTabLoading={isTableLoading}
          onPrimaryAction={handlePrimaryAction}
          
          // 时间范围筛选
          timeRange={filters.timeRange}
          customStartDate={filters.customStartDate}
          customEndDate={filters.customEndDate}
          onTimeRangeChange={filters.setTimeRange}
          
          // 状态筛选
          paymentStatus={filters.paymentStatus}
          fulfillmentStatus={filters.fulfillmentStatus}
          onPaymentStatusChange={filters.setPaymentStatus}
          onFulfillmentStatusChange={filters.setFulfillmentStatus}
          
          // 排序
          sortBy={filters.sortBy}
          sortDirection={filters.sortDirection}
          onSortChange={filters.setSortOptions}
        />
        
        <OrderTable
          orders={orders}
          isLoading={isTableLoading}
          isFetching={isRefreshing}
          itemsPerPage={ITEMS_PER_PAGE}
        />
        
        {/* 页面切换时的加载提示 */}
        {isPageChanging && (
          <div style={{ 
            padding: '0.5rem', 
            textAlign: 'center', 
            backgroundColor: '#f6f6f7',
            borderTop: '1px solid #e1e3e5'
          }}>
            <small>正在加载...</small>
          </div>
        )}
        
        <OrderPagination
          total={total}
          currentPage={filters.currentPage}
          totalPages={totalPages}
          onPrevious={handlePreviousPage}
          onNext={handleNextPage}
          disabled={isTableLoading}
        />
      </Card>
    </Page>
  );
}