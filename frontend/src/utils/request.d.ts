export interface IExtParam {
  // 判定取消类型方法
  cancelType?: CancelTypesEnum
}

export enum CancelTypesEnum {
  ALL = "all",// 根据 method data 等判定重读请求
  PATH = "path",// 根据路径判定重复请求
}
