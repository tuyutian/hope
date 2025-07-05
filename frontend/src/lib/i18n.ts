import LanguageDetector from 'i18next-browser-languagedetector';
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

i18n
  .use(LanguageDetector) // 嗅探当前浏览器语言
  .use(initReactI18next) // init i18next
  .init({
    // 引入资源文件
    resources: {},
    // 选择默认语言，选择内容为上述配置中的key，即en/zh/ja
    fallbackLng: 'en',
    debug: false,
    // partialBundledLanguages:true,
    interpolation: {
      escapeValue: false, // not needed for react as it escapes by default
    },
  })
  .then((r) => { console.log(r); });

const get = (key: string, ...variables: any[]) => {
  const variablesObject = variables.reduce((acc, variable, index) => {
    acc[`v${index + 1}`] = variable;
    return acc;
  }, {} as Record<string, any>);

  return i18n.t(key, variablesObject);
};

const obj = {
  get: get
}

export default obj;