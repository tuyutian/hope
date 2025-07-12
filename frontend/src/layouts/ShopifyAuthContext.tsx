
import { AuthContext } from "@/stores/context";
import { DefaultUser, User } from "@/types/user.ts";
import { ReactNode, useState, useEffect } from "react";
import {GetSessionData} from "@/api";
import {useShopifyBridge} from "@/hooks/useShopifyBridge";

declare global {
  interface Window {
    hideLoadingState: () => void;
  }
}

interface ShopifyAuthContextProps {
  children: ReactNode;
}

export function ShopifyAuthContext({ children }: ShopifyAuthContextProps) {
  // 使用 React 19 的 useState 钩子初始化用户数据
  const [user, setUser] = useState<User>(DefaultUser);
  const shopify = useShopifyBridge()
  // 初始化用户数据的异步函数
  const initializeUser = async () => {
    try {
      // 这里可以调用 API 获取用户数据
      const userData = await GetSessionData();

      setUser({
        shop:userData.data.shop,
        userGuide:userData.data.guide_step,
      });
    } catch (err) {
      console.error('初始化用户数据错误:', err);
    } finally {
      if (shopify){
        shopify.loading(false)
      }
      // 隐藏加载状态
      if (window.hideLoadingState) {
        window.hideLoadingState();
      }
    }
  };

  // 使用 useEffect 在组件挂载时初始化用户数据
  useEffect(() => {
    console.log(123);
    void initializeUser();
  }, []);


  return (
    <AuthContext value={{user:user,setUser:setUser}}>
      {children}
    </AuthContext>
  );
}