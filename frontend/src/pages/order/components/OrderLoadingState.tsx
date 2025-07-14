import { Page, Card } from '@shopify/polaris';
import { OrderTable } from './OrderTable';

interface OrderLoadingStateProps {
  itemsPerPage: number;
}

export const OrderLoadingState = ({ itemsPerPage }: OrderLoadingStateProps) => {
  return (
    <Page title="Protection Orders">
      <Card>
        <OrderTable
          orders={[]}
          isLoading
          itemsPerPage={itemsPerPage}
        />
      </Card>
    </Page>
  );
};
