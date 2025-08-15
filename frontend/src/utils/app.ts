export const isProductionEnv = function () {
  // return true;
  return import.meta.env.PROD;
};

export const handleContact = function (message: string = "") {
  return open(`mailto:support@protectifyapp.com?subject=Contact&body=${message}`);
};
