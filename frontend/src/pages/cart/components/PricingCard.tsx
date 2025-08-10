import type { PricingCardProps } from "@/types/cart.ts";
import { BlockStack, Card, Divider, Select, Text, TextField } from "@shopify/polaris";
import PricingTable from "@/pages/cart/components/PricingTable.tsx";
import React from "react";

export default function PricingCard({
  pricingSettings,
  moneySymbol,
  errors,
  onSettingsChange,
  onPricingChange,
  onAddItem,
  onDeleteItem,
}: PricingCardProps) {
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
        />

        <Select
          label="Pricing rule"
          options={[
            { label: "App for all cart vale", value: "0" },
            { label: "Apply for different cart value range", value: "1" },
          ]}
          value={pricingSettings.pricingRule}
          onChange={value => onSettingsChange({ pricingRule: value })}
        />

        {pricingSettings.pricingRule === "1" && (
          <PricingTable
            pricingType={pricingSettings.pricingType}
            priceSelect={pricingSettings.priceSelect}
            tiersSelect={pricingSettings.tiersSelect}
            moneySymbol={moneySymbol}
            errors={errors}
            onPriceChange={(index, field, value) => onPricingChange(index, field, value, "price")}
            onTierChange={(index, field, value) => onPricingChange(index, field, value, "tier")}
            onDeletePrice={index => onDeleteItem(index, "price")}
            onDeleteTier={index => onDeleteItem(index, "tier")}
            onAddPrice={() => onAddItem("price")}
            onAddTier={() => onAddItem("tier")}
          />
        )}

        {pricingSettings.pricingRule === "1" ? (
          <TextField
            autoComplete="off"
            label="Other Value Range"
            value={pricingSettings.restValuePrice}
            onChange={value => onSettingsChange({ restValuePrice: value })}
            prefix={moneySymbol}
          />
        ) : (
          <TextField
            label=""
            autoComplete="off"
            value={pricingSettings.pricingType === "1" ? pricingSettings.allTiersValue : pricingSettings.allPriceValue}
            onChange={value =>
              onSettingsChange(
                pricingSettings.pricingType === "1" ? { allTiersValue: value } : { allPriceValue: value }
              )
            }
            prefix={pricingSettings.pricingType === "1" ? "%" : moneySymbol}
          />
        )}

        <Divider />
      </BlockStack>
    </Card>
  );
}
