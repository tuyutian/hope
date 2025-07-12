import {BlockStack, Button, Card, InlineStack, Page, Text,} from "@shopify/polaris";
import BillingTable from "@/pages/billing/components/BillingTable.tsx";
import {useEffect, useState} from "react";
import {GetCurrentPeriod} from "@/api";
import BillingCardSkeleton from "@/pages/billing/components/BillingCardSkeleton.tsx";


export default function Billing() {
  const [currentPeriod, setCurrentPeriod] = useState<{ billingCycle: string, amount: string }>({
    billingCycle: "-",
    amount: "-"
  });
  const [loading, setLoading] = useState(true);
  useEffect(() => {
    async function init() {
      const res = await GetCurrentPeriod();
      if (res.code === 0) {
        const data = res.data;
        setCurrentPeriod({
          //Feb 25, 2025 - Mar 27, 2025
          billingCycle: `${data.period_start} - ${data.period_end}`,
          amount: `$ ${data.amount}`
        });
      }
    }

    setLoading(true);
    void init().finally(() => setLoading(false));
  }, []);
  return <Page title="Billings">
    <BlockStack gap="400">
      {loading ? <BillingCardSkeleton /> : <Card>
        <InlineStack align="space-between" blockAlign="start">
          <Text as="h2" variant="headingSm">Current protection billing </Text>
          <Button url="/billing/detail" disabled={currentPeriod.billingCycle === " - "}>View details</Button>
        </InlineStack>
        <BlockStack gap="100">
          <Text as="p" tone="subdued">
            Billing cycle: {currentPeriod.billingCycle}
          </Text>
          <Text as="p" variant="headingXl">
            {currentPeriod.amount}
          </Text>
        </BlockStack>
      </Card>}
      <BillingTable />
    </BlockStack>
  </Page>;
}