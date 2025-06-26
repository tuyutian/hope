import {create} from "zustand";

type State = {
  domain: string;
  plan: number;
  name: string;
  email: string;
  authToken: string;
  token: string;
}
type Action = {
  updateAuthToken: (authToken:State['authToken']) => void;
  updateUserToken: (authToken:State['token']) => void;
  setDomain: (domain:State['domain']) => void;
}
const useUserStore = create<State&Action>((set) => {
  return {
    domain: "",
    plan: 0,
    name: "",
    email: "",
    token: "",
    authToken: "",
    setDomain: (domain:string) => set((state) => {
      state.domain = domain
      return state
    }),
    updateAuthToken: (authToken:string) => set((state) => {
      state.authToken = authToken
      return state
    }),
    updateUserToken: (token:string) => set((state) => {
      state.token = token
      return state
    }),
  };
});

export default useUserStore;