import {
  BlockStack,
  Box,
  Button,
  Card,
  Collapsible,
  Icon,
  Tooltip,
} from "@shopify/polaris";
import React, { useEffect, useState } from "react";
import { XIcon, ChevronUpIcon, ChevronDownIcon } from "@shopify/polaris-icons";

const AbilityBox = ({  data, setPageData, toastFun }) => {
  const abilityList = {
    level1: [
      {
        title: "Cart Page Protection Plan",
        dsc: "Enjoy exclusive partner perks with your Shopify Balance account. Get offers from Google, Apple, UPS, and other industry-leading partners to run your business.",
        isShow: true,
        button: "Setting Now",
        icon:  "/image/in_card01.png",
        link: "/home",
        id: 1,
      },
      {
        title: "Checkout Page Protection Plan",
        dsc: "Enjoy exclusive partner perks with your Shopify Balance account. Get offers from Google, Apple, UPS, and other industry-leading partners to run your business.",
        isShow: true,
        button: "Setting Now",
        link: "/home",
        id: 2,
        icon:  "/image/in_card01.png",
      },
    ],

  };
  const [relAbilityList, setRelAbilityList] = useState([]); //当前展示的数组

  useEffect(() => {
    let isMount = false;
    let list = [];
    const list_key = "level1";
    list = abilityList[list_key]
        .slice(0, 2);
    setRelAbilityList(list);

    return () => {
      isMount = true;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [data]);

  const handleOffCard = async (id) => {
    try {
    } catch (error) {}
  };

  const handleViewDetail = (data) => {
    if (data.target) {
      window.open(data.target, "_blank");
    }
    if (data.button === "Consult Us") {
      const sendMsg = `Hi, I would like to inquire about Dropshipman Academy`;

      try {
        // 填充信息并激活对话弹窗
        // window.Intercom('showNewMessage', sendMsg)
      } catch (error) {
        console.info(error);
      }
    }
  };
  const handleShow = (id) => {
    const list = relAbilityList.map((item) => {
      if (item.id === id) {
        return {
          ...item,
          isShow: !item.isShow,
        };
      } else {
        return item;
      }
    });
    setRelAbilityList(list);
  };
  return (
    <div>
      <Box className="ability_box">
        <Box className="ability_box_content ability_box_content_pc">
          <BlockStack gap={400}>
            {relAbilityList.map((item, ind) => (
              <Card key={ind}>
                <Box className="instruction_card">
                  {/*<Box className="instruction_card_close">*/}
                  {/*  <Button*/}
                  {/*    variant="tertiary"*/}
                  {/*    onClick={() => handleOffCard(item.id)}*/}
                  {/*  >*/}
                  {/*    <Tooltip content="Dismiss">*/}
                  {/*      <Icon source={XIcon} tone="base" />*/}
                  {/*    </Tooltip>*/}
                  {/*  </Button>*/}
                  {/*</Box>*/}

                  <Box className="instruction_card_l">
                    <Box className="instruction_card_l_title">
                      {item.title}
                    </Box>
                    <Box className="instruction_card_l_dsc">
                      {item.dsc}
                    </Box>
                    <Box className="instruction_card_l_button">
                      <Button onClick={() => handleViewDetail(item)}>
                        {item.button}
                      </Button>
                    </Box>
                  </Box>
                  <Box className="instruction_card_r">
                      <img
                        width={168}
                        height={123}
                        src={item.icon}
                        alt="Dropshipman"
                      />
                  </Box>
                </Box>
              </Card>
            ))}
          </BlockStack>
        </Box>
        <Box className="ability_box_content instruction_card_mobile">
          <BlockStack gap={400}>
            {relAbilityList.map((item, ind) => (
              <Card key={ind}>
                <Box
                  className="instruction_card"
                  style={{ display: "block" }}
                >
                  <Box
                    className="instruction_card_close_button"
                    style={{
                      display: "flex",
                      justifyContent: "space-between",
                      alignItem: "center",
                    }}
                  >
                    <Box className="instruction_card_l_title">
                      {item.title}
                    </Box>
                    {/*<Box style={{ marginRight: "8px" }}>*/}
                    {/*  <Button*/}
                    {/*    variant="tertiary"*/}
                    {/*    onClick={() => handleOffCard(item.id)}*/}
                    {/*  >*/}
                    {/*    <Tooltip content="Dismiss">*/}
                    {/*      <Icon source={XIcon} tone="base" />*/}
                    {/*    </Tooltip>*/}
                    {/*  </Button>*/}
                    {/*</Box>*/}
                    {item.isShow ? (
                      <Button
                        variant="tertiary"
                        onClick={() => {
                          handleShow(item.id);
                        }}
                      >
                        <Icon source={ChevronUpIcon} tone="base" />
                      </Button>
                    ) : (
                      <Button
                        variant="tertiary"
                        onClick={() => {
                          handleShow(item.id);
                        }}
                      >
                        <Icon source={ChevronDownIcon} tone="base" />
                      </Button>
                    )}
                  </Box>

                  <Collapsible
                    open={item.isShow}
                    id="basic-collapsible"
                    transition={{
                      duration: "500ms",
                      timingFunction: "ease-in-out",
                    }}
                    expandOnPrint
                  >
                    <Box className="instruction_card_l">
                      <Box className="instruction_card_l_dsc">
                        {item.dsc}
                      </Box>
                      <Box className="instruction_card_l_button">
                        <Button onClick={() => handleViewDetail(item)}>
                          {item.button}
                        </Button>
                      </Box>
                    </Box>
                    <Box className="instruction_card_r">
                        {" "}
                        <img
                          width={168}
                          height={123}
                          src={item.icon}
                          alt="Dropshipman"
                        />
                    </Box>
                  </Collapsible>
                </Box>
              </Card>
            ))}
          </BlockStack>
        </Box>
      </Box>
    </div>
  );
};

export default AbilityBox;
