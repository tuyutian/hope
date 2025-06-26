import {create} from "zustand";

type Store = {
  domain: string;
  plan: number;
  name: string;
  email: string;
  phone: string;
}
const userStore = create<Store>((set) => {
  return {
    domain: "",
    plan: 0,
    name: "",
    email: "",
    phone: "",
    setDomain: (domain:string) => set((state) => {
      state.domain = domain
      return state
    })
  };
});