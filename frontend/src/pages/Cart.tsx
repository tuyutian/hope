import {
  Banner,
  BlockStack,
  Box,
  Card,
  ContextualSaveBar,
  Divider,
  Frame,
  InlineStack,
  Layout,
  Page,
  Select,
  Text,
  TextField,
} from "@shopify/polaris";
import {useCallback, useEffect, useState, useTransition} from "react";
import {GetUserConf, rqGetCartSetting, rqPostUpdateCartSetting} from "@/api";
import SkeletonScreen from "@/pages/cart/components/Skeleton";
import SketchPickerWithInput from "@/components/form/SketchPickerWithInput.tsx";
import type {OptionDescriptor} from "@shopify/polaris/build/ts/src/types";

import IconSelector from "@/pages/cart/components/IconSelector";
import PricingTable from "@/pages/cart/components/PricingTable";
import CollectionSelector from "@/pages/cart/components/CollectionSelector";
import CartDemo from "@/pages/cart/components/CartDemo";
import PublishWidget from "@/pages/cart/components/PublishWidget";
import {getMessageState} from "@/stores/messageStore.ts";

export default function ShippingProtectionSettings() {
  const toastMessage = getMessageState().toastMessage;
  // All existing state variables remain the same
  const [planTitle, setPlanTitle] = useState("Plan Title");
  const [iconVisibility, setIconVisibility] = useState("0");
  const [insuranceVisibility, setInsuranceVisibility] = useState("0");
  const [selectButton, setSelectButton] = useState("0");
  const [addonTitle, setAddonTitle] = useState("Shipping Protection");
  const [enabledDescription, setEnabledDescription] = useState("After purchasing this insurance, we will resolve all after-sales issues related to this order for you.");
  const [disabledDescription, setDisabledDescription] = useState("After purchasing this insurance, we will resolve all after-sales issues related to this order for you.");
  const [footerText, setFooterText] = useState("");
  const [footerUrl, setFooterUrl] = useState("");
  const [toastContent, setToastContent] = useState("");
  const [toastError, setToastError] = useState(false);
  const [optInColor, setOptInColor] = useState("#fffff");
  const [optOutColor, setOptOutColor] = useState("#fffff");
  const [toastActive, setToastActive] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});

  // Protection Pricing state
  const [pricingType, setPricingType] = useState("0");
  const [pricingRule, setPricingRule] = useState("0");
  const [priceSelect, setPriceSelect] = useState([
    {min: "1.00", max: "100.00", price: "3.00"},
    {min: "101.00", max: "1000.00", price: "10.00"},
  ]);

  const [tiersSelect, setTiersSelect] = useState([
    {min: "1.00", max: "100.00", percentage: "10"},
    {min: "101.00", max: "1000.00", percentage: "20"},
  ]);

  const [restValuePrice, setRestValuePrice] = useState("10.00");

  const [allPriceValue, setAllPriceValue] = useState("");
  const [allTiersValue, setAllTiersValue] = useState("");

  const [productTypeInput, setProductTypeInput] = useState("");

  const [collectionInput, setCollectionInput] = useState("");
  const [collectionOptions, setCollectionOptions] = useState<OptionDescriptor[]>([]);
  const [moneySymbol, setMoneySymbol] = useState("$");


  const [selectedCollections, setSelectedCollections] = useState<string[]>([]);
  const [checkboxInput, setCheckboxInput] = useState(false);
  const [switchValue, setSwitchValue] = useState(false);

  const [isLoading, setIsLoading] = useState(true);
  const [saveLoading, setSaveLoading] = useState(false);
  const [dirty, setDirty] = useState(false);


  const [icons, setIcons] = useState([
    {
      id: 1,
      src: "https://img.icons8.com/color/48/shield.png",
      selected: true,
    },
    {
      id: 2,
      src: "https://maxst.icons8.com/vue-static/faceswapper/hero/faces/2.jpg",
      selected: false,
    },
  ]);

  const [isPending, startTransition] = useTransition();

  // All existing handlers remain the same
  const handleIconClick = (id: number) => {
    setIcons(icons.map(icon => ({
      ...icon,
      selected: icon.id === id,
    })));
    handleOnDiscard();
  };

  const selectedIcon = icons.find(icon => icon.selected);

  // All existing validation and handler functions remain the same...

  // 验证函数
  const validateFields = () => {
    const newErrors: Record<string, string> = {};
    if (insuranceVisibility === "1") {

      if (pricingType === "1" && pricingRule === "1") {
        // 检查每个Tier
        tiersSelect.forEach((tier, index) => {
          if (tier.min || tier.max || tier.percentage) {
            if (!tier.min) newErrors[`tier_min_${index}`] = "Please fill in";
            if (!tier.max) newErrors[`tier_max_${index}`] = "Please fill in";
            if (!tier.percentage) newErrors[`tier_percentage_${index}`] = "Please fill in";
          }
        });
      }

      if (pricingType === "0" && pricingRule === "1") {
        // 检查每个Price
        priceSelect.forEach((price, index) => {
          if (price.min || price.max || price.price) {
            if (!price.min) newErrors[`price_min_${index}`] = "Please fill in";
            if (!price.max) newErrors[`price_max_${index}`] = "Please fill in";
            if (!price.price) newErrors[`price_price_${index}`] = "Please fill in";
          }
        });
      }

      //if (!planTitle.trim()) newErrors.planTitle = 'Plan Title is required';
      if (!addonTitle.trim()) newErrors.addonTitle = "Add-on Title is required";
      if (!enabledDescription.trim()) newErrors.enabledDescription = "Enabled Description is required";
      if (!disabledDescription.trim()) newErrors.disabledDescription = "Disabled Description is required";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  // 字段变化时清除错误
  const handleFieldChange = (setter: (value: string) => void, field: string) => (value: string) => {
    setter(value);
    if (errors[field]) {
      setErrors(prev => ({...prev, [field]: ""}));
    }
    handleOnDiscard();
  };

  const removeTag = (valueToRemove: string) => {
    setSelectedCollections(prev => prev.filter((val) => val !== valueToRemove));
    handleOnDiscard();
  };

  // All existing useEffect and API functions remain the same...
  useEffect(() => {
    void getUserInfoData();
    void getCartData();
  }, []);

  const getUserInfoData = async () => {
    const res = await GetUserConf();
    if (res.code !== 0 || !res.data) return;
    const data = res.data;
    if (Array.isArray(data.collections)) {
      const collections = data.collections.map((collection: {
        label: string;
        value: number;
      }) => ({label: collection.label, value: String(collection.value)}));
      setCollectionOptions(collections);
    }
    if (data.money_symbol) setMoneySymbol(data.money_symbol);
  };

  const getCartData = async () => {
    try {
      setIsLoading(true);
      const res = await rqGetCartSetting();
      if (res.code !== 0 || !res.data) return;
      const data = res.data;

      // Set all state from API response (same as before)
      if (data.plan_title) setPlanTitle(data.plan_title);
      if (data.addon_title) setAddonTitle(data.addon_title);
      if (data.enabled_desc) setEnabledDescription(data.enabled_desc);
      if (data.disabled_desc) setDisabledDescription(data.disabled_desc);
      if (data.foot_text) setFooterText(data.foot_text);
      if (data.foot_url) setFooterUrl(data.foot_url);
      if (data.in_color) setOptInColor(data.in_color);
      if (data.out_color) setOptOutColor(data.out_color);
      if (typeof data.other_money === "number") setRestValuePrice(String(data.other_money));

      if (typeof data.all_price === "number") setAllPriceValue(String(data.all_price));
      if (typeof data.all_tiers === "number") setAllTiersValue(String(data.all_tiers));

      if (typeof data.show_cart === "number") setInsuranceVisibility(String(data.show_cart));
      if (typeof data.show_cart_icon === "number") setIconVisibility(String(data.show_cart_icon));
      if (typeof data.select_button === "number") setSelectButton(String(data.select_button));

      if (data.product_type) setProductTypeInput(data.product_type);

      if (Array.isArray(data.product_collection)) setSelectedCollections(data.product_collection);

      if (typeof data.pricing_type === "number") setPricingType(String(data.pricing_type));

      if (typeof data.price_rule === "number") setPricingRule(String(data.price_rule));

      if (Array.isArray(data.price_select)) setPriceSelect(data.price_select);

      if (Array.isArray(data.tiers_select)) setTiersSelect(data.tiers_select);

      if (Array.isArray(data.icons) && data.icons.length > 0) setIcons(data.icons);

    } catch (error) {
      console.error("Error fetching cart data:", error);
    } finally {
      setIsLoading(false);
    }

  };

  // All existing handler functions remain the same...
  // New handlers for Protection Pricing
  const handleAddTier = () => {
    if (tiersSelect.length < 20) {
      setTiersSelect([...tiersSelect, {min: "", max: "", percentage: ""}]);
      handleOnDiscard();
    }
  };

  const handleAddPrice = () => {
    if (priceSelect.length < 20) {
      setPriceSelect([...priceSelect, {min: "", max: "", price: ""}]);
      handleOnDiscard();
    }
  };


  const handleSave = useCallback(async () => {
    if (!validateFields()) {
      setToastActive(true);
      return;
    }

    setSaveLoading(true);
    const payload = {
      planTitle,
      iconVisibility: Number(iconVisibility),
      insuranceVisibility: Number(insuranceVisibility),
      selectButton: Number(selectButton),
      addonTitle,
      enabledDescription,
      disabledDescription,
      footerText,
      footerUrl,
      optInColor,
      optOutColor,
      pricingType: Number(pricingType),
      pricingRule: Number(pricingRule),
      priceSelect,
      tiersSelect,
      restValuePrice,
      allPrice: allPriceValue,
      allTiers: allTiersValue,
      productTypeInput,
      selectedCollections,
      icons
    };

    startTransition(async function () {
      try {

        const res = await rqPostUpdateCartSetting(payload);
        setToastContent(res?.code === 0 ? "Saved successfully" : "Saved Fail");
        setToastError(res?.code !== 0);
        setToastActive(true);

      } catch (error) {
        console.error("Error saving settings:", error);
        setToastContent("Service Error");
        setToastError(true);
        setToastActive(true);
      } finally {
        setSaveLoading(false);
      }
    });
    handleDiscard();
  }, [
    planTitle, iconVisibility, insuranceVisibility, selectButton, addonTitle,
    enabledDescription, disabledDescription, footerText, footerUrl, optInColor,
    optOutColor, pricingType, pricingRule, priceSelect, tiersSelect, restValuePrice,
    allPriceValue, allTiersValue, productTypeInput, selectedCollections, icons,
  ]);

  useEffect(() => {
    if (toastActive) {
      toastMessage(Object.keys(errors).length ? "Please fix validation errors" : toastContent, 5000, toastError);
    }
  }, [errors, toastActive, toastContent, toastError]);

  const handleTierChange = (index: number, field: string, value: string, error: string) => {
    let validValue = value;

    // 如果值为空字符串，允许删除（空值也有效）
    if (validValue === "") {
      const updatedTiers = [...tiersSelect];
      updatedTiers[index] = {...updatedTiers[index], [field]: ""};
      setTiersSelect(updatedTiers);
    } else {
      // 正则表达式验证输入是数字并且不超过100
      validValue = validValue.replace(/[^0-9.]/g, "");
      const validNumber = /^[0-9]+(\.[0-9]{0,2})?$/;

      if (validNumber.test(validValue) && (parseFloat(validValue) <= 100)) {
        const updatedTiers = [...tiersSelect];
        updatedTiers[index] = {...updatedTiers[index], [field]: validValue};
        setTiersSelect(updatedTiers);
      }
    }
    handleOnDiscard();
    if (errors[error]) {
      setErrors(prev => ({...prev, [error]: ""}));
    }
  };

  const handleDeleteTier = (index: number) => {
    const updatedTiers = tiersSelect.filter((_, i) => i !== index);
    setTiersSelect(updatedTiers);
    handleOnDiscard();
  };

  const handleDeletePrice = (index: number) => {
    const updatedPrices = priceSelect.filter((_, i) => i !== index);
    setPriceSelect(updatedPrices);
    handleOnDiscard();
  };

  const handlePriceChange = (index: number, field: string, value: string, error: string) => {
    let validValue = value;

    // 如果值为空字符串，允许删除
    if (validValue === "") {
      const updatedPrices = [...priceSelect];
      updatedPrices[index] = {...updatedPrices[index], [field]: ""};
      setPriceSelect(updatedPrices);
    } else {
      // 正则表达式验证输入是数字并且不超过100
      validValue = validValue.replace(/[^0-9.]/g, "");
      const validNumber = /^[0-9]+(\.[0-9]{0,2})?$/;

      if (validNumber.test(validValue)) {
        const updatedPrices = [...priceSelect];
        updatedPrices[index] = {...updatedPrices[index], [field]: validValue};
        setPriceSelect(updatedPrices);
      }
    }
    handleOnDiscard();
    if (errors[error]) {
      setErrors(prev => ({...prev, [error]: ""}));
    }
  };

  // 切换勾选状态
  const handleCheckbox = () => {
    handleOnDiscard();
    setCheckboxInput(prev => !prev);
  };

// 切换 switch 状态
  const onChangeSwitch = (value: boolean) => {
    setSwitchValue(value);
  };

  const onInsuranceVisibility = (value: boolean) => {
    handleOnDiscard();
    setInsuranceVisibility(value ? "1" : "0");
  };

  const handleDiscard = () => {
    setDirty(false);
  };

  const handleOnDiscard = () => {
    if (!dirty) setDirty(true);
  };

  if (isLoading) {
    return <SkeletonScreen />;
  }

  return (
    <Frame>
      <div style={{position: "relative"}}>
        {dirty && (
          <ContextualSaveBar
            message="Unsaved changes"
            saveAction={{
              onAction: handleSave,
              loading: saveLoading,
              disabled: isPending,
            }}
            discardAction={{
              onAction: handleDiscard,
            }}
          />
        )}

        <Page
          title="Create Protection Plan (Cart Page)"
          primaryAction={{content: "Save Changes", onAction: handleSave, loading: saveLoading}}
        >
          <BlockStack gap="400">
            <Banner title="App embed is not enabled" tone="warning" secondaryAction={{
              content: "Enable embed"
            }}>
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
                      insuranceVisibility={insuranceVisibility}
                      onInsuranceVisibilityChange={onInsuranceVisibility}
                    />

                    <Card padding="400">
                      <BlockStack gap="300">
                        <Text variant="headingSm" as="h6">Widget Style</Text>

                        <BlockStack gap="200">
                          <Select
                            label="Icon Visibility"
                            options={[
                              {label: "Show Icon", value: "1"},
                              {label: "Hide Icon", value: "0"},
                            ]}
                            value={iconVisibility}
                            onChange={(value) => {
                              setIconVisibility(value);
                              handleOnDiscard();
                            }}
                          />
                        </BlockStack>

                        <BlockStack gap="200">
                          <IconSelector
                            icons={icons}
                            onIconClick={handleIconClick}
                          />
                        </BlockStack>

                        <BlockStack gap="200">
                          <Select
                            options={[
                              {label: "Switch", value: "0"},
                              {label: "Checkbox", value: "1"},
                            ]}
                            value={selectButton}
                            onChange={(value) => {
                              handleOnDiscard();
                              setSelectButton(value);
                            }}
                            label="Select Button"
                          />
                        </BlockStack>

                        <InlineStack align="space-between" gap="300" wrap={false}>
                          <Box>
                            <div>Opt-in action button</div>
                            <SketchPickerWithInput defaultColor={optInColor} onChange={setOptInColor} />

                          </Box>
                          <Box>
                            <div>Opt-out action button</div>
                            <SketchPickerWithInput defaultColor={optOutColor} onChange={setOptOutColor} />
                          </Box>
                        </InlineStack>
                      </BlockStack>
                    </Card>

                    <Card padding="400">
                      <BlockStack gap="200">
                        <Text variant="headingSm" as="h6">Content</Text>
                        <TextField
                          autoComplete="off"
                          label="Add-on title"
                          value={addonTitle}
                          onChange={handleFieldChange(setAddonTitle, "addonTitle")}
                          error={errors.addonTitle}
                          maxLength={50}
                        />
                        <TextField
                          autoComplete="off"

                          label="Enabled description"
                          value={enabledDescription}
                          onChange={handleFieldChange(setEnabledDescription, "enabledDescription")}
                          multiline={4}
                          maxLength={200}
                          error={errors.enabledDescription}
                        />
                        <TextField
                          autoComplete="off"
                          label="Disabled description"
                          value={disabledDescription}
                          onChange={handleFieldChange(setDisabledDescription, "disabledDescription")}
                          multiline={4}
                          maxLength={200}
                          error={errors.disabledDescription}
                        />
                        <TextField
                          autoComplete="off"
                          label="Footer link text"
                          value={footerText}
                          onChange={setFooterText}
                          maxLength={50}
                        />
                        <TextField
                          autoComplete="off"
                          label="Footer link URL"
                          value={footerUrl}
                          onChange={setFooterUrl}
                          maxLength={150}
                        />
                        <Text tone="subdued" variant="bodySm" as="span">Note: leave blank for no link</Text>
                      </BlockStack>
                    </Card>

                    {/* Protection Pricing Card */}
                    <Card padding="400">
                      <BlockStack gap="400">
                        <Text variant="headingSm" as="h6">Protection Pricing</Text>
                        <BlockStack gap="200">
                          <Select
                            label="Pricing type"
                            options={[
                              {label: "Fixed", value: "0"},
                              {label: "Pricing rule", value: "1"},
                            ]}
                            value={pricingType}
                            onChange={(value) => {
                              handleOnDiscard();
                              setPricingType(value);
                            }}
                          />
                        </BlockStack>


                        <BlockStack gap="200">
                          <Select
                            label="Pricing rule"
                            options={[
                              {label: "App for all cart vale", value: "0"},
                              {label: "Apply for different cart value range", value: "1"},
                            ]}
                            value={pricingRule}
                            onChange={(value) => {
                              handleOnDiscard();
                              setPricingRule(value);
                            }}
                          />
                        </BlockStack>

                        {pricingRule === "1" && (
                          <PricingTable
                            pricingType={pricingType}
                            priceSelect={priceSelect}
                            tiersSelect={tiersSelect}
                            moneySymbol={moneySymbol}
                            errors={errors}
                            onPriceChange={handlePriceChange}
                            onTierChange={handleTierChange}
                            onDeletePrice={handleDeletePrice}
                            onDeleteTier={handleDeleteTier}
                            onAddPrice={handleAddPrice}
                            onAddTier={handleAddTier}
                          />
                        )}

                        {pricingRule === "1" ? (
                          <TextField
                            autoComplete="off"
                            label="Other Value Range"
                            value={restValuePrice}
                            onChange={(value) => {
                              handleOnDiscard();
                              setRestValuePrice(value);
                            }}
                            prefix="$"
                          />
                        ) : (
                          <TextField
                            label=""
                            autoComplete="off"
                            value={pricingType === "1" ? allTiersValue : allPriceValue}
                            onChange={(value) => {
                              handleOnDiscard();
                              if (pricingType === "1") {
                                setAllTiersValue(value);
                              } else {
                                setAllPriceValue(value);
                              }
                            }}
                            prefix={pricingType === "1" ? "%" : moneySymbol}
                          />
                        )}


                        <Divider />

                      </BlockStack>
                    </Card>

                    <Card padding="400">
                      <BlockStack gap="200">

                        <TextField
                          autoComplete="off"
                          label="Product Types"
                          value={productTypeInput}
                          onChange={setProductTypeInput}
                          maxLength={50}
                        />
                        <CollectionSelector
                          collectionOptions={collectionOptions}
                          selectedCollections={selectedCollections}
                          collectionInput={collectionInput}
                          onCollectionInputChange={setCollectionInput}
                          onCollectionSelect={(value) => {
                            setSelectedCollections([...new Set([...selectedCollections, value])]);
                            setCollectionInput("");
                          }}
                          onRemoveCollection={removeTag}
                        />
                      </BlockStack>
                    </Card>

                  </BlockStack>
                </Box>
              </Layout.Section>

              <Layout.Section variant="oneHalf">
                <Box>
                  <CartDemo
                    iconVisibility={iconVisibility}
                    selectedIcon={selectedIcon}
                    addonTitle={addonTitle}
                    enabledDescription={enabledDescription}
                    disabledDescription={disabledDescription}
                    footerText={footerText}
                    footerUrl={footerUrl}
                    selectButton={selectButton}
                    switchValue={switchValue}
                    checkboxInput={checkboxInput}
                    optInColor={optInColor}
                    optOutColor={optOutColor}
                    onSwitchChange={onChangeSwitch}
                    onCheckboxChange={handleCheckbox}
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