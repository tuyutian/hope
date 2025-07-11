import {Box,} from "@shopify/polaris";
import React from "react";

import NewPersonBox from "./NewPersonBox";
import FulfilledOrders from "./FulfilledOrders";
import AbilityBox from "./AbilityBox";

const DashboardClass = () => {
    return <s-page>
        <Box className="home_title home_title_mobile">
            <Box style={{display: 'flex', alignItems: 'center', flexWrap: 'wrap', fontWeight: 650, fontSize: '16px'}}>
                <h1 style={{fontWeight: 650, fontSize: '16px', marginRight: "8px"}}>👋 Hi, Welcome to Goodcare Protection
                </h1>
            </Box>
        </Box>

        <Box className="order_statistics" style={{marginBottom: "20px"}}>
            <FulfilledOrders />
        </Box>

        <Box className="guide" style={{marginBottom: "20px"}}>
            <NewPersonBox
            />
        </Box>

        <Box className="card_guidance" style={{marginBottom: "20px"}}>
            <AbilityBox
            />
        </Box>
    </s-page>
};

export default DashboardClass;
