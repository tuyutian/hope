import {BlockStack, Card, SkeletonBodyText, SkeletonPage, SkeletonTabs} from "@shopify/polaris";

export default function BillingTableSkeleton() {
  return <SkeletonPage>
    <Card padding="0">
      <BlockStack gap="300">
        <SkeletonTabs />
        <SkeletonBodyText lines={10} />
      </BlockStack>
    </Card>
  </SkeletonPage>;
}