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
  const handleNavigate = (event: React.MouseEvent, path: string) => {
    event.preventDefault();

    navigate({pathname: path})
  };
  return (
    <div >
      <s-heading>
        {appBridge ? (<NavMenu>
          <a href="/" onClick={(e)=>{ handleNavigate(e,""); }} rel="home">
            Home
          </a>
          {Object.values(links).map((link,index) => {
            return link.url !=="/" &&<a key={index} href={link.url} onClick={(e)=>{ handleNavigate(e,link.url); }}>
              {link.name}
            </a>
          })}

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
