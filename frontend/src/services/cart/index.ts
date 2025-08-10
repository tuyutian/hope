import { BaseApiService } from "../base";
import type { ApiResponse } from "@/types/api.ts";
import { CartSettingsData, UpdateCartSettingsParams } from "@/types/cart.ts";

export class CartService extends BaseApiService {
  constructor() {
    super("v1/setting/");
  }

  // 获取购物车设置
  getSettings(): Promise<ApiResponse<CartSettingsData>> {
    return this.get("cart");
  }

  // 更新购物车设置
  updateSettings(params: UpdateCartSettingsParams): Promise<ApiResponse> {
    return this.post("cart", params);
  }

  uploadLogo(image: File): Promise<ApiResponse<{ id: number; src: string }>> {
    const formData = new FormData();
    formData.append("image", image);
    return this.post("upload_logo", formData, {
      timeout: 60000, // 文件上传需要更长时间
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });
  }
}

// 导出实例
export const cartService = new CartService();
