import React, { useTransition } from "react";
import { Banner, BlockStack, Layout, Page, Text } from "@shopify/polaris";
import SkeletonScreen from "@/pages/cart/components/Skeleton";
import CartDemo from "@/pages/cart/components/CartDemo";
import PublishWidget from "@/pages/cart/components/PublishWidget";
import { useCartSettings } from "@/hooks/useCartSettings";
import ContentCard from "@/pages/cart/components/ContentCard.tsx";
import PricingCard from "@/pages/cart/components/PricingCard.tsx";
import WidgetStyleCard from "@/pages/cart/components/WidgetStyleCard.tsx";
import ProductCard from "@/pages/cart/components/ProductCard.tsx";
import PageSaveBar from "@/components/form/PageSaveBar.tsx";
import "@/styles/cart.css";
import { ResourceItem } from "@/types/cart.ts";
import { userService } from "@/services/user";

export default function ShippingProtectionSettings() {
  const {
    widgetSettings,
    pricingSettings,
    productSettings,
    moneySymbol,
    errors,
    isLoading,
    dirty,
    setWidgetSettings,
    setPricingSettings,
    setProductSettings,
    setErrors,
    saveSettings,
    hasSubscribe,
    discardChanges,
  } = useCartSettings();
  const [isPending, startTransition] = useTransition();
  const [isHanding, startSubscription] = useTransition();
  // 字段变化处理器
  const handleFieldChange =
    (setter: (value: any) => void, field: string, transform?: (value: any) => any) => (value: any) => {
      const transformedValue = transform ? transform(value) : value;
      setter(transformedValue);
      if (errors[field]) {
        setErrors(prev => ({ ...prev, [field]: "" }));
      }
    };

  // 图标选择处理
  const handleIconClick = (id: number) => {
    setProductSettings(prev => ({
      ...prev,
      icons: prev.icons.map(icon => ({
        ...icon,
        selected: icon.id === id,
      })),
    }));
  };

  // 定价表格处理器
  const handlePricingChange = (index: number, field: string, value: string, type: "price" | "tier") => {
    const errorKey = `${type}_${field}_${index}`;

    if (type === "price") {
      setPricingSettings(prev => ({
        ...prev,
        priceSelect: prev.priceSelect.map((item, i) =>
          i === index
            ? {
                ...item,
                [field]: value,
              }
            : item
        ),
      }));
    } else {
      setPricingSettings(prev => ({
        ...prev,
        tiersSelect: prev.tiersSelect.map((item, i) =>
          i === index
            ? {
                ...item,
                [field]: value,
              }
            : item
        ),
      }));
    }

    if (errors[errorKey]) {
      setErrors(prev => ({ ...prev, [errorKey]: "" }));
    }
  };

  // 添加/删除定价项
  const handleAddPricingItem = (type: "price" | "tier") => {
    if (type === "price") {
      setPricingSettings(prev => ({
        ...prev,
        priceSelect: [...prev.priceSelect, { min: "", max: "", price: "" }],
      }));
    } else {
      setPricingSettings(prev => ({
        ...prev,
        tiersSelect: [...prev.tiersSelect, { min: "", max: "", percentage: "" }],
      }));
    }
  };

  const handleDeletePricingItem = (index: number, type: "price" | "tier") => {
    if (type === "price") {
      setPricingSettings(prev => ({
        ...prev,
        priceSelect: prev.priceSelect.filter((_, i) => i !== index),
      }));
    } else {
      setPricingSettings(prev => ({
        ...prev,
        tiersSelect: prev.tiersSelect.filter((_, i) => i !== index),
      }));
    }
  };

  // 集合选择处理
  const handleCollectionSelect = (value: ResourceItem) => {
    setProductSettings(prev => ({
      ...prev,
      selectedCollections: [...new Set([...prev.selectedCollections, value])],
      collectionInput: "",
    }));
  };

  const handleOnlyCollection = (check: boolean) => {
    setProductSettings(prev => ({
      ...prev,
      onlyInCollection: check,
    }));
  };
  // 集合选择处理
  const handleRemoveCollection = (valueToRemove: number) => {
    setProductSettings(prev => ({
      ...prev,
      selectedCollections: prev.selectedCollections.filter(val => val.id !== valueToRemove),
    }));
  };

  // 处理ResourcePicker选择的集合
  const handleCollectionChange = (resources: ResourceItem[]) => {
    setProductSettings(prev => ({
      ...prev,
      selectedCollections: [...new Set([...prev.selectedCollections, ...resources])],
    }));
  };

  if (isLoading) {
    return <SkeletonScreen />;
  }

  const selectedIcon = productSettings.icons.find(icon => icon.selected);

  function handlePublishWidget() {
    if (!hasSubscribe) {
      startSubscription(async function () {
        const res = await userService.startSubscription();
        if (res.code === 0) {
          startSubscription(function () {
            open(res.data, "_self");
          });
        }
      });
      return;
    }
    handleFieldChange(
      (value: boolean) =>
        setWidgetSettings(prev => ({
          ...prev,
          insuranceVisibility: value ? "1" : "0",
        })),
      "insuranceVisibility"
    );
  }

  return (
    <Page
      title="Create Protection Plan (Cart Page)"
      primaryAction={{
        content: "Save Changes",
        onAction: () => startTransition(saveSettings),
        loading: isPending,
      }}
    >
      <PageSaveBar dirty={dirty} onSave={saveSettings} onDiscard={discardChanges} />
      <BlockStack gap="400">
        <Banner title="App embed is not enabled" tone="warning" secondaryAction={{ content: "Enable embed" }}>
          <Text as="p">
            Shipping Protection widget were published, but the app embed does not appear to be enabled. Please enable
            that to display the widget on your storefront cart.
          </Text>
        </Banner>

        <Layout>
          <Layout.Section variant="oneHalf">
            <div className="settings-scroll-container">
              <BlockStack gap="400">
                {/* 您的设置组件 */}
                <PublishWidget
                  loading={isHanding}
                  insuranceVisibility={widgetSettings.insuranceVisibility}
                  onInsuranceVisibilityChange={handlePublishWidget}
                />

                <WidgetStyleCard
                  widgetSettings={widgetSettings}
                  icons={productSettings.icons}
                  onWidgetSettingsChange={handleFieldChange(
                    (value: any) => setWidgetSettings(prev => ({ ...prev, ...value })),
                    "widgetSettings"
                  )}
                  onIconClick={handleIconClick}
                />

                <ContentCard
                  widgetSettings={widgetSettings}
                  errors={errors}
                  onFieldChange={handleFieldChange(
                    (value: any) => setWidgetSettings(prev => ({ ...prev, ...value })),
                    "content"
                  )}
                />

                <PricingCard
                  pricingSettings={pricingSettings}
                  moneySymbol={moneySymbol}
                  errors={errors}
                  onSettingsChange={handleFieldChange(
                    (value: any) => setPricingSettings(prev => ({ ...prev, ...value })),
                    "pricing"
                  )}
                  onPricingChange={handlePricingChange}
                  onAddItem={handleAddPricingItem}
                  onDeleteItem={handleDeletePricingItem}
                />

                <ProductCard
                  productSettings={productSettings}
                  onCollectionSelect={handleCollectionSelect}
                  onlyCollection={handleOnlyCollection}
                  onRemoveCollection={handleRemoveCollection}
                  onCollectionChange={handleCollectionChange}
                />
              </BlockStack>
            </div>
          </Layout.Section>

          <Layout.Section variant="oneHalf">
            <div className="preview-sticky-container">
              <CartDemo
                iconVisibility={widgetSettings.iconVisibility}
                selectedIcon={selectedIcon}
                addonTitle={widgetSettings.addonTitle}
                enabledDescription={widgetSettings.enabledDescription}
                disabledDescription={widgetSettings.disabledDescription}
                footerText={widgetSettings.footerText}
                footerUrl={widgetSettings.footerUrl}
                selectButton={widgetSettings.selectButton}
                optInColor={widgetSettings.optInColor}
                optOutColor={widgetSettings.optOutColor}
                switchValue={false}
                checkboxInput={false}
                onSwitchChange={() => {}}
                onCheckboxChange={() => {}}
              />
            </div>
          </Layout.Section>
        </Layout>
      </BlockStack>
    </Page>
  );
}
