import { FulfillmentCardProps } from "@/types/cart.ts";
import { Box, Card, ChoiceList, Text } from "@shopify/polaris";

export default function FulfillmentCard({ fulfillmentSettings, onFulfillmentTypeChange }: FulfillmentCardProps) {
  const handleChoiceChange = (value: string[]) => {
    onFulfillmentTypeChange(value[0]);
  };
  return (
    <Card>
      <Box paddingBlockEnd="300">
        <Text as="h2" variant="headingSm">
          Fulfillment Rules
        </Text>
      </Box>
      <ChoiceList
        title="When to mark the protection products as fulfilled"
        allowMultiple={false}
        onChange={handleChoiceChange}
        choices={fulfillmentSettings.fulfillmentOptions}
        selected={[fulfillmentSettings.fulfillmentType]}
      />
    </Card>
  );
}
