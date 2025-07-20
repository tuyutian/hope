import { ProductCardProps, ResourceItem } from "@/types/cart.ts";
import { Autocomplete, BlockStack, Button, Card, Checkbox, FormLayout, InlineStack, Tag } from "@shopify/polaris";
import { useShopifyBridge } from "@/hooks/useShopifyBridge.ts";
import { useCallback, useMemo, useState } from "react";

export default function ProductCard({
  productSettings,
  collectionOptions,
  onProductChange,
  onCollectionSelect,
  onRemoveProduct,
  onRemoveCollection,
  onCollectionChange,
}: ProductCardProps) {
  const shopify = useShopifyBridge();
  const deselectedOptions = useMemo(
    () => [
      { value: "rustic", label: "Rustic" },
      { value: "antique", label: "Antique" },
      { value: "vinyl", label: "Vinyl" },
      { value: "vintage", label: "Vintage" },
      { value: "refurbished", label: "Refurbished" },
    ],
    []
  );
  const [inputValue, setInputValue] = useState("");
  const [typeChoiceLoading, setTypeChoiceLoading] = useState(true);
  const handleCollectionPicker = useCallback(async () => {
    if (!shopify) {
      console.error("Shopify bridge not available");
      return;
    }

    const initialProductsIds = productSettings.selectedCollections.map(item => {
      return { id: `gid://shopify/Collection/${item.id}` };
    });
    try {
      // 使用新的 resourcePicker API
      const selected = await shopify.resourcePicker({
        type: "collection",
        multiple: true,
        selectionIds: initialProductsIds,
      });

      if (selected && selected.length > 0) {
        // 处理选择结果
        const selectedCollections = selected.map(item => ({
          id: Number(item.id.replace("gid://shopify/Collection/", "")),
          title: item.title,
        }));

        // 如果提供了外部处理函数，则调用它
        if (onCollectionChange) {
          onCollectionChange(selectedCollections);
        } else {
          // 否则使用默认的处理逻辑
          selectedCollections.forEach(collection => {
            onCollectionSelect(collection);
          });
        }
      }
    } catch (error) {
      console.error("Error opening resource picker:", error);
    }
  }, [productSettings.selectedCollections]);
  // 获取集合标签
  const getCollectionLabel = (collectionId: string) => {
    const option = collectionOptions.find(opt => opt.value === collectionId);
    return option?.label || collectionId;
  };

  const [options, setOptions] = useState(deselectedOptions);

  const updateText = useCallback(
    (value: string) => {
      setInputValue(value);

      if (value === "") {
        setOptions(deselectedOptions);
        return;
      }

      const filterRegex = new RegExp(value, "i");
      const resultOptions = deselectedOptions.filter(option => option.label.match(filterRegex));

      setOptions(resultOptions);
    },
    [deselectedOptions]
  );

  const verticalContentMarkup =
    productSettings.selectProductTypes.length > 0 ? (
      <InlineStack gap="100" align="center">
        {productSettings.selectProductTypes.map(option => {
          let tagLabel = "";
          tagLabel = option.title.replace("_", " ");
          return (
            <Tag key={`option${option}`} onRemove={() => onRemoveProduct(String(option.id))}>
              {tagLabel}
            </Tag>
          );
        })}
      </InlineStack>
    ) : null;
  const textField = (
    <Autocomplete.TextField
      onChange={updateText}
      label="Tags"
      value={inputValue}
      placeholder="Vintage, cotton, summer"
      verticalContent={verticalContentMarkup}
      autoComplete="off"
    />
  );

  const selectProductTypes = useMemo(function () {
    return productSettings.selectProductTypes.map(item => String(item.id));
  }, []);

  const handleSelect = useCallback(
    function (ids: string[]) {
      const selectProductTypes: ResourceItem[] = [];
      ids.map(id => {
        const item = options.find(item => item.value === String(id));
        if (item) {
          selectProductTypes.push({ title: item.label, id: Number(item.value) });
        }
      });
      onProductChange(selectProductTypes);
    },
    [onProductChange, options]
  );

  return (
    <Card padding="400">
      <FormLayout>
        <BlockStack gap="200">
          <Autocomplete
            options={options}
            selected={selectProductTypes}
            onSelect={ids => handleSelect(ids)}
            loading={typeChoiceLoading}
            textField={textField}
          />
        </BlockStack>
        <BlockStack gap="200">
          <InlineStack gap="200">
            <Checkbox label="By product collections" />
            <Button onClick={handleCollectionPicker}>Browse</Button>
          </InlineStack>
          {productSettings.selectedCollections.length > 0 && (
            <InlineStack gap="100" wrap>
              {productSettings.selectedCollections.map(item => (
                <Tag key={item.id} onRemove={() => onRemoveCollection(item.id)}>
                  {getCollectionLabel(item.title)}
                </Tag>
              ))}
            </InlineStack>
          )}
        </BlockStack>
      </FormLayout>
    </Card>
  );
}
