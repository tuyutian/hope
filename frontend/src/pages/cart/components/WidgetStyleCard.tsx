import { WidgetStyleCardProps } from "@/types/cart.ts";
import { BlockStack, Box, Card, InlineStack, Select, Text } from "@shopify/polaris";
import IconSelector from "@/pages/cart/components/IconSelector.tsx";
import SketchPickerWithInput from "@/components/form/SketchPickerWithInput.tsx";

export default function WidgetStyleCard({
  widgetSettings,
  icons,
  onWidgetSettingsChange,
  onIconClick,
}: WidgetStyleCardProps) {
  return (
    <Card padding="400">
      <BlockStack gap="300">
        <Text variant="headingSm" as="h6">
          Widget Style
        </Text>

        <Select
          label="Icon Visibility"
          options={[
            { label: "Show Icon", value: "1" },
            { label: "Hide Icon", value: "0" },
          ]}
          value={widgetSettings.iconVisibility}
          onChange={value => onWidgetSettingsChange({ iconVisibility: value })}
        />

        <IconSelector icons={icons} onIconClick={onIconClick} />

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
