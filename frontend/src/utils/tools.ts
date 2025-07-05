export function hexToRgb(color: string) {
  if (color.length === 4) {
    const repeatHex = (hex1: number, hex2: number) =>
      color.slice(hex1, hex2).repeat(2);
    const red = parseInt(repeatHex(1, 2), 16);
    const green = parseInt(repeatHex(2, 3), 16);
    const blue = parseInt(repeatHex(3, 4), 16);

    return {red, green, blue};
  }

  const red = parseInt(color.slice(1, 3), 16);
  const green = parseInt(color.slice(3, 5), 16);
  const blue = parseInt(color.slice(5, 7), 16);

  return {red, green, blue};
}

export function formatTimestampToUSDate(timestamp:number) {
  const date = new Date(timestamp * 1000); // 转换为毫秒
  const options:Intl.DateTimeFormatOptions = {
    weekday: "short", // 简短的星期几
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    hour12: true, // 12小时制
    timeZone: "America/New_York"
  };

  // 使用美国东部时间 (Eastern Time Zone)
  return date.toLocaleString("en-US", options);
}

export function getShopNameFromDomain(shopDomain: string) {
  const shopMatch = /(\S*).myshopify.com/.exec(shopDomain);
  if (shopMatch) {
    return shopMatch[1];
  }
  return shopDomain;
}

export function isEmpty(value: any): boolean {
  if (Array.isArray(value)) {
    return value.length === 0;
  } else {
    return value === "" || value == null;
  }
}

export const getWeekDaysConfig = (t: any) => [
  {label: t("001046", "Monday"), value: "Mon"},
  {label: t("001047", "Tuesday"), value: "Tue"},
  {label: t("001048", "Wednesday"), value: "Wed"},
  {label: t("001049", "Thursday"), value: "Thu"},
  {label: t("001050", "Friday"), value: "Fri"},
  {label: t("001051", "Saturday"), value: "Sat"},
  {label: t("001052", "Sunday"), value: "Sun"},
];
