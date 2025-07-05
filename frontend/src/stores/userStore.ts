import { create } from 'zustand'

type State = {
  domain: string;
  plan: number;
  name: string;
  email: string;
  authToken: string;
  token: string;
  lang: string;
}

type Action = {
  updateAuthToken: (authToken: State['authToken']) => void;
  updateUserToken: (token: State['token']) => void;
  setDomain: (domain: State['domain']) => void;
  setLang: (lang: State['lang']) => void;
}

const useUserStore = create<State & Action>((set) => {
  return {
    domain: "",
    plan: 0,
    lang: "en",
    name: "",
    email: "",
    token: "",
    authToken: "",
    setDomain: (domain: string) => set((state) => ({
      ...state,
      domain
    })),
    setLang: (lang: string) => set((state) => ({
      ...state,
      lang
    })),
    updateAuthToken: (authToken: string) => set((state) => ({
      ...state,
      authToken
    })),
    updateUserToken: (token: string) => set((state) => ({
      ...state,
      token
    })),
  };
});

// 导出默认的 hook
export default useUserStore

// 导出 store 实例用于直接访问
export const userStore = useUserStore

// 导出获取状态的函数（非 hook）
export const getUserState = () => useUserStore.getState()