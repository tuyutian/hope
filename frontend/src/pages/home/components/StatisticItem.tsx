import React from "react";
import { BlockStack, Box, Text } from "@shopify/polaris";
import HomeTooltip from "./HomeTooltip";

interface StatisticItemProps {
  title: string;
  value: number | string;
  prefix?: string;
  color?: string;
  tooltipContent?: React.ReactNode;
}

const StatisticItem: React.FC<StatisticItemProps> = ({ title, value, prefix = "", tooltipContent }) => {
  return (
    <div
      style={{ width: "100%", padding: "12px" }}
      className=" Polaris-Button Polaris-Button--pressable Polaris-Button--variantTertiary Polaris-Button--sizeMedium Polaris-Button--textAlignCenter "
    >
      <BlockStack gap="100">
        <HomeTooltip
          width={244}
          title={title}
          text={
            tooltipContent || (
              <Box padding="200">
                <Text as="p" variant="bodyMd" fontWeight="medium">
                  {title}
                </Text>
                <Box />
              </Box>
            )
          }
        />
        <Text as="span" variant="headingLg">
          {prefix}
          {value}
        </Text>
      </BlockStack>
    </div>
  );
};

export default StatisticItem;
