import { Page, Card, Button, Text } from '@shopify/polaris';
import { useOrderActions } from '@/hooks/order/useOrderActions.ts';

interface OrderErrorBoundaryProps {
  error: Error;
  onRetry?: () => void;
}

export const OrderErrorBoundary = ({ error, onRetry }: OrderErrorBoundaryProps) => {
  const { refreshOrders } = useOrderActions();

  const handleRetry = () => {
    if (onRetry) {
      onRetry();
    } else {
      refreshOrders();
    }
  };

  return (
    <Page title="Protection Orders">
      <Card>
        <div style={{ 
          padding: '2rem', 
          textAlign: 'center',
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          gap: '1rem'
        }}>
          <Text variant="headingMd" as="h2">
            加载订单时出错
          </Text>
          <Text as="p" variant="bodyMd" tone="subdued">
            {error?.message || '未知错误'}
          </Text>
          <Button variant="primary" onClick={handleRetry}>
            重试
          </Button>
        </div>
      </Card>
    </Page>
  );
};
