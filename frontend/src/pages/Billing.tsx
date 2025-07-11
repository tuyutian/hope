import {BlockStack, Button, Card, InlineStack, Page, Text,} from "@shopify/polaris";
import BillingTable from "@/pages/billing/components/BillingTable.tsx";



export default function Billing() {


  return <Page title="Billings">
    <BlockStack gap="400">
      <Card>
        <InlineStack align="space-between" blockAlign="start">
          <Text as="h2" variant="headingSm">Current protection billing </Text>
          <Button url="/billing/detail">View details</Button>
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
      <BillingTable />
    </BlockStack>
  </Page>;
}