import { IndexFilters, useSetIndexFiltersMode } from "@shopify/polaris";
import { useCallback } from "react";
import { TimeRangeFilter } from "./TimeRangeFilter";
import { PaymentStatusFilter } from "./PaymentStatusFilter";
import { FulfillmentStatusFilter } from "./FulfillmentStatusFilter";
import type { 
  TimeRange, 
  PaymentStatus, 
  FulfillmentStatus, 
  SortOption, 
  SortDirection 
} from "@/types/order";
import { FULL_SORT_OPTIONS } from "@/constants/orderFilters";
import type { IndexFiltersPrimaryAction } from "@shopify/polaris/build/ts/src/components/IndexFilters/types";

interface OrderFiltersProps {
  // 基础筛选
  queryValue: string;
  onQueryChange: (value: string) => void;
  onQueryClear: () => void;
  onClearAll: () => void;
  onCancel: () => void;
  selected: number;
  onSelect: (index: number) => void;
  isTabLoading: boolean;
  onPrimaryAction: () => Promise<boolean>;
  
  // 时间范围筛选
  timeRange: TimeRange;
  customStartDate?: string;
  customEndDate?: string;
  onTimeRangeChange: (range: TimeRange, startDate?: string, endDate?: string) => void;
  
  // 状态筛选
  paymentStatus?: PaymentStatus;
  fulfillmentStatus?: FulfillmentStatus;
  onPaymentStatusChange: (status?: PaymentStatus) => void;
  onFulfillmentStatusChange: (status?: FulfillmentStatus) => void;
  
  // 排序
  sortBy?: SortOption;
  sortDirection?: SortDirection;
  onSortChange: (sortBy?: SortOption, sortDirection?: SortDirection) => void;
}

export const OrderFilters = ({
  queryValue,
  onQueryChange,
  onQueryClear,
  onClearAll,
  onCancel,
  selected,
  onSelect,
  isTabLoading,
  onPrimaryAction,
  timeRange,
  customStartDate,
  customEndDate,
  onTimeRangeChange,
  paymentStatus,
  fulfillmentStatus,
  onPaymentStatusChange,
  onFulfillmentStatusChange,
  sortBy,
  sortDirection,
  onSortChange,
}: OrderFiltersProps) => {
  const { mode, setMode } = useSetIndexFiltersMode();

  const itemStrings = ['All', 'Paid', 'Refund'];
  const tabs = itemStrings.map((item, index) => ({
    content: item,
    index,
    onAction: () => {},
    id: `${item}-${index}`,
    isLocked: index === 0,
  }));

  const primaryAction: IndexFiltersPrimaryAction = {
    type: 'save',
    onAction: onPrimaryAction,
    disabled: false,
    loading: false,
  };

  // 限制搜索框最多200个字符
  const handleQueryChange = useCallback((value: string) => {
    if (value.length <= 200) {
      onQueryChange(value);
    }
  }, [onQueryChange]);

  // 配置排序选项
  const sortOptions = [
    {
      label: 'Order Number',
      value: FULL_SORT_OPTIONS.ORDER_NUMBER_ASC,
      directionLabel: 'A-Z',
    },
    {
      label: 'Order Number',
      value: FULL_SORT_OPTIONS.ORDER_NUMBER_DESC,
      directionLabel: 'Z-A',
    },
    {
      label: 'Payment Date',
      value: FULL_SORT_OPTIONS.PAYMENT_DATE_ASC,
      directionLabel: 'Oldest-Newest',
    },
    {
      label: 'Payment Date',
      value: FULL_SORT_OPTIONS.PAYMENT_DATE_DESC,
      directionLabel: 'Newest-Oldest',
    },
    {
      label: 'Protection Fee',
      value: FULL_SORT_OPTIONS.PROTECTION_FEE_ASC,
      directionLabel: 'Lowest-Highest',
    },
    {
      label: 'Protection Fee',
      value: FULL_SORT_OPTIONS.PROTECTION_FEE_DESC,
      directionLabel: 'Highest-Lowest',
    },
  ];

  // 处理排序选择
  const handleSortChange = useCallback((sortValue: string[]) => {
    if (sortValue.length === 0) {
      onSortChange(undefined, undefined);
      return;
    }

    const selectedSort = sortValue[0];
    const [sortByValue, direction] = selectedSort.split(' ');
    
    onSortChange(
      sortByValue as SortOption,
      direction as SortDirection
    );
  }, [onSortChange]);

  // 当前选中的排序值
  const selectedSortValue = sortBy && sortDirection 
    ? [`${sortBy} ${sortDirection}`]
    : [];

  const appliedFilters = [];
  
  // 添加时间范围过滤标签
  if (timeRange && timeRange !== 'last_30_days') {
    const timeLabel = timeRange === 'custom' && customStartDate && customEndDate 
      ? `${customStartDate} - ${customEndDate}`
      : timeRange.replace(/_/g, ' ');
    
    appliedFilters.push({
      key: 'timeRange',
      label: `Time: ${timeLabel}`,
      onRemove: () => onTimeRangeChange('last_30_days'),
    });
  }
  
  // 添加支付状态过滤标签
  if (paymentStatus) {
    const paymentLabel = paymentStatus.replace(/_/g, ' ').toLowerCase();
    appliedFilters.push({
      key: 'paymentStatus',
      label: `Payment: ${paymentLabel}`,
      onRemove: () => onPaymentStatusChange(undefined),
    });
  }
  
  // 添加履约状态过滤标签
  if (fulfillmentStatus) {
    const fulfillmentLabel = fulfillmentStatus.replace(/_/g, ' ').toLowerCase();
    appliedFilters.push({
      key: 'fulfillmentStatus',
      label: `Fulfillment: ${fulfillmentLabel}`,
      onRemove: () => onFulfillmentStatusChange(undefined),
    });
  }

  const filters = [
    {
      key: 'timeRange',
      label: 'Time Range',
      filter: (
        <TimeRangeFilter
          selectedRange={timeRange}
          customStartDate={customStartDate}
          customEndDate={customEndDate}
          onRangeChange={onTimeRangeChange}
        />
      ),
      shortcut: true,
    },
    {
      key: 'paymentStatus',
      label: 'Payment Status',
      filter: (
        <PaymentStatusFilter
          paymentStatus={paymentStatus}
          onPaymentStatusChange={onPaymentStatusChange}
        />
      ),
      shortcut: true,
    },
    {
      key: 'fulfillmentStatus',
      label: 'Fulfillment Status',
      filter: (
        <FulfillmentStatusFilter
          fulfillmentStatus={fulfillmentStatus}
          onFulfillmentStatusChange={onFulfillmentStatusChange}
        />
      ),
      shortcut: true,
    },
  ];

  return (
    <IndexFilters
      primaryAction={primaryAction}
      onClearAll={onClearAll}
      mode={mode}
      setMode={setMode}
      queryValue={queryValue}
      queryPlaceholder="Search orders (max 200 characters)"
      onQueryChange={handleQueryChange}
      onQueryClear={onQueryClear}
      tabs={tabs}
      selected={selected}
      onSelect={onSelect}
      cancelAction={{
        onAction: onCancel,
        disabled: false,
        loading: false,
      }}
      filters={filters}
      appliedFilters={appliedFilters}
      sortOptions={sortOptions}
      sortSelected={selectedSortValue}
      onSort={handleSortChange}
      hideQueryField={false}
      disabled={isTabLoading}
      canCreateNewView={false}
    />
  );
};