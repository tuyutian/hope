import { IndexTable, SkeletonBodyText } from "@shopify/polaris";

interface OrderSkeletonProps {
  rows?: number;
  columns?: number;
}

export const OrderSkeleton = ({ rows = 20, columns = 6 }: OrderSkeletonProps) => {
  return (
    <>
      {Array(rows).fill(0).map((row, index) => (
        <IndexTable.Row position={index} id={row.id} key={`skeleton-${index}`}>
          {Array(columns).fill(0).map((_, cellIndex) => (
            <IndexTable.Cell key={`skeleton-cell-${cellIndex}`}>
              <SkeletonBodyText lines={1} />
            </IndexTable.Cell>
          ))}
        </IndexTable.Row>
      ))}
    </>
  );
};
