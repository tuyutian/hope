import { AxiosRequestConfig} from "axios";
import request, {IResponse} from "~/utils/request";
import {FilterParams} from "@/types/billing.ts";

/** 登录模块 */

export const ckeditorConfig = {
  toolbar: [
    "heading",
    "|",
    "bold",
    "italic",
    "link",
    "bulletedList",
    "numberedList",
    "blockQuote",
    "insertTable",
    "undo",
    "redo"
  ],
  heading: {
    options: [
      {model: "paragraph", title: "Paragraph", class: "ck-heading_paragraph"},
      {model: "heading1", view: "h1", title: "Heading 1", class: "ck-heading_heading1"},
      {model: "heading2", view: "h2", title: "Heading 2", class: "ck-heading_heading2"},
      {model: "heading3", view: "h3", title: "Heading 3", class: "ck-heading_heading3"},
      {model: "heading4", view: "h4", title: "Heading 4", class: "ck-heading_heading4"},
      // { model: 'heading5', view: 'h5', title: 'Heading 5', class: 'ck-heading_heading5' },
      // { model: 'heading6', view: 'h6', title: 'Heading 6', class: 'ck-heading_heading6' },
    ]
  },
  //   ckfinder: {
  //       uploadUrl: API_DOMAIN + '/api/v1/ckUpload'
  //   }
};
export const rqPostLogin = (params: any): Promise<IResponse> => request.post("api/v1/auth/login", params); // 登录
export const rqGetDashboard = (days:string): Promise<IResponse> => request.get(`api/v1/order/dashboard?days=${days}`); // 查询订单统计

export const rqGetCartSetting = (params?:any): Promise<IResponse> => request.get("api/v1/setting/cart", params); // 查询购物车配置

export const rqPostUpdateCartSetting = (params:any): Promise<IResponse> => request.post("api/v1/setting/cart", params); // 修改购物车配置

export const GetUserConf = (): Promise<IResponse> => request.get("api/v1/user/conf"); // 查询用户配置

export const GetSessionData= (): Promise<IResponse> => request.get("api/v1/user/session"); // 查询用户配置

export const rqGetOrderList = (params:any): Promise<IResponse> => request.get("api/v1/order/list", params); // 获取订单列表

export const reqGetUserAuthInstall = (params:any): Promise<IResponse> => request.get("api/v1/auth/register", params);

export const GetBillingData = (params:FilterParams): Promise<IResponse> => request.post("api/v1/billing/list",params);
export const GetBillingDetailData = (params:FilterParams): Promise<IResponse> => request.post("api/v1/billing/details",params);
export const GetCurrentPeriod = (): Promise<IResponse> => request.get("api/v1/billing/current");

export const UpdateDashboardGuide= (name:string,open:boolean): Promise<IResponse> => request.post('api/v1/user/step',{
  'name':name,
  'open':open,
})
