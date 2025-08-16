import React from "react";
import { BlockStack, Card, Text, TextField, Box, Link } from "@shopify/polaris";
import { WidgetSettings } from "@/types/cart";
import { handleContact } from "@/utils/app.ts";

interface CSSCardProps {
  widgetSettings: WidgetSettings;
  onFieldChange: (value: Partial<WidgetSettings>) => void;
}

export default function CSSCard({ widgetSettings, onFieldChange }: CSSCardProps) {
  return (
    <Card>
      <Box paddingBlockEnd="300">
        <Text as="h2" variant="headingSm">
          CSS
        </Text>
      </Box>
      <BlockStack gap="300">
        <Text as="p" variant="bodyMd" tone="subdued">
          If you would like to adjust the styling of the widgets in your store,{" "}
          <Link onClick={handleContact}>contact us</Link> and we will add CSS code here to make custom changes. This
          won&#39;t affect your store theme.{" "}
        </Text>

        <TextField
          autoComplete="off"
          label=""
          value={widgetSettings.css}
          onChange={value => onFieldChange({ css: value })}
          multiline={10}
          placeholder="/* Add your custom CSS here */
.protectify-widget {
  /* Example styles */
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}"
          helpText=""
        />
      </BlockStack>
    </Card>
  );
}
