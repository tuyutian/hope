import { useState, useCallback } from 'react';
import {
  ChoiceList,
  TextField,
  ButtonGroup,
  Button,
  BlockStack,
  Text,
  Divider
} from '@shopify/polaris';
import type { TimeRange } from '@/types/order';
import {TIME_RANGES} from "@/constants/orderFilters.ts";

interface TimeRangeFilterProps {
  selectedRange: TimeRange;
  customStartDate?: string;
  customEndDate?: string;
  onRangeChange: (range: TimeRange, startDate?: string, endDate?: string) => void;
}

export const TimeRangeFilter = ({
                                  selectedRange,
                                  customStartDate,
                                  customEndDate,
                                  onRangeChange,
                                }: TimeRangeFilterProps) => {
  const [showCustomPicker, setShowCustomPicker] = useState(selectedRange === TIME_RANGES.CUSTOM);
  const [tempStartDate, setTempStartDate] = useState(customStartDate || '');
  const [tempEndDate, setTempEndDate] = useState(customEndDate || '');

  const timeRangeChoices = [
    { label: 'Today', value: TIME_RANGES.TODAY },
    { label: 'Last 7 days', value: TIME_RANGES.LAST_7_DAYS },
    { label: 'Last 14 days', value: TIME_RANGES.LAST_14_DAYS },
    { label: 'Last 30 days', value: TIME_RANGES.LAST_30_DAYS },
    { label: 'Last 60 days', value: TIME_RANGES.LAST_60_DAYS },
    { label: 'Last 90 days', value: TIME_RANGES.LAST_90_DAYS },
    { label: 'Last 180 days', value: TIME_RANGES.LAST_180_DAYS },
    { label: 'Last 360 days', value: TIME_RANGES.LAST_360_DAYS },
    { label: 'Custom range', value: TIME_RANGES.CUSTOM },
  ];

  const handleRangeSelect = useCallback((value: string[]) => {
    const range = value[0] as TimeRange;
    if (range === TIME_RANGES.CUSTOM) {
      setShowCustomPicker(true);
    } else {
      setShowCustomPicker(false);
      onRangeChange(range);
    }
  }, [onRangeChange]);

  const handleCustomRangeApply = useCallback(() => {
    if (tempStartDate && tempEndDate) {
      onRangeChange(TIME_RANGES.CUSTOM, tempStartDate, tempEndDate);
    }
  }, [tempStartDate, tempEndDate, onRangeChange]);

  const handleCustomRangeCancel = useCallback(() => {
    setShowCustomPicker(false);
    setTempStartDate(customStartDate || '');
    setTempEndDate(customEndDate || '');
    // 如果当前是 custom 状态，重置为默认值
    if (selectedRange === TIME_RANGES.CUSTOM) {
      onRangeChange(TIME_RANGES.LAST_30_DAYS);
    }
  }, [customStartDate, customEndDate, selectedRange, onRangeChange]);

  return (
    <BlockStack gap="400">
      <ChoiceList
        title="Time Range"
        choices={timeRangeChoices}
        selected={[selectedRange]}
        onChange={handleRangeSelect}
      />

      {showCustomPicker && (
        <BlockStack gap="300">
          <Divider />
          <Text variant="headingMd" as="h4">
            Custom Date Range
          </Text>

          <TextField
            label="Start Date"
            type="date"
            value={tempStartDate}
            onChange={setTempStartDate}
            autoComplete="off"
          />

          <TextField
            label="End Date"
            type="date"
            value={tempEndDate}
            onChange={setTempEndDate}
            autoComplete="off"
          />

          <ButtonGroup>
            <Button onClick={handleCustomRangeCancel}>
              Cancel
            </Button>
            <Button
              variant="primary"
              onClick={handleCustomRangeApply}
              disabled={!tempStartDate || !tempEndDate}
            >
              Apply
            </Button>
          </ButtonGroup>
        </BlockStack>
      )}
    </BlockStack>
  );
};