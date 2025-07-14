import { useMemo } from 'react';

interface ChartOption {
  title: Record<string, any>;
  tooltip: {
    trigger: string;
  };
  color: string[];
  legend: {
    bottom: string;
    itemWidth: number;
    itemHeight: number;
    data: {
      name: string;
      icon: string;
    }[];
  };
  grid: {
    left: string;
    right: string;
    bottom: string;
    containLabel: boolean;
  };
  toolbox: Record<string, any>;
  xAxis: {
    type: string;
    boundaryGap: boolean;
    axisTick: {
      show: boolean;
    };
    data: string[];
  };
  yAxis: {
    offset: number;
  };
  series: {
    name: string;
    type: string;
    smooth: boolean;
    lineStyle?: {
      color: string;
      width: number;
      type: string;
    };
    data: number[];
  }[];
}

interface OrderStatisticsTableItem {
  date: string;
  sales: number;
  refund: number;
}

export const useChartOptions = (data: OrderStatisticsTableItem[]): ChartOption => {
  return useMemo(() => {
    const dates = data.map(item => item.date);
    const sales = data.map(item => item.sales);
    const refund = data.map(item => item.refund);

    return {
      title: {},
      tooltip: {
        trigger: "axis",
      },
      color: ["rgb(41, 171, 228)", "rgb(132, 203, 234)"],
      legend: {
        bottom: "0%",
        itemWidth: 16,
        itemHeight: 1,
        data: [
          {
            name: "Insured Sales",
            icon: "path://M5 10c0-.414.336-.75.75-.75h8.5c.414 0 .75.336.75.75s-.336.75-.75.75h-8.5c-.414 0-.75-.336-.75-.75",
          },
          {
            name: "Refund Amount",
            icon: "path://M6 10a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0Z M11.5 10a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0Z M17 10a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0Z",
          },
        ],
      },
      grid: {
        left: "3%",
        right: "4%",
        bottom: "10%",
        containLabel: true,
      },
      toolbox: {},
      xAxis: {
        type: "category",
        boundaryGap: false,
        axisTick: {
          show: false,
        },
        data: dates,
      },
      yAxis: {
        offset: 10,
      },
      series: [
        {
          name: "Insured Sales",
          type: "line",
          smooth: true,
          data: sales,
        },
        {
          name: "Refund Amount",
          type: "line",
          smooth: true,
          lineStyle: {
            color: "rgb(132, 203, 234)",
            width: 2,
            type: "dotted",
          },
          data: refund,
        },
      ],
    };
  }, [data]);
};
