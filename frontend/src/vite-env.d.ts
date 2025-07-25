/// <reference types="vite/client" />

interface ImportMetaEnv {
  VITE_SERVICE_NAME: string;
  VITE_SHOPIFY_APP_KEY: string;
  VITE_API_BASE_URL: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
