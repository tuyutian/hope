import { BlockStack, Button, Card, InlineStack, Page, Text } from "@shopify/polaris";
import BillingTable from "@/pages/billing/components/BillingTable.tsx";
import { useEffect, useState } from "react";
import { billingService } from "@/api";
import BillingCardSkeleton from "@/pages/billing/components/BillingCardSkeleton.tsx";
import { formatUnixTimestampRange } from "@/utils/dateUtils";
import { useNavigate } from "react-router";

export default function Billing() {
  const [currentPeriod, setCurrentPeriod] = useState<{ billingCycle: string; amount: string }>({
    billingCycle: "-",
    amount: "-",
  });
  const [currentPeriodStart, setCurrentPeriodStart] = useState(0);
  const [currentPeriodEnd, setCurrentPeriodEnd] = useState(0);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();
  useEffect(() => {
    async function init() {
      const res = await billingService.getCurrentPeriod();
      if (res.code === 0 && res.data) {
        const data = res.data;
        if (data.period_start === 0 && data.period_end === 0) {
          setCurrentPeriod({
            billingCycle: " - ",
            amount: " - ",
          });
          return;
        }
        setCurrentPeriodStart(data.period_start);
        setCurrentPeriodEnd(data.period_end);
        setCurrentPeriod({
          billingCycle: formatUnixTimestampRange(data.period_start, data.period_end),
          amount: `$ ${data.amount}`,
        });
      }
    }

    setLoading(true);
    void init().finally(() => setLoading(false));
  }, []);
  const handleViewDetails = async (minTime: number, maxTime: number) => {
    if (minTime === 0 && maxTime === 0) {
      return;
    }
    await navigate(`/billing/detail?minTime=${minTime}&maxTime=${maxTime}`);
  };
  return (
    <Page title="Billings">
      <BlockStack gap="400">
        {loading ? (
          <BillingCardSkeleton />
        ) : (
          <Card>
            <InlineStack align="space-between" blockAlign="start">
              <Text as="h2" variant="headingSm">
                Current protection billing{" "}
              </Text>
              <Button
                onClick={() => handleViewDetails(currentPeriodStart, currentPeriodEnd)}
                disabled={currentPeriod.billingCycle === " - "}
              >
                View details
              </Button>
            </InlineStack>
            <BlockStack gap="100">
              <Text as="p" tone="subdued">
                Billing cycle: {currentPeriod.billingCycle}
              </Text>
              <Text as="p" variant="headingXl">
                {currentPeriod.amount}
              </Text>
            </BlockStack>
          </Card>
        )}
        <BillingTable />
      </BlockStack>
    </Page>
  );
}
