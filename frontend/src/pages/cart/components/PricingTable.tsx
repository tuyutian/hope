import React from "react";
import { Button, DataTable, TextField, InlineStack } from "@shopify/polaris";
import { PlusIcon, DeleteIcon } from "@shopify/polaris-icons";
import { getMessageState } from "@/stores/messageStore";

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
  onPriceChange: (index: number, field: string, value: string, type: "price" | "tier") => void;
  onTierChange: (index: number, field: string, value: string, type: "price" | "tier") => void;
  onDeletePrice: (index: number, type: "price" | "tier") => void;
  onDeleteTier: (index: number, type: "price" | "tier") => void;
  onAddPrice: (type: "price" | "tier") => void;
  onAddTier: (type: "price" | "tier") => void;
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
  const toastMessage = getMessageState().toastMessage;

  // 输入验证函数
  const validatePriceInput = (value: string) => {
    const regex = /^\d{0,7}(\.\d{0,2})?$/;
    return regex.test(value) && value.length <= 10;
  };

  // 检查价格范围重叠
  const checkRangeOverlap = (ranges: any[], newIndex: number, field: string, value: string) => {
    if (!value) return false;

    const updatedRanges = ranges.map((range, index) => (index === newIndex ? { ...range, [field]: value } : range));

    for (let i = 0; i < updatedRanges.length; i++) {
      for (let j = i + 1; j < updatedRanges.length; j++) {
        const range1 = updatedRanges[i];
        const range2 = updatedRanges[j];

        const min1 = parseFloat(range1.min) || 0;
        const max1 = parseFloat(range1.max) || 0;
        const min2 = parseFloat(range2.min) || 0;
        const max2 = parseFloat(range2.max) || 0;

        if (min1 <= max2 && min2 <= max1) {
          return true;
        }
      }
    }
    return false;
  };

  // 处理价格范围输入
  const handlePriceRangeChange = (index: number, field: string, value: string, type: "price" | "tier") => {
    if (!validatePriceInput(value)) return;

    const ranges = type === "price" ? priceSelect : tiersSelect;

    // 检查重叠
    if ((field === "min" || field === "max") && checkRangeOverlap(ranges, index, field, value)) {
      toastMessage("Price ranges cannot overlap", 3000, true);
      return;
    }

    if (type === "price") {
      onPriceChange(index, field, value, type);
    } else {
      onTierChange(index, field, value, type);
    }
  };

  // 处理添加新区间
  const handleAddTier = (type: "price" | "tier") => {
    const currentRanges = type === "price" ? priceSelect : tiersSelect;

    if (currentRanges.length >= 5) {
      toastMessage("Up to 5 supported", 3000, true);
      return;
    }

    if (type === "price") {
      onAddPrice(type);
    } else {
      onAddTier(type);
    }
  };

  // 处理删除区间
  const handleDelete = (index: number, type: "price" | "tier") => {
    if (type === "price") {
      onDeletePrice(index, type);
    } else {
      onDeleteTier(index, type);
    }
    toastMessage("Deleted successfully", 2000, false);
  };

  if (pricingType === "1") {
    // 百分比定价
    const rows = tiersSelect.map((tier, index) => [
      <TextField
        label=""
        autoComplete="off"
        key={`min-${index}`}
        value={tier.min}
        onChange={value => handlePriceRangeChange(index, "min", value, "tier")}
        prefix={moneySymbol}
        error={errors[`tier_min_${index}`]}
        placeholder="0.00"
      />,
      <TextField
        label=""
        autoComplete="off"
        key={`max-${index}`}
        value={tier.max}
        onChange={value => handlePriceRangeChange(index, "max", value, "tier")}
        prefix={moneySymbol}
        error={errors[`tier_max_${index}`]}
        placeholder="100.00"
      />,
      <TextField
        label=""
        autoComplete="off"
        key={`percentage-${index}`}
        value={tier.percentage}
        onChange={value => handlePriceRangeChange(index, "percentage", value, "tier")}
        suffix="%"
        error={errors[`tier_percentage_${index}`]}
        placeholder="10"
      />,
      <Button
        key={`delete-${index}`}
        onClick={() => handleDelete(index, "tier")}
        variant="plain"
        tone="critical"
        icon={DeleteIcon}
        accessibilityLabel="Delete tier"
      />,
    ]);

    return (
      <>
        <DataTable
          columnContentTypes={["text", "text", "text", "text"]}
          headings={["Min Cart Value", "Max Cart Value", "Protection Percentage", ""]}
          rows={rows}
        />
        <InlineStack>
          <Button icon={PlusIcon} onClick={() => handleAddTier("tier")}>
            Add Tier
          </Button>
        </InlineStack>
      </>
    );
  } else {
    // 固定价格定价
    const rows = priceSelect.map((tier, index) => [
      <TextField
        label=""
        autoComplete="off"
        key={`min-${index}`}
        value={tier.min}
        onChange={value => handlePriceRangeChange(index, "min", value, "price")}
        prefix={moneySymbol}
        error={errors[`price_min_${index}`]}
        placeholder="0.00"
      />,
      <TextField
        label=""
        autoComplete="off"
        key={`max-${index}`}
        value={tier.max}
        onChange={value => handlePriceRangeChange(index, "max", value, "price")} // 修复了这里的拼写错误
        prefix={moneySymbol}
        error={errors[`price_max_${index}`]}
        placeholder="100.00"
      />,
      <TextField
        label=""
        autoComplete="off"
        key={`price-${index}`}
        value={tier.price}
        onChange={value => handlePriceRangeChange(index, "price", value, "price")}
        prefix={moneySymbol}
        error={errors[`price_price_${index}`]}
        placeholder="5.00"
      />,
      <Button
        key={`delete-${index}`}
        onClick={() => handleDelete(index, "price")}
        variant="plain"
        tone="critical"
        icon={DeleteIcon}
        accessibilityLabel="Delete price range"
      />,
    ]);

    return (
      <>
        <DataTable
          columnContentTypes={["text", "text", "text", "text"]}
          headings={["Min Cart Value", "Max Cart Value", "Protection Price", ""]}
          rows={rows}
        />
        <InlineStack>
          <Button icon={PlusIcon} onClick={() => handleAddTier("price")}>
            Add Tier
          </Button>
        </InlineStack>
      </>
    );
  }
};

export default PricingTable;
