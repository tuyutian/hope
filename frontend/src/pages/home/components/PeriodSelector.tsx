import React from 'react';
import { Button, OptionList, Popover } from '@shopify/polaris';
import { CalendarIcon } from '@shopify/polaris-icons';
import intl from '@/lib/i18n';

interface PeriodOption {
  label: string;
  value: string;
}

interface PeriodSelectorProps {
  selectedPeriod: string;
  currentDateLabel: string;
  onPeriodChange: (period: string) => void;
  popoverActive: boolean;
  onTogglePopover: () => void;
}

const PeriodSelector: React.FC<PeriodSelectorProps> = ({
  selectedPeriod,
  currentDateLabel,
  onPeriodChange,
  popoverActive,
  onTogglePopover,
}) => {
  const options: PeriodOption[] = [
    { label: intl.get('Last 30 days') as string, value: '30' },
    { label: intl.get('Last 90 days') as string, value: '90' },
    { label: intl.get('Last 365 days') as string, value: '365' },
  ];

  const handleSelectChange = (values: string[]) => {
    if (values.length > 0) {
      onPeriodChange(values[0]);
    }
  };

  return (
    <div className="statistical_order_select">
      <Popover
        active={popoverActive}
        activator={
          <Button icon={CalendarIcon} onClick={onTogglePopover}>
            {currentDateLabel}
          </Button>
        }
        autofocusTarget="first-node"
        onClose={onTogglePopover}
      >
        <OptionList
          onChange={handleSelectChange}
          role="menuitem"
          options={options}
          selected={[selectedPeriod]}
        />
      </Popover>
    </div>
  );
};

export default PeriodSelector;
