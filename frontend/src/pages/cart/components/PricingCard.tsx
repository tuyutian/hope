import type { PricingCardProps } from "@/types/cart.ts";
import { BlockStack, Card, Divider, Select, Text, TextField } from "@shopify/polaris";
import PricingTable from "@/pages/cart/components/PricingTable.tsx";
import React from "react";

// 输入验证函数
const validatePriceInput = (value: string) => {
  // 只允许数字和一个小数点，最多保留2位小数，最长10个字符
  const regex = /^\d{0,7}(\.\d{0,2})?$/;
  return regex.test(value) && value.length <= 10;
};

export default function PricingCard({
  pricingSettings,
  moneySymbol,
  errors,
  onSettingsChange,
  onPricingChange,
  onAddItem,
  onDeleteItem,
}: PricingCardProps) {
  // 处理单一价格/百分比输入
  const handleSingleValueChange = (value: string) => {
    if (!validatePriceInput(value)) return;

    if (pricingSettings.pricingType === "1") {
      onSettingsChange({ allTiersValue: value });
    } else {
      onSettingsChange({ allPriceValue: value });
    }
  };

  // 处理超出范围价格/百分比输入
  const handleOutOfRangeValueChange = (value: string) => {
    if (!validatePriceInput(value)) return;

    if (pricingSettings.pricingType === "1") {
      onSettingsChange({ outTierValue: value });
    } else {
      onSettingsChange({ outPriceValue: value });
    }
  };

  return (
    <Card padding="400">
      <BlockStack gap="400">
        <Text variant="headingSm" as="h6">
          Protection Pricing
        </Text>

        <Select
          label="Pricing type"
          options={[
            { label: "Fixed", value: "0" },
            { label: "Percentage", value: "1" },
          ]}
          value={pricingSettings.pricingType}
          onChange={value => onSettingsChange({ pricingType: value })}
          helpText={
            pricingSettings.pricingType === "0"
              ? "Set insurance fees as fixed amounts"
              : "Set insurance fees as a percentage of product value"
          }
        />

        <Select
          label="Pricing rule"
          options={[
            { label: "Apply for all cart value", value: "0" },
            { label: "Apply for different cart value ranges", value: "1" },
          ]}
          value={pricingSettings.pricingRule}
          onChange={value => onSettingsChange({ pricingRule: value })}
          helpText={
            pricingSettings.pricingRule === "0"
              ? "Single pricing rule applied to all cart values"
              : "Different pricing rules for different cart value ranges"
          }
        />

        {/* 单一价格/百分比设置 */}
        {pricingSettings.pricingRule === "0" && (
          <TextField
            label={pricingSettings.pricingType === "1" ? "Protection percentage" : "Protection price"}
            autoComplete="off"
            value={pricingSettings.pricingType === "1" ? pricingSettings.allTiersValue : pricingSettings.allPriceValue}
            onChange={handleSingleValueChange}
            prefix={pricingSettings.pricingType === "1" ? "%" : moneySymbol}
            helpText={
              pricingSettings.pricingType === "1" ? "Enter percentage (e.g., 10 for 10%)" : "Enter fixed amount"
            }
            error={errors.allPriceValue || errors.allTiersValue}
          />
        )}

        {/* 价格范围表格 */}
        {pricingSettings.pricingRule === "1" && (
          <>
            <PricingTable
              pricingType={pricingSettings.pricingType}
              priceSelect={pricingSettings.priceSelect}
              tiersSelect={pricingSettings.tiersSelect}
              moneySymbol={moneySymbol}
              errors={errors}
              onPriceChange={onPricingChange}
              onTierChange={onPricingChange}
              onDeletePrice={onDeleteItem}
              onDeleteTier={onDeleteItem}
              onAddPrice={onAddItem}
              onAddTier={onAddItem}
            />

            {/* 超出范围设置 */}
            <TextField
              label="The Rest Value Range"
              autoComplete="off"
              value={pricingSettings.pricingType === "1" ? pricingSettings.outTierValue : pricingSettings.outPriceValue}
              onChange={handleOutOfRangeValueChange}
              prefix={pricingSettings.pricingType === "1" ? "%" : moneySymbol}
              helpText="Pricing rule for cart values not covered by the ranges above"
              error={errors.outPriceValue || errors.outTierValue}
            />
          </>
        )}

        <Divider />
      </BlockStack>
    </Card>
  );
}
