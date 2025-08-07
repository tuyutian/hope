import { WidgetStyleCardProps } from "@/types/cart.ts";
import { BlockStack, Box, Card, InlineStack, Select, Text } from "@shopify/polaris";
import IconSelector from "@/pages/cart/components/IconSelector.tsx";
import SketchPickerWithInput from "@/components/form/SketchPickerWithInput.tsx";
import CustomSwitch from "@/pages/cart/components/CustomSwitch.tsx";
import React from "react";

export default function WidgetStyleCard({
  widgetSettings,
  icons,
  onWidgetSettingsChange,
  onIconClick,
  onIconUpload,
}: WidgetStyleCardProps) {
  return (
    <Card padding="400">
      <BlockStack gap="300">
        <Text variant="headingSm" as="h6">
          Widget Style
        </Text>
        <Box>
          <InlineStack gap="200" align="space-between">
            <Text as="p" variant="bodyMd" fontWeight="medium">
              Icon visibility
            </Text>
            <CustomSwitch
              onChange={value => onWidgetSettingsChange({ iconVisibility: value ? "1" : "0" })}
              checked={widgetSettings.iconVisibility === "1"}
            />
          </InlineStack>
        </Box>

        <IconSelector icons={icons} onIconClick={onIconClick} onIconUpload={onIconUpload} />

        <Select
          label="Select Button"
          options={[
            { label: "Switch", value: "0" },
            { label: "Checkbox", value: "1" },
          ]}
          value={widgetSettings.selectButton}
          onChange={value => onWidgetSettingsChange({ selectButton: value })}
        />

        <InlineStack align="space-between" gap="300" wrap={false}>
          <Box>
            <div>Opt-in action button</div>
            <SketchPickerWithInput
              defaultColor={widgetSettings.optInColor}
              onChange={color => onWidgetSettingsChange({ optInColor: color })}
            />
          </Box>
          <Box>
            <div>Opt-out action button</div>
            <SketchPickerWithInput
              defaultColor={widgetSettings.optOutColor}
              onChange={(color: string) => onWidgetSettingsChange({ optOutColor: color })}
            />
          </Box>
        </InlineStack>
      </BlockStack>
    </Card>
  );
}
