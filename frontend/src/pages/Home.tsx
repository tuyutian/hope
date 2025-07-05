import {Text,} from "@shopify/polaris";

import NewPersonBox from "@/pages/Home/NewPersonBox";
import FulfilledOrders from "@/pages/Home/FulfilledOrders";
import AbilityBox from "@/pages/Home/AbilityBox";
import "@/pages/Home/index.css"

const DashboardClass = () => {

  return (
    <s-page>
      <s-stack gap="base">
      <s-box >
        <s-box style={{display: "flex", alignItems: "center", flexWrap: "wrap", fontWeight: 650, fontSize: "16px"}}>
          <Text fontWeight="bold" as="h1" variant="headingLg">👋 Hi, Welcome to Goodcare Protection
          </Text>
        </s-box>
      </s-box>

      <s-box>
        <FulfilledOrders />
      </s-box>

      <s-box >
        <NewPersonBox />
      </s-box>

      <s-box >
        <AbilityBox />
      </s-box>
      </s-stack>
    </s-page>
  );
};

export default DashboardClass;
