import { useState, useCallback } from 'react';
import { 
  Button, 
  Popover, 
  OptionList,
  Card,
  Text,
  Divider
} from '@shopify/polaris';
import { SortIcon } from '@shopify/polaris-icons';
import { 
  SortOption,
  SortDirection 
} from '@/types/order';
import {SORT_DIRECTION, SORT_OPTIONS} from "@/constants/orderFilters.ts";

interface SortSelectorProps {
  sortBy?: SortOption;
  sortDirection?: SortDirection;
  onSortChange: (sortBy?: SortOption, sortDirection?: SortDirection) => void;
}

export const SortSelector = ({
  sortBy,
  sortDirection,
  onSortChange,
}: SortSelectorProps) => {
  const [popoverActive, setPopoverActive] = useState(false);

  const sortOptions = [
    { value: '', label: 'Default' },
    { value: SORT_OPTIONS.ORDER_NUMBER, label: 'Order Number' },
    { value: SORT_OPTIONS.PAYMENT_DATE, label: 'Payment Date' },
    { value: SORT_OPTIONS.PROTECTION_FEE, label: 'Protection Fee' },
  ];

  const directionOptions = [
    { value: SORT_DIRECTION.ASC, label: 'Ascending' },
    { value: SORT_DIRECTION.DESC, label: 'Descending' },
  ];

  const handleSortBySelect = useCallback((values: string[]) => {
    const value = values[0];
    if (value) {
      onSortChange(value as SortOption, sortDirection || SORT_DIRECTION.ASC);
    } else {
      onSortChange(undefined, undefined);
    }
  }, [onSortChange, sortDirection]);

  const handleDirectionSelect = useCallback((values: string[]) => {
    const value = values[0] as SortDirection;
    if (sortBy) {
      onSortChange(sortBy, value);
    }
  }, [onSortChange, sortBy]);

  const handleClearSort = useCallback(() => {
    onSortChange(undefined, undefined);
  }, [onSortChange]);

  const getSortLabel = () => {
    if (sortBy) {
      const sortOption = sortOptions.find(opt => opt.value === sortBy);
      const direction = sortDirection === SORT_DIRECTION.DESC ? 'Desc' : 'Asc';
      return `${sortOption?.label} (${direction})`;
    }
    return 'Sort';
  };

  const activator = (
    <Button
      onClick={() => setPopoverActive(!popoverActive)}
      disclosure
      icon={SortIcon}
      tone={sortBy ? 'success' : undefined}
    >
      {getSortLabel()}
    </Button>
  );

  return (
    <Popover
      active={popoverActive}
      activator={activator}
      onClose={() => setPopoverActive(false)}
      preferredAlignment="left"
    >
      <Card>
        <div style={{ padding: '1rem', minWidth: '200px' }}>
          <Text variant="headingMd" as="h3">
            Sort By
          </Text>
          <div style={{ marginTop: '0.5rem' }}>
            <OptionList
              options={sortOptions}
              selected={sortBy ? [sortBy] : ['']}
              onChange={handleSortBySelect}
            />
          </div>
          
          {sortBy && (
            <>
              <Divider />
              <div style={{ marginTop: '1rem' }}>
                <Text variant="headingMd" as="h3">
                  Direction
                </Text>
                <div style={{ marginTop: '0.5rem' }}>
                  <OptionList
                    options={directionOptions}
                    selected={sortDirection ? [sortDirection] : [SORT_DIRECTION.ASC]}
                    onChange={handleDirectionSelect}
                  />
                </div>
              </div>
            </>
          )}
          
          {sortBy && (
            <>
              <Divider />
              <div style={{ marginTop: '1rem' }}>
                <Button onClick={handleClearSort} tone="critical">
                  Clear Sort
                </Button>
              </div>
            </>
          )}
        </div>
      </Card>
    </Popover>
  );
};
