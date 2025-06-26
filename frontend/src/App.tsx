import React from "react";
import {useShopifyBridge, isShopifyEmbedded} from "@/hooks/useShopifyBridge.ts";
import {isProductionEnv} from "@/utils/app.ts";
import useUserStore from "@/stores/userStore.ts";
import {redirectRemote} from "@/utils/shopify.ts";

function App() {
  const shopify = useShopifyBridge();
  const params = new URLSearchParams(
    location.search
  ); /* Shopify 重定向过来，携带了这些参数 */
  const authToken = params.get("authToken"); /* 超管访问携带 */
  const token = params.get("token"); /* 非shopify 访问模式 */
  const shop = params.get("shop");
  const useAuthToken = !!(authToken && authToken.length >= 32);
  const useAdminToken = !!(token && token.length >= 32);
  if (shopify) {
    shopify.loading(true);
  }
  const {updateAuthToken, updateUserToken} = useUserStore(state => ({
    updateAuthToken: state.updateAuthToken,
    updateUserToken: state.updateUserToken
  }));
  if (!isShopifyEmbedded() && !isProductionEnv()) {
    console.log(
      "%c开发环境,使用用户token",
      "background-color:#00a0ac; color: white; padding: 4px 8px;"
    );
    updateAuthToken(import.meta.env.VITE_TEST_TOKEN as string);
  }
  // 用户进入
  if (useAuthToken) {
    updateAuthToken(authToken);
  } else if (useAdminToken) {
    updateUserToken(token);
  }
  if (
    !isShopifyEmbedded() &&
    isProductionEnv() &&
    !useAdminToken &&
    !useAuthToken
  ) {
    let redirectUrl = "https://tms.trackingmore.net";
    if (shop) {
      redirectUrl += `?shop=${shop}`;
    }
    redirectRemote(redirectUrl, true);
  }
  return (
    <s-page>
      <s-section>
        <s-text>Hello World</s-text>
      </s-section>
    </s-page>
  );
}

export default App;
