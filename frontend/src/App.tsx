import React, { useEffect } from "react";
import { isShopifyEmbedded, useShopifyBridge } from "@/hooks/useShopifyBridge";
import { isProductionEnv } from "@/utils/app.ts";
import { getUserState } from "@/stores/userStore.ts";
import { redirectRemote } from "@/utils/shopify.ts";
import Router from "./routes/Router.tsx";
import "@shopify/polaris/build/esm/styles.css";
import { PolarisProvider } from "@/components/providers/PolarisProvider.tsx";
import { ShopifyAuthProvider } from "@/components/providers/ShopifyAuthProvider.tsx";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ThemeProvider } from "@/components/providers/ThemeProvider.tsx";
const queryClient = new QueryClient();

function App() {
  const shopify = useShopifyBridge();
  const params = new URLSearchParams(location.search);
  const authToken = params.get("authToken");
  const token = params.get("token");
  const shop = params.get("shop");
  const useAuthToken = !!(authToken && authToken.length >= 32);
  const useAdminToken = !!(token && token.length >= 32);

  // 分别获取函数，避免创建新对象
  const updateAuthToken = getUserState().updateAuthToken;
  const updateUserToken = getUserState().updateUserToken;

  useEffect(() => {
    if (shopify) {
      shopify.loading(true);
    }

    // 开发环境设置
    if (!isShopifyEmbedded() && !isProductionEnv()) {
      console.log("%c开发环境,使用用户token", "background-color:#00a0ac; color: white; padding: 4px 8px;");
      updateAuthToken(import.meta.env.VITE_TEST_TOKEN as string);
    }

    // 用户进入
    if (useAuthToken) {
      updateAuthToken(authToken);
    } else if (useAdminToken) {
      updateUserToken(token);
    }

    // 重定向逻辑
    if (!isShopifyEmbedded() && isProductionEnv() && !useAdminToken && !useAuthToken) {
      let redirectUrl = `${import.meta.env.VITE_API_BASE_URL}/protectify/api/`;
      if (shop) {
        redirectUrl += `?shop=${shop}`;
      }
      redirectRemote(redirectUrl, true);
    }
  }, [authToken, token, shop, useAuthToken, useAdminToken, updateAuthToken, updateUserToken, shopify]);

  return (
    // Provide the client to your App
    <QueryClientProvider client={queryClient}>
      <PolarisProvider>
        <ShopifyAuthProvider>
          <ThemeProvider defaultTheme="light" storageKey="protectify-theme-key">
            <Router />
          </ThemeProvider>
        </ShopifyAuthProvider>
      </PolarisProvider>
    </QueryClientProvider>
  );
}

export default App;
