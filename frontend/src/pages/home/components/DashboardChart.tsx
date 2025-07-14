import React from "react";
import {Collapsible} from "@shopify/polaris";
import ReactECharts from "echarts-for-react";
import * as echarts from "echarts/core";
import "echarts/lib/chart/lines";
import "echarts/lib/component/tooltip";
import "echarts/lib/component/title";
import {useChartOptions} from "@/hooks/home/useChartOptions.ts";

interface OrderStatisticsTableItem {
  date: string;
  sales: number;
  refund: number;
}

interface DashboardChartProps {
  data: OrderStatisticsTableItem[];
  isCollapsed: boolean;
}

const DashboardChart = ({
                          data,
                          isCollapsed,
                        }: DashboardChartProps) => {
  const chartOptions = useChartOptions(data);

  return <Collapsible
    open={isCollapsed}
    id="dashboard-chart-collapsible"
    transition={{
      duration: "500ms",
      timingFunction: "ease-in-out",
    }}
    expandOnPrint
  >
    <div className="statistical_order_line">
      <ReactECharts
        style={{height: "400px"}}
        echarts={echarts}
        option={chartOptions}
        notMerge
        lazyUpdate
        theme="theme_name"
        onChartReady={() => {
        }}
        opts={{renderer: "canvas"}}
      />
    </div>
  </Collapsible>;

};

export default DashboardChart;
