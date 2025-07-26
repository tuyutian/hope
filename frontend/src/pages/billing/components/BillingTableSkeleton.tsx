import {Box, Card, SkeletonBodyText, Text} from "@shopify/polaris";

export default function BillingTableSkeleton() {
  return <Card padding="0">
    <Box padding="400">
      <Text as="h2" variant="headingSm">
        Past protection billing
      </Text>
    </Box>
    <Box padding="400">
      <SkeletonBodyText lines={10} />
    </Box>
  </Card>;
}