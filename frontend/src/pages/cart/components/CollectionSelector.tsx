import React from "react";
import { Autocomplete, InlineStack, Tag, Text } from "@shopify/polaris";
import type { OptionDescriptor } from "@shopify/polaris/build/ts/src/types";

interface CollectionSelectorProps {
  collectionOptions: OptionDescriptor[];
  selectedCollections: string[];
  collectionInput: string;
  onCollectionInputChange: (value: string) => void;
  onCollectionSelect: (value: string) => void;
  onRemoveCollection: (value: string) => void;
}

const CollectionSelector: React.FC<CollectionSelectorProps> = ({
  collectionOptions,
  selectedCollections,
  collectionInput,
  onCollectionInputChange,
  onCollectionSelect,
  onRemoveCollection,
}) => {
  return (
    <>
      <Text as="p" variant="bodySm">Collections</Text>
      <Autocomplete
        options={collectionOptions.filter(opt => !selectedCollections.includes(opt.value))}
        selected={[]}
        onSelect={(selectedValue) => {
          const value = selectedValue[0];
          onCollectionSelect(value);
        }}
        textField={
          <Autocomplete.TextField
            autoComplete="off"
            label=""
            value={collectionInput}
            onChange={onCollectionInputChange}
            onBlur={() => onCollectionInputChange("")}
            placeholder="Choose from list"
          />
        }
      />
      <InlineStack gap="100" wrap>
        {selectedCollections.map((value) => {
          const label = collectionOptions.find((o) => o.value === value)?.label || value;
          return (
            <Tag key={value} onRemove={() => onRemoveCollection(value)}>
              {label}
            </Tag>
          );
        })}
      </InlineStack>
    </>
  );
};

export default CollectionSelector;
