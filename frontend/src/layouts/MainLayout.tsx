import {Outlet, useLocation, useNavigate} from "react-router";
import {useShopifyBridge} from "@/hooks/useShopifyBridge";
import {useEffect, useState} from "react";
import {NavMenu} from "@shopify/app-bridge-react";
import i18n from "@/lib/i18n.ts";
import {TabProps, Tabs} from "@shopify/polaris";

const links = [
  {
    name: i18n.get("Home") as string,
    url: "/",
  },
  {
    name: i18n.get("Protection Page") as string,
    url: "/cart",
  },
  {
    name: i18n.get("Orders") as string,
    url: "/orders",
  },
  {
    name: i18n.get("Billing") as string,
    url: "/billing",
  }
];

const MainLayout = () => {
  const appBridge = useShopifyBridge();
  const {pathname, search, hash} = useLocation();
  const navigate = useNavigate();
  const [selected, setSelected] = useState(links.findIndex((link) => link.url === pathname) || 0);
  const tabs: TabProps[] = links.map((link) => {
    return {
      id: link.url,
      content: link.name,
      panelID: link.url,
    } as TabProps;
  });
  const handleTabChange = async (id: number) => {
    setSelected(id);
    await navigate(links[id].url);
  };
  useEffect(() => {
    if (appBridge) {
      const fullPath = `${pathname}${search}${hash}`;
      history.replaceState(null, "", fullPath);
      appBridge.loading(false);
    }
    // 跳转后，返回页面顶部
    window.scrollTo(0, 0);
  }, [pathname]);

  return (
    <div >
      <s-heading>
        {appBridge ? (<NavMenu>
          <a href="/cart">{i18n.get("Protection Page") as string}</a>
          <a href="/orders">{i18n.get("Orders") as string}</a>
          <a href="/billing">{i18n.get("Billing") as string}</a>
        </NavMenu>) : <div>
          <Tabs
            tabs={tabs}
            selected={selected}
            onSelect={handleTabChange}
            disclosureText="More views"
          />
        </div>}
      </s-heading>

      <div className="pb-10">
        <Outlet />
      </div>

    </div>
  );
};

export default MainLayout;
