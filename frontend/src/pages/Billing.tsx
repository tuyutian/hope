import {Card, Text, Page, InlineStack, Button, BlockStack, DataTable, Box} from "@shopify/polaris";

export default function Billing() {
  return <Page title="Billings">
    <BlockStack gap="400">
    <Card>
      <InlineStack align="space-between" blockAlign="start">
        <Text as="h2" variant="headingSm" >Current protection billing </Text>
        <Button >View details</Button>
      </InlineStack>
      <BlockStack gap="100">
        <Text as="p" tone="subdued">
          Billing cycle: Feb 25, 2025 - Mar 27, 2025
        </Text>
        <Text as="p" variant="headingXl">
          $ 2,000
        </Text>
      </BlockStack>
    </Card>
    <Card padding="0">
      <BlockStack gap="300">
        <Box padding="400">
          <Text as="h2" variant="headingSm" >
            Past Protection Billing
          </Text>
        </Box>
        <DataTable columnContentTypes={["text"]} headings={["Bill number","Bill cycle","Payment status","Amount","Action"]} rows={[]} />
      </BlockStack>
    </Card>
    </BlockStack>
  </Page>;
}