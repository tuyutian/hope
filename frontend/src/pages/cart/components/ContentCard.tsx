import {ContentCardProps} from "@/types/cart.ts";
import {BlockStack, Card,Text, TextField} from "@shopify/polaris";

export default function ContentCard({ widgetSettings, errors, onFieldChange }: ContentCardProps) {
  return (
    <Card padding="400">
      <BlockStack gap="200">
        <Text variant="headingSm" as="h6">Content</Text>

        <TextField
          autoComplete="off"
          label="Add-on title"
          value={widgetSettings.addonTitle}
          onChange={(value) => onFieldChange({ addonTitle: value })}
          error={errors.addonTitle}
          maxLength={50}
        />

        <TextField
          autoComplete="off"
          label="Enabled description"
          value={widgetSettings.enabledDescription}
          onChange={(value) => onFieldChange({ enabledDescription: value })}
          multiline={4}
          maxLength={200}
          error={errors.enabledDescription}
        />

        <TextField
          autoComplete="off"
          label="Disabled description"
          value={widgetSettings.disabledDescription}
          onChange={(value) => onFieldChange({ disabledDescription: value })}
          multiline={4}
          maxLength={200}
          error={errors.disabledDescription}
        />

        <TextField
          autoComplete="off"
          label="Footer link text"
          value={widgetSettings.footerText}
          onChange={(value) => onFieldChange({ footerText: value })}
          maxLength={50}
        />

        <TextField
          autoComplete="off"
          label="Footer link URL"
          value={widgetSettings.footerUrl}
          onChange={(value) => onFieldChange({ footerUrl: value })}
          maxLength={150}
        />

        <Text tone="subdued" variant="bodySm" as="span">
          Note: leave blank for no link
        </Text>
      </BlockStack>
    </Card>
  );
}