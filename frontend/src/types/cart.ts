export interface WidgetSettings {
  planTitle: string;
  iconVisibility: string;
  protectifyVisibility: string;
  selectButton: string;
  addonTitle: string;
  enabledDescription: string;
  disabledDescription: string;
  footerText: string;
  footerUrl: string;
  optInColor: string;
  optOutColor: string;
}

export interface PricingSettings {
  pricingType: string;
  pricingRule: string;
  priceSelect: PriceRange[];
  tiersSelect: TierRange[];
  restValuePrice: string;
  allPriceValue: string;
  allTiersValue: string;
}

export interface ProductSettings {
  selectedCollections: ResourceItem[];
  icons: IconType[];
  onlyInCollection: boolean;
}

export interface PriceRange {
  min: string;
  max: string;
  price: string;
}

export interface TierRange {
  min: string;
  max: string;
  percentage: string;
}

export interface IconType {
  id: number;
  src: string;
  selected: boolean;
}

export interface CartSettingsHook {
  widgetSettings: WidgetSettings;
  pricingSettings: PricingSettings;
  productSettings: ProductSettings;
  moneySymbol: string;
  errors: Record<string, string>;
  isLoading: boolean;
  hasSubscribe: boolean;
  dirty: boolean;
  setWidgetSettings: (setter: (prev: WidgetSettings) => WidgetSettings) => void;
  setPricingSettings: (setter: (prev: PricingSettings) => PricingSettings) => void;
  setProductSettings: (setter: (prev: ProductSettings) => ProductSettings) => void;
  setErrors: (setter: (prev: Record<string, string>) => Record<string, string>) => void;
  saveSettings: () => Promise<void>;
  discardChanges: () => void;
  validateFields: () => boolean;
}

// 组件Props类型
export interface WidgetStyleCardProps {
  widgetSettings: WidgetSettings;
  icons: IconType[];
  onWidgetSettingsChange: (value: Partial<WidgetSettings>) => void;
  onIconClick: (id: number) => void;
  onIconUpload: (file: File) => void;
}

export interface ContentCardProps {
  widgetSettings: WidgetSettings;
  errors: Record<string, string>;
  onFieldChange: (value: Partial<WidgetSettings>) => void;
}

export interface PricingCardProps {
  pricingSettings: PricingSettings;
  moneySymbol: string;
  errors: Record<string, string>;
  onSettingsChange: (value: Partial<PricingSettings>) => void;
  onPricingChange: (index: number, field: string, value: string, type: "price" | "tier") => void;
  onAddItem: (type: "price" | "tier") => void;
  onDeleteItem: (index: number, type: "price" | "tier") => void;
}
export type ResourceItem = { id: number; title: string };
export interface ProductCardProps {
  productSettings: ProductSettings;
  onCollectionSelect: (value: ResourceItem) => void;
  onlyCollection: (value: boolean) => void;
  onRemoveCollection: (value: number) => void;
  onCollectionChange?: (resources: Array<ResourceItem>) => void;
}

export interface CartSettingsData {
  plan_title: string;
  addon_title: string;
  enabled_desc: string;
  disabled_desc: string;
  foot_text: string;
  foot_url: string;
  in_color: string;
  out_color: string;
  show_cart: number;
  show_cart_icon: number;
  select_button: number;
  pricing_type: number;
  price_rule: number;
  price_select: any[];
  tiers_select: any[];
  other_money: number;
  all_price: number;
  all_tiers: number;
  product_collection: any[];
  icons: any[];
  in_collection: boolean;
}

export interface UpdateCartSettingsParams {
  planTitle: string;
  iconVisibility: number;
  protectifyVisibility: number;
  selectButton: number;
  addonTitle: string;
  enabledDescription: string;
  disabledDescription: string;
  footerText: string;
  footerUrl: string;
  optInColor: string;
  optOutColor: string;
  pricingType: number;
  pricingRule: number;
  priceSelect: any[];
  tiersSelect: any[];
  restValuePrice: string;
  allPrice: string;
  allTiers: string;
  selectedCollections: any[];
  icons: any[];
  onlyInCollection: boolean;
}
