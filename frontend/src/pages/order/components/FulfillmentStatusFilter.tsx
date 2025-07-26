import { ChoiceList } from '@shopify/polaris';
import { useCallback } from 'react';
import { FULFILLMENT_STATUS } from '@/constants/orderFilters';
import type { FulfillmentStatus } from '@/types/order';

interface FulfillmentStatusFilterProps {
  fulfillmentStatus?: FulfillmentStatus;
  onFulfillmentStatusChange: (status?: FulfillmentStatus) => void;
}

export const FulfillmentStatusFilter = ({
  fulfillmentStatus,
  onFulfillmentStatusChange,
}: FulfillmentStatusFilterProps) => {
  const fulfillmentStatusChoices = [
    { label: 'All Fulfillment Status', value: '' },
    { label: 'Fulfilled', value: FULFILLMENT_STATUS.FULFILLED },
    { label: 'Unfulfilled', value: FULFILLMENT_STATUS.UNFULFILLED },
    { label: 'Partial Fulfilled', value: FULFILLMENT_STATUS.PARTIAL_FULFILLED },
  ];

  const handleFulfillmentStatusChange = useCallback((value: string[]) => {
    const status = value[0];
    onFulfillmentStatusChange(status ? (status as FulfillmentStatus) : undefined);
  }, [onFulfillmentStatusChange]);

  return (
    <ChoiceList
      title="Fulfillment Status"
      choices={fulfillmentStatusChoices}
      selected={[fulfillmentStatus || '']}
      onChange={handleFulfillmentStatusChange}
    />
  );
};
