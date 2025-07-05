import {Outlet, useLocation} from "react-router";
import {useShopifyBridge} from "@/hooks/useShopifyBridge";
import {useEffect} from "react";
import {NavMenu} from "@shopify/app-bridge-react";
import i18n from "@/lib/i18n.ts";

const MainLayout = () => {
  const appBridge = useShopifyBridge();
  const {pathname, search, hash} = useLocation();
  useEffect(() => {
    if (appBridge) {
      const fullPath = `${pathname}${search}${hash}`;
      history.replaceState(null, "", fullPath);
      appBridge.loading(false);
    }
    // 跳转后，返回页面顶部
    window.scrollTo(0, 0);
  }, [appBridge, hash, pathname, search]);
  return (
    <div className="layout-container">
      {appBridge ? (<NavMenu>
        <a href="/plans">{i18n.get("Protection Plans") as string}</a>
        <a href="/cart">{i18n.get("Protection Page") as string}</a>
        <a href="/orders">{i18n.get("Orders") as string}</a>
        <a href="/billing">{i18n.get("Billing") as string}</a>
      </NavMenu>):<s-heading />}


      <main className="main-content">
        <Outlet />
      </main>

      <footer className="footer">
        <p>© {new Date().getFullYear()} 我的 React 19 应用</p>
      </footer>
    </div>
  );
};

export default MainLayout;
