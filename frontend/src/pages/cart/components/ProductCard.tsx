import { ProductCardProps } from "@/types/cart.ts";
import { BlockStack, Button, Card, Checkbox, FormLayout, InlineStack, Tag } from "@shopify/polaris";
import { useShopifyBridge } from "@/hooks/useShopifyBridge.ts";
import { useCallback } from "react";

export default function ProductCard({
  productSettings,
  onCollectionSelect,
  onlyCollection,
  onRemoveCollection,
  onCollectionChange,
}: ProductCardProps) {
  const shopify = useShopifyBridge();
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

  const handleSelect = useCallback(function (check: boolean) {
    onlyCollection(check);
  }, []);

  return (
    <Card padding="400">
      <FormLayout>
        <BlockStack gap="200">
          <InlineStack gap="200">
            <Checkbox
              onChange={handleSelect}
              checked={productSettings.onlyInCollection}
              label="By product collections"
            />
            <Button onClick={handleCollectionPicker}>Browse</Button>
          </InlineStack>
          {productSettings.selectedCollections.length > 0 && (
            <InlineStack gap="100" wrap>
              {productSettings.selectedCollections.map(item => (
                <Tag key={item.id} onRemove={() => onRemoveCollection(item.id)}>
                  {item.title}
                </Tag>
              ))}
            </InlineStack>
          )}
        </BlockStack>
      </FormLayout>
    </Card>
  );
}
