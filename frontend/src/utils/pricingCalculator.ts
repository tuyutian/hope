import { PricingSettings } from "@/types/cart";

/**
 * 计算保险价格
 * @param productPrice 商品价格
 * @param pricingSettings 定价设置
 * @returns 计算出的保险价格
 */
export function calculateProtectionPrice(productPrice: number, pricingSettings: PricingSettings): number {
  const {
    pricingType,
    pricingRule,
    priceSelect,
    tiersSelect,
    allPriceValue,
    allTiersValue,
    outPriceValue,
    outTierValue,
  } = pricingSettings;

  // pricingType: "0" = 固定价格, "1" = 百分比
  // pricingRule: "0" = 单一价格/百分比, "1" = 价格范围

  if (pricingType === "0") {
    // 固定价格模式
    if (pricingRule === "0") {
      // 单一固定价格
      return parseFloat(allPriceValue) || 0;
    } else {
      // 价格范围模式
      for (const range of priceSelect) {
        const min = parseFloat(range.min);
        const max = parseFloat(range.max);
        const price = parseFloat(range.price);

        if (productPrice >= min && productPrice <= max) {
          return price || 0;
        }
      }
      // 如果不在任何范围内，返回超出范围价格
      return parseFloat(outPriceValue) || 0;
    }
  } else {
    // 百分比模式
    if (pricingRule === "0") {
      // 单一百分比
      const percentage = parseFloat(allTiersValue) || 0;
      return (productPrice * percentage) / 100;
    } else {
      // 百分比范围模式
      for (const range of tiersSelect) {
        const min = parseFloat(range.min);
        const max = parseFloat(range.max);
        const percentage = parseFloat(range.percentage);

        if (productPrice >= min && productPrice <= max) {
          return (productPrice * (percentage || 0)) / 100;
        }
      }
      // 如果不在任何范围内，返回超出范围百分比计算
      const outPercentage = parseFloat(outTierValue) || 0;
      return (productPrice * outPercentage) / 100;
    }
  }
}

/**
 * 格式化价格显示
 * @param price 价格数值
 * @param currency 货币符号
 * @returns 格式化后的价格字符串
 */
export function formatPrice(price: number, currency: string = "$"): string {
  return `${price.toFixed(2)} ${currency === "$" ? "USD" : currency}`;
}
