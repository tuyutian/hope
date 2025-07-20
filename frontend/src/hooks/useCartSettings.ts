import { useState, useEffect, useCallback, useRef } from "react";
import { GetUserConf, rqGetCartSetting, rqPostUpdateCartSetting } from "@/api";
import type { WidgetSettings, PricingSettings, ProductSettings, CartSettingsHook } from "@/types/cart";
import { getMessageState } from "@/stores/messageStore";
import type { OptionDescriptor } from "@shopify/polaris/build/ts/src/types";

export function useCartSettings(): CartSettingsHook {
  const toastMessage = getMessageState().toastMessage;

  // 添加 ref 来存储初始数据
  const initialDataRef = useRef<{
    widgetSettings: WidgetSettings;
    pricingSettings: PricingSettings;
    productSettings: ProductSettings;
  } | null>(null);

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
    selectProductTypes: [],
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

    return hasWidgetChanges || hasPricingChanges || hasProductChanges;
  }, [widgetSettings, pricingSettings, productSettings]);

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
    };
  }, [widgetSettings, pricingSettings, productSettings]);

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
      selectProductTypes: data.product_type || prev.selectProductTypes,
      selectedCollections: Array.isArray(data.product_collection) ? data.product_collection : prev.selectedCollections,
      icons: Array.isArray(data.icons) && data.icons.length > 0 ? data.icons : prev.icons,
    }));
  };

  const validateFields = useCallback(() => {
    const newErrors: Record<string, string> = {};

    if (widgetSettings.insuranceVisibility === "1") {
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
        selectProductTypes: productSettings.selectProductTypes,
        selectedCollections: productSettings.selectedCollections,
        icons: productSettings.icons,
      };

      const res = await rqPostUpdateCartSetting(payload);
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
  }, [widgetSettings, pricingSettings, productSettings, toastMessage, validateFields, saveInitialData]);

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
  }, [widgetSettings, pricingSettings, productSettings, updateDirtyState]);

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
    dirty,

    // 更新函数
    setWidgetSettings,
    setPricingSettings,
    setProductSettings,
    setErrors,

    // 操作函数
    saveSettings,
    discardChanges,
    validateFields,
  };
}
