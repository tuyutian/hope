import { IndexTable, Text } from "@shopify/polaris";
import { OrderBadge } from "./OrderBadge.tsx";

interface OrderItem {
  id: string;
  order_name: string;
  date: string;
  item: number;
  paymentStatus: string;
  total: string;
  protectionFee: string;
}

interface OrderTableRowProps {
  order: OrderItem;
  index: number;
  selected: boolean;
}

export const OrderTableRow = ({ order, index, selected }: OrderTableRowProps) => {
  const { id, order_name, date, item, paymentStatus, total, protectionFee } = order;

  return (
    <IndexTable.Row
      id={id}
      key={id}
      selected={selected}
      position={index}
    >
      <IndexTable.Cell>
        <Text as="span" alignment="center">{order_name}</Text>
      </IndexTable.Cell>
      <IndexTable.Cell>
        <Text as="span" alignment="center">{date}</Text>
      </IndexTable.Cell>
      <IndexTable.Cell>
        <Text as="span" alignment="center" numeric>{item}</Text>
      </IndexTable.Cell>
      <IndexTable.Cell>
        <Text as="span" alignment="center" numeric>
          <OrderBadge status={paymentStatus} />
        </Text>
      </IndexTable.Cell>
      <IndexTable.Cell >
        <Text as="span" alignment="center" numeric>{total}</Text>
      </IndexTable.Cell>
      <IndexTable.Cell>
        <Text as="span" alignment="center" numeric>{protectionFee}</Text>
      </IndexTable.Cell>
    </IndexTable.Row>
  );
};
