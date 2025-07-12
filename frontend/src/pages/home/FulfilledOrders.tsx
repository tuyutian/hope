import {Box, Button, Card, Collapsible, Icon, InlineStack, OptionList, Popover,} from "@shopify/polaris";
import React, {useEffect, useState} from "react";
import {CalendarIcon, ChevronDownIcon, ChevronUpIcon,} from "@shopify/polaris-icons";
// import {
//   LineChart,
//   Line,
//   XAxis,
//   YAxis,
//   CartesianGrid,
//   Tooltip,
//   Legend,
//   ResponsiveContainer,
// } from "recharts";
import intl from "@/lib/i18n";
import Tooltips from "./Tooltips";
import ReactECharts from "echarts-for-react";
import * as echarts from "echarts/core";
import "echarts/lib/chart/lines";
import "echarts/lib/component/tooltip";
import "echarts/lib/component/title";
import {rqGetDashboard} from "@/api/index.js";

export default function FulfilledOrders(props) {
  const Data = {
    orderStaticsTable: [
      {
        date: "0",
        sales: 0,
        refund: 0,
      },
    ],
    orderStatics: {
      refund: 0,
      orders: 0,
      total: 0,
      sales: 0,
    }
  };


  const [orderData, setOrderData] = useState(Data);
  const [selected, setSelected] = useState(["30"]);
  const [currentDate, setCurrentDate] = useState("Last 30 days");
  const [getStartClo, setGetStartClo] = useState(true); //
  const options = [
    {label: intl.get("Last 30 days"), value: "30"},
    {label: intl.get("Last 90 days"), value: "90"},
    {label: intl.get("Last 365 days"), value: "365"},
  ];
  const [opction, setOpction] = useState({
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
          icon: "path://path://M6 10a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0Z M11.5 10a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0Z M17 10a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0Z",
        },
      ],
    },
    grid: {
      left: "3%",
      right: "4%",
      bottom: "10%",
      containLabel: true,
    },
    toolbox: {
      // feature: {
      //   saveAsImage: {}
      // }
    },
    xAxis: {
      type: "category",
      boundaryGap: false,
      axisTick: {
        show: false,
      },
      data: ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"],
    },

    yAxis: {
      // type: "value",
      // min:5,
      // show:true,
      // dataMin:5,
      offset: 10,
    },
    series: [
      {
        name: "Insured Sales",
        type: "line",
        smooth: true,
        data: [0, 0, 0, 0, 0, 0, 0],
      },
      {
        name: "Refund Amount",
        type: "line",
        smooth: true,
        data: [0, 0, 0, 0, 0, 0, 0],
      },
    ],
  });

  const getOrderData = async (value:string) => {
    const res = await rqGetDashboard(value);
    if (res.code === 0) {
      setOrderData((prevValue) => ({
        ...prevValue,
        orderStaticsTable: res.data.order_statistics_table,
        orderStatics: res.data.order_statistics,
      }));


      const data = res.data.order_statistics_table;
      const dates = data.map((item) => item.date);
      const sales = data.map((item) => item.sales);
      const refund = data.map((item) => item.refund);
      setOpction({
        ...opction,
        xAxis: {
          type: "category",
          boundaryGap: false,
          axisTick: {
            show: false,
          },

          data: dates,
        },
        series: [
          {
            name: "Insured Sales",
            type: "line",
            smooth: true,
            data: sales,
          },
          {
            type: "line",
            smooth: true,
            name: "Refund Amount",
            lineStyle: {
              color: "rgb(132, 203, 234)", // 线的颜色
              width: 2, // 线宽
              type: "dotted", // 线的类型，默认是'solid'，可以设置为'dashed'或'dotted'等
            },
            data: refund,
          },
        ],
      });
    }
  };

  const getDefalutData = (data) => {
    const svg = document.querySelector(".recharts-surface");
    if (svg) {
      svg.setAttribute(
        "viewBox",
        `0 0 ${svg.getAttribute("width")} ${
          Number(svg.getAttribute("height")) + 15
        }`
      );
    }
  };

  useEffect(() => {
    let isMount = false;
    if (!isMount) {
      getDefalutData();
    }
    return () => {
      isMount = true;
    };
  }, [orderData]);

  useEffect(() => {
    let isMount = false;
    if (!isMount) {
      getOrderData("30");
    }
    return () => {
      isMount = true;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handleSelectChange = (val) => {
    setSelected(val);
    getOrderData(val[0]);
    if (val[0] === "30") {
      setCurrentDate("Last 30 days");
    }
    if (val[0] === "90") {
      setCurrentDate("Last 90 days");
    }
    if (val[0] === "360") {
      setCurrentDate("Last 360 days");
    }
  };


  const [popoverActive, setPopoverActive] = useState(false);

  const togglePopoverActive = () => {
    setPopoverActive(!popoverActive);
  };

  return (
    <div className="statistical_order">
      <div className="statistical_order_top">

        <div className="statistical_order_select">
          <Popover
            active={popoverActive}
            activator={
              <Button icon={CalendarIcon} onClick={togglePopoverActive}>
                {currentDate}
              </Button>
            }
            autofocusTarget="first-node"
            onClose={togglePopoverActive}
          >
            <OptionList
              onChange={handleSelectChange}
              actionRole="menuitem"
              options={options}
              selected={selected}
            />
          </Popover>
        </div>
      </div>
      <Card>
        <div className="statistical_order_lines">
          <div className="statistical_order_line_mobile">
            <div className="statistical_order_line_item">
              <InlineStack wrap={false} direction="row">
                <div
                  className="statistical_order_line_title"
                  style={{paddingLeft: "24px"}}
                >
                  <div
                    style={{
                      fontWeight: 550,
                      marginBottom: "8px",

                      width: "max-content",
                      borderBottom: "2px dotted rgba(204, 204, 204, 1)",
                    }}
                  >
                    <Tooltips
                      width={244}
                      title={intl.get("Pretection Orders")}
                      text={
                        <Box style={{padding: "8px", width: "244px"}}>
                          <Box>
                            <strong>Pretection Orders</strong>
                          </Box>
                          <Box />
                        </Box>
                      }
                     />
                  </div>
                  <span
                    className="statistical_order_line_num"
                    style={{fontSize: "20px", fontWeight: 650}}
                  >
                    {orderData.orderStatics?.orders}
                  </span>
                </div>
                <div className="statistical_order_line_title">
                  <div
                    style={{
                      fontWeight: 550,
                      marginBottom: "8px",
                      fontSize: "13px",
                      width: "max-content",
                      borderBottom: "2px dotted rgba(204, 204, 204, 1)",
                    }}
                  >
                    <Tooltips
                      width={244}
                      title={intl.get("Insured Sales")}
                      text={
                        <Box style={{padding: "8px", width: "244px"}}>
                          <Box>
                            <strong>Insured Sales</strong>
                          </Box>
                          <Box />
                        </Box>
                      }
                     />
                  </div>
                  <span
                    className="statistical_order_line_num"
                    style={{fontSize: "20px", fontWeight: 650}}
                  >
                    $ {orderData.orderStatics?.sales}
                  </span>
                </div>
                <div className="statistical_order_line_title">
                  <div
                    style={{
                      fontWeight: 550,
                      marginBottom: "8px",
                      fontSize: "13px",
                      width: "max-content",
                      borderBottom: "2px dotted rgba(204, 204, 204, 1)",
                    }}
                  >
                    <Tooltips
                      width={244}
                      title={intl.get("Refund Amount")}
                      text={
                        <Box style={{padding: "8px", width: "244px"}}>
                          <Box>
                            <strong>Refund Amount</strong>
                          </Box>
                          <Box />
                        </Box>
                      }
                     />
                  </div>
                  <span
                    className="statistical_order_line_num"
                    style={{fontSize: "20px", fontWeight: 650}}
                  >
                    $ {orderData.orderStatics?.refund}
                  </span>
                </div>
                <div className="statistical_order_line_title">
                  <div
                    style={{
                      fontWeight: 550,
                      marginBottom: "8px",
                      fontSize: "13px",
                      width: "max-content",
                      borderBottom: "2px dotted rgba(204, 204, 204, 1)",
                    }}
                  >
                    <Tooltips
                      width={244}
                      title={intl.get("Total Revenue")}
                      text={
                        <Box style={{padding: "8px", width: "244px"}}>
                          <Box>
                            <strong>Total Revenue</strong>
                          </Box>
                          <Box />
                          <Box
                            style={{
                              padding: "6px 7px 6px 7px",
                              borderRadius: "4px",
                              backgroundColor: "rgba(243, 243, 243, 1)",
                              color: "rgba(0, 91, 211, 1)",
                              marginTop: "5px",
                            }}
                           />
                        </Box>
                      }
                     />
                  </div>
                  <span
                    className="statistical_order_line_num"
                    style={{
                      fontSize: "20px",
                      fontWeight: 650,
                      color:
                        orderData.orderStatics?.total >= 0
                          ? "#303030"
                          : "#f00",
                    }}
                  >
                    $ {orderData.orderStatics?.total}
                  </span>
                </div>

              </InlineStack>
            </div>

            <div className="statistical_order_line_arrowup">
              {getStartClo ? (
                <Button
                  variant="tertiary"
                  onClick={() => {
                    setGetStartClo(false);
                  }}
                >
                  <Icon source={ChevronUpIcon} tone="base" />
                </Button>
              ) : (
                <Button
                  variant="tertiary"
                  onClick={() => {
                    setGetStartClo(true);
                  }}
                >
                  <Icon source={ChevronDownIcon} tone="base" />
                </Button>
              )}
            </div>
          </div>

          <Collapsible
            open={getStartClo}
            id="basic-collapsible"
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
                option={opction}
                notMerge
                lazyUpdate
                theme="theme_name"
                onChartReady={() => {
                }}
                onEvents={() => {
                }}
                opts={{renderer: "canvas"}}
              />

            </div>
          </Collapsible>
        </div>
      </Card>
    </div>
  );
}
