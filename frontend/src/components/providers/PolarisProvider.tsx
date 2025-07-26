import React, {ReactNode, useCallback, useEffect, useMemo} from "react";
import {AppProvider} from "@shopify/polaris";
import enTranslations from "@shopify/polaris/locales/en.json";
import cnTranslations from "@shopify/polaris/locales/zh-CN.json";
import esTranslations from "@shopify/polaris/locales/es.json";
import frTranslations from "@shopify/polaris/locales/fr.json";
import deTranslations from "@shopify/polaris/locales/de.json";
import itTranslations from "@shopify/polaris/locales/it.json";
import "@shopify/polaris/build/esm/styles.css";

import {useNavigate} from "react-router";
import {LinkLikeComponentProps} from "@shopify/polaris/build/ts/src/utilities/link/types";
import {useShopifyBridge} from "@/hooks/useShopifyBridge";
import {getUserState} from "@/stores/userStore.ts";

function AppBridgeLink({url, children, external, ...rest}: LinkLikeComponentProps) {
  const navigate = useNavigate();
  const handleClick = useCallback(() => {
    navigate(url);
  }, [navigate, url]);

  const IS_EXTERNAL_LINK_REGEX = /^(?:[a-z][a-z\d+.-]*:|\/\/)/;

  if (external || IS_EXTERNAL_LINK_REGEX.test(url)) {
    return (
      <a {...rest} href={url} target="_blank" rel="noopener noreferrer">
        {children}
      </a>
    );
  }

  return (
    <a {...rest} onClick={handleClick}>
      {children}
    </a>
  );
}

/**
 * Sets up the AppProvider from Polaris.
 * @desc PolarisProvider passes a custom link component to Polaris.
 * The Link component handles navigation within an embedded app.
 * Prefer using this vs any other method such as an anchor.
 * Use it by importing Link from Polaris, e.g:
 *
 * ```
 * import {Link} from '@shopify/polaris'
 *
 * function MyComponent() {
 *  return (
 *    <div><Link url="/tab2">Tab 2</Link></div>
 *  )
 * }
 * ```
 *
 * PolarisProvider also passes translations to Polaris.
 *
 */
export const PolarisProvider = ({children}: { children: ReactNode }) => {
  const appBridge = useShopifyBridge();
  const setLang = getUserState().setLang;
  const lang = getUserState().lang;

  useEffect(() => {
    if (appBridge) {
      const changeLang = async function changeLang(userLang: string) {
        //TODO 只请求非英文语言数据加快反应
        // return appService.getLangData(lang).then(languageData => {
        //     if (languageData === null) return;
        //     const currentLang:string = languageData.lang;
        //     const langData = languageData.data;
        //     const languages = [];
        //     for (const set in langData) {
        //         languages.push({value: set, label: langData[set].name});
        //         i18n.addResourceBundle(set, "translation", langData[set].translation);
        //     }
        //     stores.commonStore.setLanguages(languages);
        //     stores.commonStore.setLang(currentLang);
        //     void i18n.changeLanguage(currentLang);
        //     setLang(currentLang)
        // });
        setLang(userLang);
      };
      void changeLang(appBridge.config.locale ?? "en").then();
    } else {
      // void appService.getLangData(appLang).then(languageData => {
      //     if (languageData === null) return;
      //     const currentLang:string = languageData.lang;
      //     const langData = languageData.data;
      //     const languages = [];
      //     for (const set in langData) {
      //         languages.push({value: set, label: langData[set].name});
      //         i18n.addResourceBundle(set, "translation", langData[set].translation);
      //     }
      //     stores.commonStore.setLanguages(languages);
      //     void i18n.changeLanguage(currentLang);
      //     setLang(currentLang)
      // })
    }
  }, [appBridge, lang]);

  const trans = useMemo(function () {
    switch (lang) {
      case "de":
        return deTranslations;
      case "es":
        return esTranslations;
      case "it":
        return itTranslations;
      case "cn":
        return cnTranslations;
      case "fr":
        return frTranslations;
      default:
        return enTranslations;
    }
  }, [lang]);

  return (
    <AppProvider i18n={trans} linkComponent={AppBridgeLink}>
      {children}
    </AppProvider>
  );
};
