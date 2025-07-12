import {CalloutCard,} from "@shopify/polaris";
import React, {useMemo} from "react";

const AbilityBox = () => {
  const abilityList = {
    level1: [
      {
        title: "Cart Page Protection Plan",
        dsc: "Enjoy exclusive partner perks with your Shopify Balance account. Get offers from Google, Apple, UPS, and other industry-leading partners to run your business.",
        isShow: true,
        button: "Setting Now",
        icon: "/image/in_card01.png",
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
        icon: "/image/in_card01.png",
      },
    ],

  };
  const relAbilityList = useMemo(() => abilityList.level1.slice(0, 2), []);

  const handleViewDetail = (data) => {
    if (data.target) {
      window.open(data.target, "_blank");
    }
    if (data.button === "Consult Us") {

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
  return <>{relAbilityList.map((item, ind) => (
    <CalloutCard
      key={ind}
      title={item.title}
      illustration={item.icon}
      primaryAction={{
        content: item.button,
        onAction: () => handleViewDetail(item),
      }}
    >
      <p>{item.dsc}</p>
    </CalloutCard>
  ))}
  </>;
};

export default AbilityBox;
