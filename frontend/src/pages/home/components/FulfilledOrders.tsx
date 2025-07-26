import React, { useState, useCallback } from 'react';
import {BlockStack, Box, Button, Card, InlineStack, Spinner} from "@shopify/polaris";
import intl from '@/lib/i18n.ts';
import { useDashboardData } from '@/hooks/home/useDashboardData';
import PeriodSelector from './PeriodSelector.tsx';
import StatisticItem from './StatisticItem.tsx';
import DashboardChart from './DashboardChart.tsx';
import {ChevronDownIcon, ChevronUpIcon} from "@shopify/polaris-icons";

interface FulfilledOrdersProps {
  // 如果有props，可以在这里定义
}

const FulfilledOrders: React.FC<FulfilledOrdersProps> = () => {
  const { data, loading, error, fetchData } = useDashboardData();
  const [selectedPeriod, setSelectedPeriod] = useState('30');
  const [currentDateLabel, setCurrentDateLabel] = useState('Last 30 days');
  const [isChartCollapsed, setIsChartCollapsed] = useState(true);
  const [popoverActive, setPopoverActive] = useState(false);

  const periodLabels = {
    '30': 'Last 30 days',
    '90': 'Last 90 days',
    '365': 'Last 365 days',
  };

  const handlePeriodChange = useCallback(async (period: string) => {
    setSelectedPeriod(period);
    setCurrentDateLabel(periodLabels[period as keyof typeof periodLabels]);
    setPopoverActive(false);
    await fetchData(period);
  }, [fetchData]);

  const handleTogglePopover = useCallback(() => {
    setPopoverActive(!popoverActive);
  }, [popoverActive]);

  const handleToggleChart = useCallback(() => {
    setIsChartCollapsed(!isChartCollapsed);
  }, [isChartCollapsed]);

  // 加载状态
  if (loading) {
    return (
      <div className="pt-3">
        <BlockStack gap="300">
          <PeriodSelector
            selectedPeriod={selectedPeriod}
            currentDateLabel={currentDateLabel}
            onPeriodChange={handlePeriodChange}
            popoverActive={popoverActive}
            onTogglePopover={handleTogglePopover}
          />
        <Card>
          <div className="flex justify-center items-center " style={{ padding: '2rem', textAlign: 'center',minHeight:"500px" }}>
            <Spinner size="large" />
          </div>
        </Card>
        </BlockStack>
      </div>
    );
  }

  // 错误状态
  if (error) {
    return (
      <div className="statistical_order">
        <Card>
          <div style={{ padding: '2rem', textAlign: 'center', color: 'red' }}>
            Error: {error}
          </div>
        </Card>
      </div>
    );
  }

  const { orderStatics } = data;

  return (
    <div className="pt-3">
        <BlockStack gap="300">
          <PeriodSelector
            selectedPeriod={selectedPeriod}
            currentDateLabel={currentDateLabel}
            onPeriodChange={handlePeriodChange}
            popoverActive={popoverActive}
            onTogglePopover={handleTogglePopover}
          />

          <Card>
            <BlockStack>
              <InlineStack gap="200" wrap={false} blockAlign="center">
                <StatisticItem
                  title={intl.get('Pretection Orders') as string}
                  value={orderStatics?.orders || 0}
                />

                <StatisticItem
                  title={intl.get('Insured Sales') as string}
                  value={orderStatics?.sales || 0}
                  prefix="$ "
                />

                <StatisticItem
                  title={intl.get('Refund Amount') as string}
                  value={orderStatics?.refund || 0}
                  prefix="$ "
                />

                <StatisticItem
                  title={intl.get('Total Revenue') as string}
                  value={orderStatics?.total || 0}
                  prefix="$ "
                  color={
                    (orderStatics?.total || 0) >= 0 ? '#303030' : '#f00'
                  }
                  tooltipContent={
                    <Box padding="200">
                      <Box>
                        <strong>Total Revenue</strong>
                      </Box>
                      <Box />
                    </Box>
                  }
                />
                <div>
                  <Button
                    variant="tertiary"
                    icon={isChartCollapsed ? ChevronUpIcon : ChevronDownIcon}
                    onClick={handleToggleChart}
                  />
                </div>
              </InlineStack>
              <DashboardChart
                data={data.orderStaticsTable}
                isCollapsed={isChartCollapsed}
              />
            </BlockStack>
          </Card>
        </BlockStack>
    </div>
  );
};

export default FulfilledOrders;