import { ChoiceList } from '@shopify/polaris';
import { useCallback } from 'react';
import type { PaymentStatus } from '@/types/order';
import {PAYMENT_STATUS} from "@/constants/orderFilters.ts";

interface PaymentStatusFilterProps {
  paymentStatus?: PaymentStatus;
  onPaymentStatusChange: (status?: PaymentStatus) => void;
}

export const PaymentStatusFilter = ({
  paymentStatus,
  onPaymentStatusChange,
}: PaymentStatusFilterProps) => {
  const paymentStatusChoices = [
    { label: 'All Payment Status', value: '' },
    { label: 'Paid', value: PAYMENT_STATUS.PAID },
    { label: 'Refunded', value: PAYMENT_STATUS.REFUNDED },
    { label: 'Partial Refunded', value: PAYMENT_STATUS.PARTIAL_REFUNDED },
  ];

  const handlePaymentStatusChange = useCallback((value: string[]) => {
    const status = value[0];
    onPaymentStatusChange(status ? (status as PaymentStatus) : undefined);
  }, [onPaymentStatusChange]);

  return (
    <ChoiceList
      title="Payment Status"
      choices={paymentStatusChoices}
      selected={[paymentStatus || '']}
      onChange={handlePaymentStatusChange}
    />
  );
};
