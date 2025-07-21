import { BaseApiService } from "../base";
import type { ApiResponse } from "@/types/api.ts";

export interface CartSettingsData {
  plan_title: string;
  addon_title: string;
  enabled_desc: string;
  disabled_desc: string;
  foot_text: string;
  foot_url: string;
  in_color: string;
  out_color: string;
  show_cart: number;
  show_cart_icon: number;
  select_button: number;
  pricing_type: number;
  price_rule: number;
  price_select: any[];
  tiers_select: any[];
  other_money: number;
  all_price: number;
  all_tiers: number;
  product_collection: any[];
  icons: any[];
  in_collection: boolean;
}

export interface UpdateCartSettingsParams {
  planTitle: string;
  iconVisibility: number;
  insuranceVisibility: number;
  selectButton: number;
  addonTitle: string;
  enabledDescription: string;
  disabledDescription: string;
  footerText: string;
  footerUrl: string;
  optInColor: string;
  optOutColor: string;
  pricingType: number;
  pricingRule: number;
  priceSelect: any[];
  tiersSelect: any[];
  restValuePrice: string;
  allPrice: string;
  allTiers: string;
  selectedCollections: any[];
  icons: any[];
  onlyInCollection: boolean;
}

export class CartService extends BaseApiService {
  constructor() {
    super("api/v1/setting/");
  }

  // 获取购物车设置
  getSettings(): Promise<ApiResponse<CartSettingsData>> {
    return this.get("cart");
  }

  // 更新购物车设置
  updateSettings(params: UpdateCartSettingsParams): Promise<ApiResponse> {
    return this.post("cart", params);
  }
}

// 导出实例
export const cartService = new CartService();
