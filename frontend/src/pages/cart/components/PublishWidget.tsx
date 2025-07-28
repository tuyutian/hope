import React from "react";
import { Box, Card, InlineStack, Link, Text } from "@shopify/polaris";
import CustomSwitch from "./CustomSwitch";

interface PublishWidgetProps {
  protectifyVisibility: string;
  loading: boolean;
  onProtectifyVisibilityChange: (value: boolean) => void;
}

const PublishWidget: React.FC<PublishWidgetProps> = ({
  protectifyVisibility,
  loading,
  onProtectifyVisibilityChange,
}) => {
  return (
    <Card padding="300">
      <InlineStack align="space-between" blockAlign="center" gap="200">
        <Text variant="headingSm" as="h6">
          Publish Widget
        </Text>
        <Box>
          <CustomSwitch
            loading={loading}
            onChange={onProtectifyVisibilityChange}
            checked={protectifyVisibility === "1"}
          />
        </Box>
      </InlineStack>
      <Text as="p">
        Please follow the 👉 <Link url="#">help docs</Link> to complete setup. If after publishing the widget, you find
        that the widget does not show up or work properly in store cart, please turn off this switch only. This way the
        widget will not have any effect in the cart, and then please contact us for a free expert adaptation.
      </Text>
    </Card>
  );
};

export default PublishWidget;
