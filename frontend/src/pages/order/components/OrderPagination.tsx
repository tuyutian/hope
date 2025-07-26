import { Pagination, Text } from "@shopify/polaris";

interface OrderPaginationProps {
  total: number;
  currentPage: number;
  totalPages: number;
  onPrevious: () => void;
  onNext: () => void;
  disabled?: boolean;
}

export const OrderPagination = ({
  total,
  currentPage,
  totalPages,
  onPrevious,
  onNext,
}: OrderPaginationProps) => {
  return (
    <div style={{ 
      padding: '1rem', 
      display: 'flex', 
      justifyContent: 'space-between', 
      alignItems: 'center' 
    }}>
      <Text variant="bodySm" as="p">
        Total: {total} orders | Page {currentPage} of {totalPages}
      </Text>
      <Pagination
        hasPrevious={currentPage > 1}
        onPrevious={onPrevious}
        hasNext={currentPage < totalPages}
        onNext={onNext}
      />
    </div>
  );
};
