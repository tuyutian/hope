import { useState, useEffect, useCallback, useRef } from "react";
import { cartService, userService } from "@/api";
import type {
  WidgetSettings,
  PricingSettings,
  ProductSettings,
  CartSettingsHook,
  FulfillmentSettings,
} from "@/types/cart";
import { getMessageState } from "@/stores/messageStore";

export function useCartSettings(): CartSettingsHook {
  const toastMessage = getMessageState().toastMessage;

  // 添加 ref 来存储初始数据
  const initialDataRef = useRef<{
    widgetSettings: WidgetSettings;
    pricingSettings: PricingSettings;
    productSettings: ProductSettings;
    fulfillmentSettings: FulfillmentSettings;
  } | null>(null);

  const [widgetSettings, setWidgetSettings] = useState<WidgetSettings>({
    planTitle: "Plan Title",
    iconVisibility: "0",
    protectifyVisibility: "0",
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
    css: "",
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
    allPriceValue: "0",
    allTiersValue: "0",
    outPriceValue: "0",
    outTierValue: "0",
  });

  const [productSettings, setProductSettings] = useState<ProductSettings>({
    selectedCollections: [],
    icons: [{ id: 1, src: "https://s.protectifyapp.com/logo.png", selected: true }],
    onlyInCollection: false,
  });
  const [fulfillmentSettings, setFulfillmentSettings] = useState<FulfillmentSettings>({
    fulfillmentRule: "0",
    fulfillmentOptions: [
      {
        label: "Mark as fulfilled when first item(s) are fulfilled",
        value: "0",
      },
      {
        label: "Mark as fulfilled when other items fulfilled",
        value: "1",
      },
      {
        label: "Mark as fulfilled immediately after purchase",
        value: "2",
      },
    ],
  });
  const [moneySymbol, setMoneySymbol] = useState("$");
  const [hasSubscribe, setHasSubscribe] = useState(false);
  const [hasEmbedInstalled, setHasEmbedInstalled] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [isLoading, setIsLoading] = useState(true);
  const [dirty, setDirty] = useState(false);

  // 深拷贝函数
  const deepClone = <T>(obj: T): T => {
    if (obj === null || typeof obj !== "object") return obj;
    if (obj instanceof Date) return new Date(obj.getTime()) as T;
    if (obj instanceof Array) return obj.map(item => deepClone(item)) as T;
    if (typeof obj === "object") {
      const cloned = {} as T;
      for (const key in obj) {
        if (Object.prototype.hasOwnProperty.call(obj, key)) {
          cloned[key] = deepClone(obj[key]);
        }
      }
      return cloned;
    }
    return obj;
  };

  // 深度比较函数
  const deepEqual = (obj1: any, obj2: any): boolean => {
    if (obj1 === obj2) return true;
    if (obj1 == null || obj2 == null) return false;
    if (typeof obj1 !== typeof obj2) return false;

    if (typeof obj1 !== "object") return obj1 === obj2;

    if (Array.isArray(obj1) !== Array.isArray(obj2)) return false;

    if (Array.isArray(obj1)) {
      if (obj1.length !== obj2.length) return false;
      for (let i = 0; i < obj1.length; i++) {
        if (!deepEqual(obj1[i], obj2[i])) return false;
      }
      return true;
    }

    const keys1 = Object.keys(obj1);
    const keys2 = Object.keys(obj2);
    if (keys1.length !== keys2.length) return false;

    for (const key of keys1) {
      if (!keys2.includes(key)) return false;
      if (!deepEqual(obj1[key], obj2[key])) return false;
    }

    return true;
  };

  // 检查是否有变化
  const checkForChanges = useCallback(() => {
    if (!initialDataRef.current) return false;

    const hasWidgetChanges = !deepEqual(widgetSettings, initialDataRef.current.widgetSettings);
    const hasPricingChanges = !deepEqual(pricingSettings, initialDataRef.current.pricingSettings);
    const hasProductChanges = !deepEqual(productSettings, initialDataRef.current.productSettings);
    const hasFulfillmentChanges = !deepEqual(fulfillmentSettings, initialDataRef.current.fulfillmentSettings);

    return hasWidgetChanges || hasPricingChanges || hasProductChanges || hasFulfillmentChanges;
  }, [widgetSettings, pricingSettings, productSettings, fulfillmentSettings]);

  // 更新dirty状态
  const updateDirtyState = useCallback(() => {
    const hasChanges = checkForChanges();
    if (dirty !== hasChanges) {
      setDirty(hasChanges);
    }
  }, [dirty, checkForChanges]);

  // 保存初始数据的函数
  const saveInitialData = useCallback(() => {
    initialDataRef.current = {
      widgetSettings: deepClone(widgetSettings),
      pricingSettings: deepClone(pricingSettings),
      productSettings: deepClone(productSettings),
      fulfillmentSettings: deepClone(fulfillmentSettings),
    };
  }, [widgetSettings, pricingSettings, productSettings, fulfillmentSettings]);

  const loadInitialData = useCallback(async () => {
    try {
      setIsLoading(true);
      await Promise.all([loadUserConfig(), loadCartData()]);
    } finally {
      setIsLoading(false);
    }
  }, []);

  const loadUserConfig = async () => {
    const res = await userService.getConfig();
    if (res.code !== 0 || !res.data) return;

    const { money_symbol, has_subscribe, has_embed_installed } = res.data;
    if (money_symbol) setMoneySymbol(money_symbol);
    if (has_embed_installed) setHasEmbedInstalled(has_embed_installed);
    setHasSubscribe(has_subscribe);
  };

  const loadCartData = async () => {
    const res = await cartService.getSettings();
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
      protectifyVisibility: String(data.show_cart),
      iconVisibility: String(data.show_cart_icon),
      selectButton: String(data.select_button),
      css: data.css || prev.css,
    }));

    // 更新pricing设置
    setPricingSettings(prev => ({
      ...prev,
      pricingType: String(data.pricing_type),
      pricingRule: String(data.pricing_rule),
      priceSelect: Array.isArray(data.price_select) ? data.price_select : prev.priceSelect,
      tiersSelect: Array.isArray(data.tiers_select) ? data.tiers_select : prev.tiersSelect,
      outPriceValue: String(data.out_price),
      outTierValue: String(data.out_tier),
      allPriceValue: String(data.all_price),
      allTiersValue: String(data.all_tiers),
    }));

    // 更新产品设置
    setProductSettings(prev => ({
      ...prev,
      selectedCollections: Array.isArray(data.product_collection) ? data.product_collection : prev.selectedCollections,
      icons: Array.isArray(data.icons) && data.icons.length > 0 ? data.icons : prev.icons,
      onlyInCollection: data.in_collection,
    }));

    //更新 fulfillment设置
    setFulfillmentSettings(prev => ({
      ...prev,
      fulfillmentRule: String(data.fulfillment_rule),
    }));
  };

  const validateFields = useCallback(() => {
    const newErrors: Record<string, string> = {};

    if (widgetSettings.protectifyVisibility === "1") {
      if (!widgetSettings.addonTitle.trim()) {
        newErrors.addonTitle = "Add-on title is required";
      }
      if (!widgetSettings.enabledDescription.trim()) {
        newErrors.enabledDescription = "Enabled description is required";
      }
      if (!widgetSettings.disabledDescription.trim()) {
        newErrors.disabledDescription = "Disabled description is required";
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

    try {
      const payload = {
        planTitle: widgetSettings.planTitle,
        iconVisibility: Number(widgetSettings.iconVisibility),
        protectifyVisibility: Number(widgetSettings.protectifyVisibility),
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
        outPrice: pricingSettings.outPriceValue,
        outTier: pricingSettings.outTierValue,
        allPrice: pricingSettings.allPriceValue,
        allTiers: pricingSettings.allTiersValue,
        selectedCollections: productSettings.selectedCollections,
        icons: productSettings.icons,
        onlyInCollection: productSettings.onlyInCollection,
        fulfillmentRule: Number(fulfillmentSettings.fulfillmentRule),
        css: widgetSettings.css,
      };

      const res = await cartService.updateSettings(payload);
      toastMessage(res?.code === 0 ? "Saved successfully" : "Saved Fail", 5000, res?.code !== 0);

      if (res?.code === 0) {
        setDirty(false);
        // 保存成功后，更新初始数据
        saveInitialData();
      }
    } catch (error) {
      console.error("Error saving settings:", error);
      toastMessage("Service Error", 5000, true);
    }
  }, [
    widgetSettings,
    pricingSettings,
    productSettings,
    fulfillmentSettings,
    toastMessage,
    validateFields,
    saveInitialData,
  ]);

  const discardChanges = useCallback(() => {
    if (initialDataRef.current) {
      // 还原到初始数据
      setWidgetSettings(deepClone(initialDataRef.current.widgetSettings));
      setPricingSettings(deepClone(initialDataRef.current.pricingSettings));
      setProductSettings(deepClone(initialDataRef.current.productSettings));

      // 清除错误和重置dirty状态
      setErrors({});
      setDirty(false);
    }
  }, []);

  // 在数据加载完成后保存初始数据
  useEffect(() => {
    if (!isLoading && initialDataRef.current === null) {
      saveInitialData();
    }
  }, [isLoading, saveInitialData]);

  // 监听数据变化，自动更新dirty状态
  useEffect(() => {
    if (initialDataRef.current) {
      updateDirtyState();
    }
  }, [widgetSettings, pricingSettings, productSettings, fulfillmentSettings, updateDirtyState]);

  useEffect(() => {
    void loadInitialData();
  }, []);

  return {
    // 状态
    widgetSettings,
    pricingSettings,
    productSettings,
    fulfillmentSettings,
    moneySymbol,
    hasSubscribe,
    hasEmbedInstalled,
    errors,
    isLoading,
    dirty,

    // 更新函数
    setWidgetSettings,
    setPricingSettings,
    setProductSettings,
    setFulfillmentSettings,
    setErrors,

    // 操作函数
    saveSettings,
    discardChanges,
    validateFields,
  };
}
