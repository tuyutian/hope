import {
  Autocomplete,
  BlockStack,
  Button,
  Card,
  ColorPicker,
  ContextualSaveBar,
  DataTable,
  Divider,
  hsbToHex,
  InlineStack,
  Layout,
  Page,
  Select,
  Tag,
  Text,
  TextField,
  Thumbnail,
  Toast,
} from "@shopify/polaris";
import {useCallback, useEffect, useRef, useState} from "react";
import {rqGetCartSetting, GetUserConf, rqPostUpdateCartSetting} from "@/api";
import SkeletonScreen from "@/pages/Cart/Skeleton";


export default function ShippingProtectionSettings() {
  // Existing state
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
  const [optInColor, setOptInColor] = useState({
    hue: 0,
    brightness: 0,
    saturation: 0,
  });
  const [optOutColor, setOptOutColor] = useState({
    hue: 0,
    brightness: 0.9,
    saturation: 0,
  });
  const [editingColor, setEditingColor] = useState(null);
  const colorPickerRef = useRef(null);
  const [toastActive, setToastActive] = useState(false);
  const [errors, setErrors] = useState({});

  // New state for Protection Pricing
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
  const [collectionOptions, setCollectionOptions] = useState([]);
  const [moneySymbol, setMoneySymbol] = useState("$");


  const [selectedCollections, setSelectedCollections] = useState([]);
  const [checkboxInput, setCheckboxInput] = useState(false);
  const [switchValue, setSwitchValue] = useState(false);

  const [isLoading, setIsLoading] = useState(true);
  const [saveLoading, setSaveLoading] = useState(false);
  const [dirty, setDirty] = useState(false);


  const [icons, setIcons] = useState([
    {
      id: 1,
      src: "https://img.icons8.com/color/48/shield.png",
      selected: true, // 默认第一个为选中
    },
    {
      id: 2,
      src: "https://maxst.icons8.com/vue-static/faceswapper/hero/faces/2.jpg",
      selected: false,
    },
  ]);

  const handleIconClick = (id) => {
    setIcons(icons.map(icon => ({
      ...icon,
      selected: icon.id === id,
    })));
    handleOnDiscard();
  };

  const selectedIcon = icons.find(icon => icon.selected);


  // 验证函数
  const validateFields = () => {
    const newErrors = {};
    if (insuranceVisibility === "1") {

      if (pricingType === "1" && pricingRule === "1") {
        // 检查每个Tier
        tiersSelect.forEach((tier, index) => {
          if (tier.min || tier.max || tier.price) {
            if (!tier.min) {
              newErrors[`tier_min_${index}`] = "Please fill in";
            }
            if (!tier.max) {
              newErrors[`tier_max_${index}`] = "Please fill in";
            }
            if (!tier.price) {
              newErrors[`tier_percentage_${index}`] = "Please fill in";
            }
          }
        });
      }

      if (pricingType === "0" && pricingRule === "1") {
        // 检查每个Price
        priceSelect.forEach((price, index) => {
          if (price.min || price.max || price.price) {
            if (!price.min) {
              newErrors[`price_min_${index}`] = "Please fill in";
            }
            if (!price.max) {
              newErrors[`price_max_${index}`] = "Please fill in";
            }
            if (!price.price) {
              newErrors[`price_price_${index}`] = "Please fill in";
            }
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
  const handleFieldChange = (setter, field) => (value) => {
    setter(value);
    if (errors[field]) {
      setErrors(prev => ({...prev, [field]: undefined}));
    }
    handleOnDiscard();
  };

  const removeTag = (valueToRemove, selected, setSelected) => {
    setSelected(selected.filter((val) => val !== valueToRemove));
    handleOnDiscard();
  };

  useEffect(() => {
    void getUserInfoData();
    void getCartData();  // 只在组件挂载时执行一次

    // 可选的清理函数：如果需要清理操作，可以在这里进行
    return () => {
      // 清理逻辑 (如果需要)
    };
  }, []);  // 空数组确保只在挂载时执行

  // Existing handlers
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (colorPickerRef.current && !colorPickerRef.current.contains(event.target)) {
        setEditingColor(null);
      }
    };
    if (editingColor) {
      document.addEventListener("mousedown", handleClickOutside);
    } else {
      document.removeEventListener("mousedown", handleClickOutside);
    }
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [editingColor]);

  const getUserInfoData = async (value) => {
    const res = await GetUserConf();
    if (res.code !== 200 || !res.data) return;
    console.log(res.data);
    const data = res.data;
    if (Array.isArray(data.collections)) setCollectionOptions(data.collections);
    if (data.money_symbol) setMoneySymbol(data.money_symbol);
  };

  const getCartData = async (value) => {
    try {
      setIsLoading(true);
      const res = await rqGetCartSetting();
      if (res.code !== 200 || !res.data) return;
      const data = res.data;
      if (data.plan_title) setPlanTitle(data.plan_title);
      if (data.addon_title) setAddonTitle(data.addon_title);
      if (data.enabled_desc) setEnabledDescription(data.enabled_desc);
      if (data.disabled_desc) setDisabledDescription(data.disabled_desc);
      if (data.foot_text) setFooterText(data.foot_text);
      if (data.foot_url) setFooterUrl(data.foot_url);
      if (data.in_color) setOptInColor(data.in_color); // 需要把HEX转HSB
      if (data.out_color) setOptOutColor(data.out_color); // 需要把HEX转HSB
      if (typeof data.other_money === "number") setRestValuePrice(String(data.other_money));

      if (typeof data.all_price === "number") setAllPriceValue(String(data.all_price));
      if (typeof data.all_tiers === "number") setAllTiersValue(String(data.all_tiers));

      if (typeof data.show_cart === "number") setInsuranceVisibility(String(data.show_cart));
      if (typeof data.show_cart_icon === "number") setIconVisibility(String(data.show_cart_icon));
      if (typeof data.select_button === "number") setSelectButton(String(data.select_button));

      if (data.product_type) setProductTypeInput(data.product_type);

      if (Array.isArray(data.product_collection)) setSelectedCollections(data.product_collection);

      if (typeof data.pricing_type === "number") setPricingType(String(data.pricing_type));

      if (typeof data.price_rule === "number") setPricingRule(String(data.pricing_rule));

      if (Array.isArray(data.price_select)) setPriceSelect(data.price_select);

      if (Array.isArray(data.tiers_select)) setTiersSelect(data.tiers_select);

      if (Array.isArray(data.icons) && data.icons.length > 0) setIcons(data.icons);

    } catch (error) {
      console.error("Error fetching cart data:", error);
    } finally {
      setIsLoading(false);
    }

  };

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
      optInColor: hsbToHex(optInColor),
      optOutColor: hsbToHex(optOutColor),
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

    try {

      const res = await rqPostUpdateCartSetting(payload);

      setToastContent(
        res?.code === 200
          ? "Saved successfully"
          : "Saved Fail"
      );

      setToastError(res?.code !== 200);
      setToastActive(true);

    } catch (error) {
      console.error("Error saving settings:", error);
      setToastContent("Service Error");
      setToastError(true);
      setToastActive(true);
    } finally {

    }
    setSaveLoading(false);
    handleDiscard();
  }, [
    planTitle,
    iconVisibility,
    insuranceVisibility,
    selectButton,
    addonTitle,
    enabledDescription,
    disabledDescription,
    footerText,
    footerUrl,
    optInColor,
    optOutColor,
    pricingType,
    pricingRule,
    priceSelect,
    tiersSelect,
    restValuePrice,
    productTypeInput,
    selectedCollections,
    errors,
  ]);

  const handleTierChange = (index, field, value, error) => {
    let validValue = value;

    // 如果值为空字符串，允许删除（空值也有效）
    if (validValue === "") {
      const updatedTiers = [...tiersSelect];
      updatedTiers[index][field] = "";
      setTiersSelect(updatedTiers);
    } else {
      // 正则表达式验证输入是数字并且不超过100
      validValue = validValue.replace(/[^0-9.]/g, ""); // 只保留数字和小数点
      const validNumber = /^[0-9]+(\.[0-9]{0,2})?$/; // 验证有效的数字格式（最多两位小数）

      if (validNumber.test(validValue) && (parseFloat(validValue) <= 100)) {
        const updatedTiers = [...tiersSelect];
        updatedTiers[index][field] = validValue; // 更新输入值
        setTiersSelect(updatedTiers);
      }
    }
    handleOnDiscard();
    if (errors[error]) {
      setErrors(prev => ({...prev, [error]: undefined}));
    }
  };

  const handleDeleteTier = (index) => {
    const updatedTiers = tiersSelect.filter((_, i) => i !== index); // 根据索引删除
    setTiersSelect(updatedTiers); // 更新状态
    handleOnDiscard();
  };

  const handleDeletePrice = (index) => {
    const updatedPrices = priceSelect.filter((_, i) => i !== index); // 根据索引删除
    setPriceSelect(updatedPrices); // 更新状态
    handleOnDiscard();
  };

  const handlePriceChange = (index, field, value, error) => {
    let validValue = value;

    // 如果值为空字符串，允许删除
    if (validValue === "") {
      const updatedPrices = [...priceSelect];
      updatedPrices[index][field] = "";
      setPriceSelect(updatedPrices);
    } else {
      // 正则表达式验证输入是数字并且不超过100
      validValue = validValue.replace(/[^0-9.]/g, ""); // 只保留数字和小数点
      const validNumber = /^[0-9]+(\.[0-9]{0,2})?$/; // 验证有效的数字格式（最多两位小数）

      if (validNumber.test(validValue)) {
        const updatedPrices = [...priceSelect];
        updatedPrices[index][field] = validValue; // 更新输入值
        setPriceSelect(updatedPrices);
      }
    }
    handleOnDiscard();
    if (errors[error]) {
      setErrors(prev => ({...prev, [error]: undefined}));
    }
  };

  // 切换勾选状态
  const handleCheckbox = () => {
    handleOnDiscard();
    setCheckboxInput(prev => !prev);
  };

// 切换 switch 状态
  const onChangeSwitch = (value) => {
    setSwitchValue(value);
  };

  const onInsuranceVisibility = (value) => {
    handleOnDiscard();
    if (value) {
      setInsuranceVisibility("1");
    } else {
      setInsuranceVisibility("0");
    }
  };

  // 右侧Demo的开关/复选框渲染
  const renderSelectionControl = () => {
    if (selectButton === "0") {
      return (
        <s-switch
          onChange={(nextChecked) => onChangeSwitch(nextChecked)}
          checked={switchValue}
          onColor={hsbToHex(optInColor)}
          offColor={hsbToHex(optOutColor)}
          uncheckedIcon={<div className="switchBtn" />}
          checkedIcon={false} // 如果你也想隐藏选中时的图标
        />
      );
    } else {

      return (
        <label className="custom-checkbox">
          <input
            type="checkbox"
            checked={checkboxInput}
            onChange={handleCheckbox}
          />
          <span className="checkmark" style={{
            backgroundColor: checkboxInput ? hsbToHex(optInColor) : hsbToHex(optOutColor),
            borderColor: checkboxInput ? hsbToHex(optInColor) : hsbToHex(optOutColor),
          }} />
        </label>
      );
    }
  };

  if (isLoading) {
    return <SkeletonScreen />;
  }

  const handleDiscard = () => {
    setDirty(false);
  };

  const handleOnDiscard = () => {
    if (!dirty) setDirty(true); // 统一标记
  };

  return (
    <div style={{position: "relative"}}>
      {/* 当有未保存更改时显示 ContextualSaveBar */}
      {dirty && (
        <ContextualSaveBar
          message="Unsaved changes"
          saveAction={{
            onAction: handleSave,
            loading: saveLoading,
            disabled: false,
          }}
          discardAction={{
            onAction: handleDiscard,
          }}
        />
      )}

      <Page title="Creat Protection Plan (Cart Page)"
            primaryAction={{content: "Save Changes", onAction: handleSave, loading: saveLoading}}>
        <Layout>
          <Layout.Section>
            <InlineStack align="start" gap="600" wrap={false}>
              <s-box width="480px">
                <BlockStack gap="400">
                  <Card padding="300">  {/* 减少内边距 */}
                    <InlineStack align="space-between" blockAlign="center" gap="200">
                      <Text variant="headingSm" as="h6">Turn on/off insurance</Text>
                      <s-box paddingLeft="200">  {/* 给开关添加左内边距 */}
                        <s-switch
                          onChange={(nextChecked) => onInsuranceVisibility(nextChecked)}
                          checked={insuranceVisibility === "1"}
                          onColor="#303030"
                          offColor="#DADADA"
                          uncheckedIcon={<div className="switchBtn" />}
                        />
                      </s-box>
                    </InlineStack>
                  </Card>


                  <Card padding="400">
                    <BlockStack gap="300">
                      <Text variant="headingSm" as="h6">Widget Style</Text>

                      <BlockStack gap="200">
                        <Text>Icon Visibility</Text>
                        <Select
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
                        <Text as="h2" variant="bodyMd">Widget Icon</Text>
                        <InlineStack id="widget-icons" wrap={false} gap="200">
                          {icons.map((icon, index) => {
                            const isSelected = icon.selected;
                            return (
                              <div
                                key={index}
                                onClick={() => handleIconClick(icon.id)}
                                style={{
                                  position: "relative",
                                  padding: "8px",
                                  borderRadius: "12px",
                                  border: isSelected ? "2px solid #5c6ac4" : "1px solid #dfe3e8",
                                  backgroundColor: isSelected ? "#f0f1f3" : "white",
                                  cursor: "pointer",
                                  //transition: 'all 0.1s ease',
                                }}
                                onMouseEnter={(e) => {
                                  if (!isSelected) e.currentTarget.style.borderColor = "#b4bcc4";
                                }}
                                onMouseLeave={(e) => {
                                  if (!isSelected) e.currentTarget.style.borderColor = "#dfe3e8";
                                }}
                              >
                                <Thumbnail
                                  source={icon.src}
                                  alt={icon.label}
                                  size="medium"
                                />
                                {isSelected && (
                                  <div style={{
                                    position: "absolute",
                                    top: "6px",
                                    right: "6px",
                                    backgroundColor: "#5c6ac4",
                                    color: "white",
                                    borderRadius: "50%",
                                    width: "18px",
                                    height: "18px",
                                    display: "flex",
                                    alignItems: "center",
                                    justifyContent: "center",
                                    fontSize: "12px",
                                    fontWeight: "bold",
                                  }}>
                                    ✓
                                  </div>
                                )}
                              </div>
                            );
                          })}
                        </InlineStack>
                      </BlockStack>

                      <BlockStack gap="200">
                        <Text>Select Button</Text>
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
                        />
                      </BlockStack>

                      <InlineStack align="space-between" gap="300">
                        <s-box>
                          <div>Opt-in action button</div>
                          <Button style={{width: "130px"}}
                                  onClick={() => setEditingColor(editingColor === "opt-in" ? null : "opt-in")}>
                            <InlineStack gap="100" align="center">
                              <div
                                style={{
                                  width: "16px",
                                  height: "16px",
                                  borderRadius: "4px",
                                  backgroundColor: hsbToHex(optInColor),
                                  border: "1px solid #ccc",
                                }}
                              />
                              <span>{hsbToHex(optInColor)}</span>
                            </InlineStack>
                          </Button>
                        </s-box>
                        <s-box>
                          <div>Opt-out action button</div>
                          <Button style={{width: "130px"}}
                                  onClick={() => setEditingColor(editingColor === "opt-out" ? null : "opt-out")}>
                            <InlineStack gap="100" align="center">
                              <div
                                style={{
                                  width: "16px",
                                  height: "16px",
                                  borderRadius: "4px",
                                  backgroundColor: hsbToHex(optOutColor),
                                  border: "1px solid #ccc",
                                }}
                              />
                              <span>{hsbToHex(optOutColor)}</span>
                            </InlineStack>
                          </Button>
                        </s-box>
                      </InlineStack>

                      {editingColor && (
                        <div ref={colorPickerRef}>
                          <ColorPicker
                            color={editingColor === "opt-in" ? optInColor : optOutColor}
                            onChange={(color) => {
                              if (editingColor === "opt-in") {
                                setOptInColor(color);
                              } else {
                                setOptOutColor(color);
                              }
                            }}
                          />
                        </div>
                      )}

                    </BlockStack>
                  </Card>

                  <Card padding="400">
                    <BlockStack gap="200">
                      <Text variant="headingSm" as="h6">Content</Text>
                      <TextField
                        label="Add-on title"
                        value={addonTitle}
                        onChange={handleFieldChange(setAddonTitle, "addonTitle")}
                        error={errors.addonTitle}
                        maxLength={50}
                      />
                      <TextField
                        label="Enabled description"
                        value={enabledDescription}
                        onChange={handleFieldChange(setEnabledDescription, "enabledDescription")}
                        multiline={4}
                        maxLength={200}
                        error={errors.enabledDescription}
                      />
                      <TextField
                        label="Disabled description"
                        value={disabledDescription}
                        onChange={handleFieldChange(setDisabledDescription, "disabledDescription")}
                        multiline={4}
                        maxLength={200}
                        error={errors.disabledDescription}
                      />
                      <TextField
                        label="Footer link text"
                        value={footerText}
                        onChange={setFooterText}
                        maxLength={50}
                      />
                      <TextField
                        label="Footer link URL"
                        value={footerUrl}
                        onChange={setFooterUrl}
                        maxLength={150}
                      />
                      <Text tone="subdued" variant="bodySm">Note: leave blank for no link</Text>
                    </BlockStack>
                  </Card>

                  {/* New Protection Pricing Card */}
                  <Card padding="400">
                    <BlockStack gap="400">
                      <Text variant="headingSm" as="h6">Protection Pricing</Text>
                      <BlockStack gap="200">
                        <Text>Pricing type</Text>
                        <Select
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
                        <Text>Pricing rule</Text>
                        <Select
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

                      {pricingRule === "1" ?
                        pricingType === "1" ? (
                          <>
                            <DataTable
                              columnContentTypes={["text", "text", "text", "text"]}
                              headings={["Min Cart Value", "Max Cart Value", "Protection Price", ""]}
                              rows={tiersSelect.map((tier, index) => [
                                <TextField
                                  key={`min-${index}`}
                                  value={tier.min}
                                  onChange={(value) => handleTierChange(index, "min", value, `tier_min_${index}`)}
                                  prefix={moneySymbol}
                                  error={errors[`tier_min_${index}`]}
                                />,
                                <TextField
                                  key={`max-${index}`}
                                  value={tier.max}
                                  onChange={(value) => handleTierChange(index, "max", value, `tier_max_${index}`)}
                                  prefix={moneySymbol}
                                  error={errors[`tier_max_${index}`]}
                                />,
                                <TextField
                                  key={`percentage-${index}`}
                                  value={tier.percentage}
                                  onChange={(value) => handleTierChange(index, "percentage", value, `tier_percentage_${index}`)}
                                  prefix="%"
                                  error={errors[`tier_percentage_${index}`]}
                                />,
                                <span key={`delete-${index}`} style={{cursor: "pointer"}}
                                      onClick={() => handleDeleteTier(index)}>X</span>, // 删除按钮
                              ])}
                            />
                            <Button onClick={handleAddTier}>Add Tier</Button>
                          </>
                        ) : (
                          <>
                            <DataTable
                              columnContentTypes={["text", "text", "text", "text"]}
                              headings={["Min Cart Value", "Max Cart Value", "Protection Price", ""]}
                              rows={priceSelect.map((tier, index) => [
                                <TextField
                                  key={`min-${index}`}
                                  value={tier.min}
                                  onChange={(value) => handlePriceChange(index, "min", value, `price_min_${index}`)}
                                  prefix={moneySymbol}
                                  error={errors[`price_min_${index}`]}
                                />,
                                <TextField
                                  key={`max-${index}`}
                                  value={tier.max}
                                  onChange={(value) => handlePriceChange(index, "max", value, `price_max_${index}`)}
                                  prefix={moneySymbol}
                                  error={errors[`price_max_${index}`]}
                                />,
                                <TextField
                                  key={`price-${index}`}
                                  value={tier.price}
                                  onChange={(value) => handlePriceChange(index, "price", value, `price_price_${index}`)}
                                  prefix={moneySymbol}
                                  error={errors[`price_price_${index}`]}
                                />,
                                <span key={`delete-${index}`} style={{cursor: "pointer"}}
                                      onClick={() => handleDeletePrice(index)}>X</span>, // 删除按钮
                              ])}
                            />
                            <Button onClick={handleAddPrice}>Add Price</Button>
                          </>
                        )
                        : <></>
                      }
                      {pricingRule === "1" ?
                        <TextField
                          label="Other Value Range"
                          value={restValuePrice}
                          onChange={(value) => {
                            handleOnDiscard();
                            setRestValuePrice(value);
                          }}
                          prefix="$"
                        />
                        :
                        <TextField
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
                      }


                      <Divider />

                    </BlockStack>
                  </Card>

                  <Card padding="400">
                    <BlockStack gap="200">

                      <TextField
                        label="Product Types"
                        value={productTypeInput}
                        onChange={setProductTypeInput}
                        maxLength={50}
                      />
                      {/* Collections */}
                      <Text variant="bodySm">Collections</Text>
                      <Autocomplete
                        options={collectionOptions.filter(opt => !selectedCollections.includes(opt.value))}
                        selected={[]}
                        onSelect={(selectedValue) => {
                          const value = selectedValue[0];
                          setSelectedCollections([...new Set([...selectedCollections, value])]);
                          setCollectionInput("");
                        }}
                        textField={
                          <Autocomplete.TextField
                            label=""
                            value={collectionInput}
                            onChange={setCollectionInput}
                            onBlur={() => setCollectionInput("")}
                            onKeyDown={(e) => {
                              if (e.key === "Enter") {
                                e.preventDefault();
                                setCollectionInput("");
                              }
                            }}
                            placeholder="Choose from list"
                          />
                        }
                      />
                      <InlineStack gap="100" wrap>
                        {selectedCollections.map((value) => {
                          const label = collectionOptions.find((o) => o.value === value)?.label || value;
                          return (
                            <Tag key={value}
                                 onRemove={() => removeTag(value, selectedCollections, setSelectedCollections)}>
                              {label}
                            </Tag>
                          );
                        })}
                      </InlineStack>
                    </BlockStack>
                  </Card>

                </BlockStack>
              </s-box>

              {/* Right Side - Demo Preview */}
              <s-box style={{position: "sticky", top: "0px", alignSelf: "flex-start", width: "480px"}}>
                {/* 主体内容 */}
                <Card padding="0">
                  <BlockStack
                    gap="400"
                    style={{
                      backgroundColor: "rgb(205 203 205)",  // 灰色背景
                      padding: "16px",
                      borderRadius: "4px 4px 0 0", // 仅圆角顶部，与 Card 对齐
                      marginBottom: "0",           // 移除底部 margin
                      width: "100%",               // 占满 Card 宽度
                      justifyContent: "space-between"
                    }}
                  >
                    <Text as="h6" variant="bodyMd" fontWeight="semibold">
                      Cart Page Demo
                    </Text>

                    <Text as="h6" variant="bodyMd">
                      View in store
                    </Text>
                  </BlockStack>

                  <s-box padding="400"> {/* 将 Card 的 padding 移到内部 s-box */}
                    <BlockStack gap="300">
                      {/* Mock Products */}
                      {[1, 2].map((item, idx) => (
                        <InlineStack gap="300" align="start" key={idx}>
                          <Thumbnail
                            source="https://img.icons8.com/plasticine/100/cat-footprint.png"
                            alt="Cat Slippers"
                            size="medium"
                          />
                          <BlockStack>
                            <Text variant="bodyMd" fontWeight="medium">Cute Cat Slippers</Text>
                            <Text>$10.00</Text>
                          </BlockStack>
                        </InlineStack>
                      ))}

                      {/* 保险勾选项 */}
                      <Card
                        padding="300"
                        style={{
                          position: "relative",
                          borderLeft: `4px solid ${hsbToHex(optInColor)}`,
                          minHeight: "120px",
                        }}
                      >
                        {/* 勾选框固定在右上角 */}
                        <s-box >

                          {renderSelectionControl()}
                        </s-box>

                        {/* 主体内容 */}
                        <s-box style={{display: "flex", gap: "16px", alignItems: "flex-start"}}>
                          {iconVisibility === "1" ? <img
                            src={selectedIcon.src}

                            alt="Protection"
                            style={{flexShrink: 0, width: "50px", height: "50px"}}
                          /> : <div />}


                          <s-box>
                            <s-box style={{maxWidth: "280px", wordBreak: "break-word"}}>
                              <Text
                                tone="success"
                                fontWeight="medium"
                                style={{color: hsbToHex(optInColor)}}
                              >
                                {addonTitle || "Shipping Protection"} (2.00 USD)
                              </Text>
                            </s-box>
                            <s-box style={{maxWidth: "300px", wordBreak: "break-word"}}>
                              <Text tone="subdued" variant="bodySm">
                                {/*{switchValue || checkboxInput === true ? "$22.00" : "$20.00"}*/}
                                {(selectButton === "0" && switchValue === true) || (selectButton === "1" && checkboxInput === true) ? enabledDescription : disabledDescription}

                              </Text>
                            </s-box>
                          </s-box>
                        </s-box>

                        {/* 底部链接 */}
                        {footerText && (
                          <s-box style={iconVisibility === "1" ? {marginTop: "5px", marginLeft: "65px"} : {
                            marginTop: "5px",
                            marginLeft: "15px"
                          }}>
                            <Text variant="bodySm" as="span" >
                              <a title={footerUrl || ""} style={{
                                color: "#0070f3",        /* 设置字体颜色为蓝色 */
                                textDecoration: "none", cursor: "pointer"
                              }} onClick={(e) => e.preventDefault()}>
                                {footerText}
                              </a>
                            </Text>
                          </s-box>
                        )}
                      </Card>

                      <Divider />

                      {/* Checkout 按钮 */}
                      <Button fullWidth
                              size="large">Checkout {switchValue || checkboxInput === true ? "22.00 USD" : "20.00 USD"}</Button>
                    </BlockStack>
                  </s-box>
                </Card>
              </s-box>
            </InlineStack>
          </Layout.Section>
        </Layout>
        {toastActive && (
          <Toast
            content={Object.keys(errors).length ? "Please fix validation errors" : toastContent}
            onDismiss={() => setToastActive(false)}
            error={Object.keys(errors).length > 0 || toastError}
          />
        )}
      </Page>

    </div>

  );
}