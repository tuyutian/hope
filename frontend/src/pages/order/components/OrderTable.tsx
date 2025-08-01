import { IndexTable, useIndexResourceState } from "@shopify/polaris";
import { OrderTableRow } from "./OrderTableRow";
import { OrderSkeleton } from "./OrderSkeleton";
import type { OrderItem } from "@/types/order";
import type { NonEmptyArray } from "@shopify/polaris/build/ts/src/types";
import { IndexTableHeading } from "@shopify/polaris/build/ts/src/components/IndexTable/IndexTable";

interface OrderTableProps {
  orders: OrderItem[];
  isLoading: boolean;
  isFetching?: boolean;
  itemsPerPage: number;
}

export const OrderTable = ({ orders, isLoading, isFetching = false, itemsPerPage }: OrderTableProps) => {
  const resourceName = {
    singular: "order",
    plural: "orders",
  };

  const { selectedResources, allResourcesSelected, handleSelectionChange } = useIndexResourceState(orders);

  const headings: NonEmptyArray<IndexTableHeading> = [
    { title: "Order", alignment: "center" as const },
    { title: "Date", alignment: "center" as const },
    { title: "Items", alignment: "center" as const },
    { title: "Payment status", alignment: "center" as const },
    { title: "Total", alignment: "center" as const },
    { title: "Protection Fee", alignment: "center" as const },
  ];

  const rowMarkup = isLoading ? (
    <OrderSkeleton rows={itemsPerPage} columns={6} />
  ) : (
    orders.map((order, index) => (
      <OrderTableRow key={order.id} order={order} index={index} selected={selectedResources.includes(order.id)} />
    ))
  );

  return (
    <div style={{ position: "relative" }}>
      <IndexTable
        resourceName={resourceName}
        itemCount={orders.length}
        selectedItemsCount={allResourcesSelected ? "All" : selectedResources.length}
        onSelectionChange={handleSelectionChange}
        headings={headings}
        loading={isFetching}
        hasMoreItems
      >
        {rowMarkup}
      </IndexTable>
    </div>
  );
};
