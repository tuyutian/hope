import React, { useTransition } from "react";
import { Banner, BlockStack, Box, ContextualSaveBar, Frame, Layout, Page, Text } from "@shopify/polaris";
import SkeletonScreen from "@/pages/cart/components/Skeleton";
import CartDemo from "@/pages/cart/components/CartDemo";
import PublishWidget from "@/pages/cart/components/PublishWidget";
import { useCartSettings } from "@/hooks/useCartSettings";
import ContentCard from "@/pages/cart/components/ContentCard.tsx";
import PricingCard from "@/pages/cart/components/PricingCard.tsx";
import WidgetStyleCard from "@/pages/cart/components/WidgetStyleCard.tsx";
import ProductCard from "@/pages/cart/components/ProductCard.tsx";

export default function ShippingProtectionSettings() {
  const [isPending, startTransition] = useTransition();
  const {
    widgetSettings,
    pricingSettings,
    productSettings,
    collectionOptions,
    moneySymbol,
    errors,
    isLoading,
    saveLoading,
    dirty,
    setWidgetSettings,
    setPricingSettings,
    setProductSettings,
    setErrors,
    saveSettings,
    markDirty,
    discardChanges,
  } = useCartSettings();

  // 字段变化处理器
  const handleFieldChange =
    (setter: (value: any) => void, field: string, transform?: (value: any) => any) => (value: any) => {
      const transformedValue = transform ? transform(value) : value;
      setter(transformedValue);
      if (errors[field]) {
        setErrors(prev => ({ ...prev, [field]: "" }));
      }
      markDirty();
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
    markDirty();
  };

  // 定价表格处理器
  const handlePricingChange = (index: number, field: string, value: string, type: "price" | "tier") => {
    const errorKey = `${type}_${field}_${index}`;

    if (type === "price") {
      setPricingSettings(prev => ({
        ...prev,
        priceSelect: prev.priceSelect.map((item, i) => (i === index ? { ...item, [field]: value } : item)),
      }));
    } else {
      setPricingSettings(prev => ({
        ...prev,
        tiersSelect: prev.tiersSelect.map((item, i) => (i === index ? { ...item, [field]: value } : item)),
      }));
    }

    if (errors[errorKey]) {
      setErrors(prev => ({ ...prev, [errorKey]: "" }));
    }
    markDirty();
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
    markDirty();
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
    markDirty();
  };

  // 集合选择处理
  const handleCollectionSelect = (value: string) => {
    setProductSettings(prev => ({
      ...prev,
      selectedCollections: [...new Set([...prev.selectedCollections, value])],
      collectionInput: "",
    }));
    markDirty();
  };

  const handleRemoveCollection = (valueToRemove: string) => {
    setProductSettings(prev => ({
      ...prev,
      selectedCollections: prev.selectedCollections.filter(val => val !== valueToRemove),
    }));
    markDirty();
  };

  if (isLoading) {
    return <SkeletonScreen />;
  }

  const selectedIcon = productSettings.icons.find(icon => icon.selected);

  return (
    <Frame>
      <div style={{ position: "relative" }}>
        {dirty && (
          <ContextualSaveBar
            message="Unsaved changes"
            saveAction={{
              onAction: () => startTransition(saveSettings),
              loading: saveLoading,
              disabled: isPending,
            }}
            discardAction={{
              onAction: discardChanges,
            }}
          />
        )}

        <Page
          title="Create Protection Plan (Cart Page)"
          primaryAction={{
            content: "Save Changes",
            onAction: () => startTransition(saveSettings),
            loading: saveLoading,
          }}
        >
          <BlockStack gap="400">
            <Banner title="App embed is not enabled" tone="warning" secondaryAction={{ content: "Enable embed" }}>
              <Text as="p">
                Shipping Protection widget were published, but the app embed does not appear to be enabled. Please
                enable that to display the widget on your storefront cart.
              </Text>
            </Banner>

            <Layout>
              <Layout.Section variant="oneHalf">
                <Box>
                  <BlockStack gap="400">
                    <PublishWidget
                      insuranceVisibility={widgetSettings.insuranceVisibility}
                      onInsuranceVisibilityChange={handleFieldChange(
                        (value: boolean) =>
                          setWidgetSettings(prev => ({
                            ...prev,
                            insuranceVisibility: value ? "1" : "0",
                          })),
                        "insuranceVisibility"
                      )}
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
                      collectionOptions={collectionOptions}
                      onProductTypeChange={handleFieldChange(
                        (value: string) =>
                          setProductSettings(prev => ({
                            ...prev,
                            productTypeInput: value,
                          })),
                        "productType"
                      )}
                      onCollectionInputChange={handleFieldChange(
                        (value: string) =>
                          setProductSettings(prev => ({
                            ...prev,
                            collectionInput: value,
                          })),
                        "collectionInput"
                      )}
                      onCollectionSelect={handleCollectionSelect}
                      onRemoveCollection={handleRemoveCollection}
                    />
                  </BlockStack>
                </Box>
              </Layout.Section>

              <Layout.Section variant="oneHalf">
                <Box>
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
                </Box>
              </Layout.Section>
            </Layout>
          </BlockStack>
        </Page>
      </div>
    </Frame>
  );
}
