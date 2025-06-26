import axios, { AxiosInstance, AxiosRequestConfig } from 'axios'

import { CancelTypesEnum, IExtParam } from './request.d'
import qs from 'qs'
import { useShopifyBridge } from '@/hooks/useShopifyBridge.ts'
import useUserStore from "@/stores/userStore.ts";
import {useMessageStore} from "@/stores/messageStore.ts";

// 返回res.data的interface

const axiosInstance: AxiosInstance = axios.create({
  timeout: import.meta.env.VITE_REQUEST_TIME_OUT,
  baseURL: import.meta.env.VITE_API_URL,
  headers: {
    Accept: 'application/json',
    'Content-Type': 'application/json;charset=utf-8',
  },
  withCredentials: true,
})

// axios实例拦截响应
axiosInstance.interceptors.response.use(
  response => {
    const toastMessage = useMessageStore(state => state.toastMessage)
    // 关闭加载
    removePending(response.config)

    if (response.status === 200) {
      return response.data
    } else {
      toastMessage({
        content: response.data.message,
        error: true,
      })
      return response.data
    }
  },
  // 请求失败
  error => {
    // 关闭加载
    if (axios.isCancel(error)) {
      // 中断 Promise 调用链
      return new Promise(() => undefined)
    }
    const { response } = error
    if (
      window.location.pathname !== '/unauthorized' &&
      response &&
      response.status === 401
    ) {
      open('/unauthorized', '_self')
    }
    const toastMessage = useMessageStore(state => state.toastMessage)
    if (response) {
      // 请求已发出，但是不在2xx的范围
      if (import.meta.env.DEV) {
        toastMessage({
          content: response.data.message,
          error: true,
        })
      }
      throw response.data
    } else {
      toastMessage({
        content: 'Internet Error,Please Try Later !',
        duration: 5000,
        error: true,
      })
      throw error
    }
  }
)

// axios实例拦截请求
axiosInstance.interceptors.request.use(
  config => {
    const extParam: IExtParam = {
      cancelType: config.data?.cancelType || CancelTypesEnum.PATH,
    }
    config.data?.cancelType && delete config.data.cancelType

    removePending(config, extParam) // 在请求开始前，对之前的请求做检查取消操作
    addPending(config, extParam) // 将当前请求添加到 pending 中
    const {authToken,userToken} = useUserStore(state => ({authToken:state.authToken,userToken: state.token}))
    const appBridge = useShopifyBridge()
    if (config.headers && authToken && authToken.length > 32) {
      config.headers.Authorization = `Bearer ${authToken}`
      return config
    }
    if (appBridge === null && import.meta.env.DEV) {
      if (config.headers) config.headers.Authorization = `Bearer ${userToken}`
      return config
    }
    if (appBridge === null) {
      if (config.headers) config.headers.Authorization = `Bearer ${userToken}`
      return config
    }
    return appBridge.idToken().then(token => {
      config.headers.Authorization = `Bearer ${token}`
      return config
    })
  },
  error => {
    console.log(error)
    return error
  }
)
// 声明一个 Map 用于存储每个请求的标识 和 取消函数
const pending = new Map()

/**
 * 添加请求
 * @param {Object} config
 * @param extParam
 */
const addPending = (config: AxiosRequestConfig, extParam?: IExtParam) => {
  const url = [config.method, config.url, qs.stringify(config.params)]
  // 如果取消的类型是 ALL 那么就要把 POST 参数加上
  if (extParam?.cancelType && extParam.cancelType === CancelTypesEnum.ALL) {
    url.push(qs.stringify(config.data))
  }

  const urlString: string = url.join('&')
  config.cancelToken =
    config.cancelToken ??
    new axios.CancelToken(cancel => {
      if (!pending.has(urlString)) {
        // 如果 pending 中不存在当前请求，则添加进去
        pending.set(urlString, cancel)
      }
    })
}

/**
 * 移除请求
 * @param {Object} config
 * @param extParam
 */
const removePending = (config: AxiosRequestConfig, extParam?: IExtParam) => {
  const url = [config.method, config.url, qs.stringify(config.params)]

  // 如果取消的类型是 ALL 那么就要把 POST 参数加上
  if (extParam?.cancelType && extParam.cancelType === CancelTypesEnum.ALL) {
    url.push(qs.stringify(config.data))
  }

  const urlString: string = url.join('&')
  if (pending.has(urlString)) {
    // 如果在 pending 中存在当前请求标识，需要取消当前请求，并且移除
    const cancel = pending.get(urlString)
    cancel(urlString)
    pending.delete(urlString)
  }
}

/**
 * 清空 pending 中的请求（在路由跳转时调用）
 */
export const clearPending = () => {
  // 路由切换前清空掉之前的请求
  // 可能会造成一些问题，暂时去除。2021-07-01 13:50
  for (const [url, cancel] of pending) {
    cancel(url)
  }
  pending.clear()
}

export default axiosInstance
