import { Badge } from "@shopify/polaris";
import type {Tone} from "@shopify/polaris/build/ts/src/components/Badge/types";

interface OrderBadgeProps {
  status: string;
}

export const OrderBadge = ({ status }: OrderBadgeProps) => {
  const getBadgeStatus = (status: string):Tone => {
    switch (status) {
      case 'PAID':
        return 'success';
      case 'PARTIALLY_PAID':
        return 'attention';
      case 'PARTIALLY_REFUNDED':
      case 'REFUNDED':
        return 'warning';
      case 'UNPAID':
      default:
        return 'critical';
    }
  };

  return <Badge tone={getBadgeStatus(status)}>
    {status.replace(/_/g, ' ')}
    </Badge>
};
