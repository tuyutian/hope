import type { OptionDescriptor } from "@shopify/polaris/build/ts/src/types";

export interface WidgetSettings {
  planTitle: string;
  iconVisibility: string;
  insuranceVisibility: string;
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
  productTypeInput: string;
  selectedCollections: string[];
  collectionInput: string;
  icons: IconType[];
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
  collectionOptions: OptionDescriptor[];
  moneySymbol: string;
  errors: Record<string, string>;
  isLoading: boolean;
  saveLoading: boolean;
  dirty: boolean;
  setWidgetSettings: (setter: (prev: WidgetSettings) => WidgetSettings) => void;
  setPricingSettings: (setter: (prev: PricingSettings) => PricingSettings) => void;
  setProductSettings: (setter: (prev: ProductSettings) => ProductSettings) => void;
  setErrors: (setter: (prev: Record<string, string>) => Record<string, string>) => void;
  saveSettings: () => Promise<void>;
  markDirty: () => void;
  discardChanges: () => void;
  validateFields: () => boolean;
}

// 组件Props类型
export interface WidgetStyleCardProps {
  widgetSettings: WidgetSettings;
  icons: IconType[];
  onWidgetSettingsChange: (value: Partial<WidgetSettings>) => void;
  onIconClick: (id: number) => void;
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
  onPricingChange: (index: number, field: string, value: string, type: 'price' | 'tier') => void;
  onAddItem: (type: 'price' | 'tier') => void;
  onDeleteItem: (index: number, type: 'price' | 'tier') => void;
}

export interface ProductCardProps {
  productSettings: ProductSettings;
  collectionOptions: OptionDescriptor[];
  onProductTypeChange: (value: string) => void;
  onCollectionInputChange: (value: string) => void;
  onCollectionSelect: (value: string) => void;
  onRemoveCollection: (value: string) => void;
}