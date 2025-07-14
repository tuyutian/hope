import { useState, useEffect, useCallback } from "react";
import { GetUserConf, rqGetCartSetting, rqPostUpdateCartSetting } from "@/api";
import type { WidgetSettings, PricingSettings, ProductSettings, CartSettingsHook } from "@/types/cart";
import { getMessageState } from "@/stores/messageStore";
import type { OptionDescriptor } from "@shopify/polaris/build/ts/src/types";

export function useCartSettings(): CartSettingsHook {
  const toastMessage = getMessageState().toastMessage;

  const [widgetSettings, setWidgetSettings] = useState<WidgetSettings>({
    planTitle: "Plan Title",
    iconVisibility: "0",
    insuranceVisibility: "0",
    selectButton: "0",
    addonTitle: "Shipping Protection",
    enabledDescription:
      "After purchasing this insurance, we will resolve all after-sales issues related to this order for you.",
    disabledDescription:
      "After purchasing this insurance, we will resolve all after-sales issues related to this order for you.",
    footerText: "",
    footerUrl: "",
    optInColor: "#fffff",
    optOutColor: "#fffff",
  });

  const [pricingSettings, setPricingSettings] = useState<PricingSettings>({
    pricingType: "0",
    pricingRule: "0",
    priceSelect: [
      { min: "1.00", max: "100.00", price: "3.00" },
      { min: "101.00", max: "1000.00", price: "10.00" },
    ],
    tiersSelect: [
      { min: "1.00", max: "100.00", percentage: "10" },
      { min: "101.00", max: "1000.00", percentage: "20" },
    ],
    restValuePrice: "10.00",
    allPriceValue: "",
    allTiersValue: "",
  });

  const [productSettings, setProductSettings] = useState<ProductSettings>({
    productTypeInput: "",
    selectedCollections: [],
    collectionInput: "",
    icons: [
      { id: 1, src: "https://img.icons8.com/color/48/shield.png", selected: true },
      { id: 2, src: "https://maxst.icons8.com/vue-static/faceswapper/hero/faces/2.jpg", selected: false },
    ],
  });

  const [collectionOptions, setCollectionOptions] = useState<OptionDescriptor[]>([]);
  const [moneySymbol, setMoneySymbol] = useState("$");
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [isLoading, setIsLoading] = useState(true);
  const [saveLoading, setSaveLoading] = useState(false);
  const [dirty, setDirty] = useState(false);

  const loadInitialData = useCallback(async () => {
    try {
      setIsLoading(true);
      await Promise.all([loadUserConfig(), loadCartData()]);
    } finally {
      setIsLoading(false);
    }
  }, []);

  const loadUserConfig = async () => {
    const res = await GetUserConf();
    if (res.code !== 0 || !res.data) return;

    const { collections, money_symbol } = res.data;
    if (Array.isArray(collections)) {
      setCollectionOptions(
        collections.map((collection: { label: string; value: number }) => ({
          label: collection.label,
          value: String(collection.value),
        }))
      );
    }
    if (money_symbol) setMoneySymbol(money_symbol);
  };

  const loadCartData = async () => {
    const res = await rqGetCartSetting();
    if (res.code !== 0 || !res.data) return;

    const data = res.data;

    // 更新widget设置
    setWidgetSettings(prev => ({
      ...prev,
      planTitle: data.plan_title || prev.planTitle,
      addonTitle: data.addon_title || prev.addonTitle,
      enabledDescription: data.enabled_desc || prev.enabledDescription,
      disabledDescription: data.disabled_desc || prev.disabledDescription,
      footerText: data.foot_text || prev.footerText,
      footerUrl: data.foot_url || prev.footerUrl,
      optInColor: data.in_color || prev.optInColor,
      optOutColor: data.out_color || prev.optOutColor,
      insuranceVisibility: typeof data.show_cart === "number" ? String(data.show_cart) : prev.insuranceVisibility,
      iconVisibility: typeof data.show_cart_icon === "number" ? String(data.show_cart_icon) : prev.iconVisibility,
      selectButton: typeof data.select_button === "number" ? String(data.select_button) : prev.selectButton,
    }));

    // 更新pricing设置
    setPricingSettings(prev => ({
      ...prev,
      pricingType: typeof data.pricing_type === "number" ? String(data.pricing_type) : prev.pricingType,
      pricingRule: typeof data.price_rule === "number" ? String(data.price_rule) : prev.pricingRule,
      priceSelect: Array.isArray(data.price_select) ? data.price_select : prev.priceSelect,
      tiersSelect: Array.isArray(data.tiers_select) ? data.tiers_select : prev.tiersSelect,
      restValuePrice: typeof data.other_money === "number" ? String(data.other_money) : prev.restValuePrice,
      allPriceValue: typeof data.all_price === "number" ? String(data.all_price) : prev.allPriceValue,
      allTiersValue: typeof data.all_tiers === "number" ? String(data.all_tiers) : prev.allTiersValue,
    }));

    // 更新产品设置
    setProductSettings(prev => ({
      ...prev,
      productTypeInput: data.product_type || prev.productTypeInput,
      selectedCollections: Array.isArray(data.product_collection) ? data.product_collection : prev.selectedCollections,
      icons: Array.isArray(data.icons) && data.icons.length > 0 ? data.icons : prev.icons,
    }));
  };

  const validateFields = useCallback(() => {
    const newErrors: Record<string, string> = {};

    if (widgetSettings.insuranceVisibility === "1") {
      if (!widgetSettings.addonTitle.trim()) {
        newErrors.addonTitle = "Add-on Title is required";
      }
      if (!widgetSettings.enabledDescription.trim()) {
        newErrors.enabledDescription = "Enabled Description is required";
      }
      if (!widgetSettings.disabledDescription.trim()) {
        newErrors.disabledDescription = "Disabled Description is required";
      }

      // 验证定价规则
      if (pricingSettings.pricingType === "1" && pricingSettings.pricingRule === "1") {
        pricingSettings.tiersSelect.forEach((tier, index) => {
          if (tier.min || tier.max || tier.percentage) {
            if (!tier.min) newErrors[`tier_min_${index}`] = "Please fill in";
            if (!tier.max) newErrors[`tier_max_${index}`] = "Please fill in";
            if (!tier.percentage) newErrors[`tier_percentage_${index}`] = "Please fill in";
          }
        });
      }

      if (pricingSettings.pricingType === "0" && pricingSettings.pricingRule === "1") {
        pricingSettings.priceSelect.forEach((price, index) => {
          if (price.min || price.max || price.price) {
            if (!price.min) newErrors[`price_min_${index}`] = "Please fill in";
            if (!price.max) newErrors[`price_max_${index}`] = "Please fill in";
            if (!price.price) newErrors[`price_price_${index}`] = "Please fill in";
          }
        });
      }
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  }, [widgetSettings, pricingSettings]);

  const saveSettings = useCallback(async () => {
    if (!validateFields()) {
      toastMessage("Please fix validation errors", 5000, true);
      return;
    }

    setSaveLoading(true);
    try {
      const payload = {
        planTitle: widgetSettings.planTitle,
        iconVisibility: Number(widgetSettings.iconVisibility),
        insuranceVisibility: Number(widgetSettings.insuranceVisibility),
        selectButton: Number(widgetSettings.selectButton),
        addonTitle: widgetSettings.addonTitle,
        enabledDescription: widgetSettings.enabledDescription,
        disabledDescription: widgetSettings.disabledDescription,
        footerText: widgetSettings.footerText,
        footerUrl: widgetSettings.footerUrl,
        optInColor: widgetSettings.optInColor,
        optOutColor: widgetSettings.optOutColor,
        pricingType: Number(pricingSettings.pricingType),
        pricingRule: Number(pricingSettings.pricingRule),
        priceSelect: pricingSettings.priceSelect,
        tiersSelect: pricingSettings.tiersSelect,
        restValuePrice: pricingSettings.restValuePrice,
        allPrice: pricingSettings.allPriceValue,
        allTiers: pricingSettings.allTiersValue,
        productTypeInput: productSettings.productTypeInput,
        selectedCollections: productSettings.selectedCollections,
        icons: productSettings.icons,
      };

      const res = await rqPostUpdateCartSetting(payload);
      toastMessage(res?.code === 0 ? "Saved successfully" : "Saved Fail", 5000, res?.code !== 0);

      if (res?.code === 0) {
        setDirty(false);
      }
    } catch (error) {
      console.error("Error saving settings:", error);
      toastMessage("Service Error", 5000, true);
    } finally {
      setSaveLoading(false);
    }
  }, [widgetSettings, pricingSettings, productSettings, toastMessage, validateFields]);

  const markDirty = useCallback(() => {
    if (!dirty) setDirty(true);
  }, [dirty]);

  const discardChanges = useCallback(() => {
    setDirty(false);
  }, []);

  useEffect(() => {
    void loadInitialData();
  }, []);

  return {
    // 状态
    widgetSettings,
    pricingSettings,
    productSettings,
    collectionOptions,
    moneySymbol,
    errors,
    isLoading,
    saveLoading,
    dirty,

    // 更新函数
    setWidgetSettings,
    setPricingSettings,
    setProductSettings,
    setErrors,

    // 操作函数
    saveSettings,
    markDirty,
    discardChanges,
    validateFields,
  };
}
