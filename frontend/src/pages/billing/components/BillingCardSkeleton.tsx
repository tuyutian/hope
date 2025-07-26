import {Card, InlineStack, SkeletonBodyText, SkeletonDisplayText} from "@shopify/polaris";

export default function BillingCardSkeleton() {
  return  <Card>
    <InlineStack align="space-between" blockAlign="start">
      <SkeletonDisplayText />
      <SkeletonDisplayText />
    </InlineStack>
    <SkeletonBodyText lines={4} />
  </Card>;
}