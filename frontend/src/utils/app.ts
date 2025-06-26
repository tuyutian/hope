export const isProductionEnv = function () {
  // return true;
  return import.meta.env.PROD;
};