import React from "react";
import { Button, DataTable, TextField } from "@shopify/polaris";

interface PriceTier {
  min: string;
  max: string;
  price: string;
}

interface PercentageTier {
  min: string;
  max: string;
  percentage: string;
}

interface PricingTableProps {
  pricingType: string;
  priceSelect: PriceTier[];
  tiersSelect: PercentageTier[];
  moneySymbol: string;
  errors: Record<string, string>;
  onPriceChange: (index: number, field: string, value: string, error: string) => void;
  onTierChange: (index: number, field: string, value: string, error: string) => void;
  onDeletePrice: (index: number) => void;
  onDeleteTier: (index: number) => void;
  onAddPrice: () => void;
  onAddTier: () => void;
}

const PricingTable: React.FC<PricingTableProps> = ({
  pricingType,
  priceSelect,
  tiersSelect,
  moneySymbol,
  errors,
  onPriceChange,
  onTierChange,
  onDeletePrice,
  onDeleteTier,
  onAddPrice,
  onAddTier,
}) => {
  if (pricingType === "1") {
    // Percentage pricing
    return (
      <>
        <DataTable
          columnContentTypes={["text", "text", "text", "text"]}
          headings={["Min Cart Value", "Max Cart Value", "Protection Price", ""]}
          rows={tiersSelect.map((tier, index) => [
            <TextField
              label=""
              autoComplete="off"
              key={`min-${index}`}
              value={tier.min}
              onChange={(value) => onTierChange(index, "min", value, `tier_min_${index}`)}
              prefix={moneySymbol}
              error={errors[`tier_min_${index}`]}
            />,
            <TextField
              label=""
              autoComplete="off"
              key={`max-${index}`}
              value={tier.max}
              onChange={(value) => onTierChange(index, "max", value, `tier_max_${index}`)}
              prefix={moneySymbol}
              error={errors[`tier_max_${index}`]}
            />,
            <TextField
              label=""
              autoComplete="off"
              key={`percentage-${index}`}
              value={tier.percentage}
              onChange={(value) => onTierChange(index, "percentage", value, `tier_percentage_${index}`)}
              prefix="%"
              error={errors[`tier_percentage_${index}`]}
            />,
            <Button
              key={`delete-${index}`}
              onClick={() => onDeleteTier(index)}
              variant="plain"
              tone="critical"
            >
              Delete
            </Button>,
          ])}
        />
        <Button onClick={onAddTier}>Add Tier</Button>
      </>
    );
  } else {
    // Fixed pricing
    return (
      <>
        <DataTable
          columnContentTypes={["text", "text", "text", "text"]}
          headings={["Min Cart Value", "Max Cart Value", "Protection Price", ""]}
          rows={priceSelect.map((tier, index) => [
            <TextField
              label=""
              autoComplete="off"
              key={`min-${index}`}
              value={tier.min}
              onChange={(value) => onPriceChange(index, "min", value, `price_min_${index}`)}
              prefix={moneySymbol}
              error={errors[`price_min_${index}`]}
            />,
            <TextField
              label=""
              autoComplete="off"
              key={`max-${index}`}
              value={tier.max}
              onChange={(value) => onPriceChange(index, "max", value, `price_max_${index}`)}
              prefix={moneySymbol}
              error={errors[`price_max_${index}`]}
            />,
            <TextField
              label=""
              autoComplete="off"
              key={`price-${index}`}
              value={tier.price}
              onChange={(value) => onPriceChange(index, "price", value, `price_price_${index}`)}
              prefix={moneySymbol}
              error={errors[`price_price_${index}`]}
            />,
            <Button
              key={`delete-${index}`}
              onClick={() => onDeletePrice(index)}
              variant="plain"
              tone="critical"
            >
              Delete
            </Button>,
          ])}
        />
        <Button onClick={onAddPrice}>Add Price</Button>
      </>
    );
  }
};

export default PricingTable;
