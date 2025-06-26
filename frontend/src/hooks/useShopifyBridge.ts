import { useAppBridge } from "@shopify/app-bridge-react";

export const useShopifyBridge = () => {
    const shopify = useAppBridge();
    if (typeof shopify !== "object") {
        return null;
    }else{
        return shopify;
    }
}


export function isShopifyEmbedded() {
    const appBridge = window.shopify;
    if (typeof appBridge !== "object") {
        return false;
    } else {
        return appBridge.environment.embedded || appBridge.environment.mobile;
    }
}