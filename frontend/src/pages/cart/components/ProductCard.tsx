import { ProductCardProps } from "@/types/cart.ts";
import { BlockStack, Card, TextField } from "@shopify/polaris";
import CollectionSelector from "@/pages/cart/components/CollectionSelector.tsx";

export default function ProductCard({
  productSettings,
  collectionOptions,
  onProductTypeChange,
  onCollectionInputChange,
  onCollectionSelect,
  onRemoveCollection,
}: ProductCardProps) {
  return (
    <Card padding="400">
      <BlockStack gap="200">
        <TextField
          autoComplete="off"
          label="Product Types"
          value={productSettings.productTypeInput}
          onChange={onProductTypeChange}
          maxLength={50}
        />

        <CollectionSelector
          collectionOptions={collectionOptions}
          selectedCollections={productSettings.selectedCollections}
          collectionInput={productSettings.collectionInput}
          onCollectionInputChange={onCollectionInputChange}
          onCollectionSelect={onCollectionSelect}
          onRemoveCollection={onRemoveCollection}
        />
      </BlockStack>
    </Card>
  );
}
