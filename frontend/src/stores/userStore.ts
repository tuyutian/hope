import { create } from "zustand";
import { UserGuide } from "@/types/user.ts";

interface UserState {
  domain: string;
  plan: number;
  name: string;
  email: string;
  authToken: string;
  token: string;
  lang: string;
  userGuide: UserGuide;
  guideShow: boolean;
}

interface UserActions {
  updateAuthToken: (authToken: string) => void;
  updateUserToken: (token: string) => void;
  setDomain: (domain: string) => void;
  setLang: (lang: string) => void;
  setUserGuide: (userGuide: UserGuide) => void;
  closeGuide: () => void;
  setGuideShow: (show: boolean) => void;
  toggleUserGuideStep: (name: keyof UserGuide, check: boolean) => void;
}

type UserStore = UserState & UserActions;

const useUserStore = create<UserStore>((set) => ({
  // Initial state
  domain: "",
  plan: 0,
  lang: "en",
  name: "",
  email: "",
  token: "",
  authToken: "",
  userGuide: {
    enabled: false,
    setting_protension: false,
    setup_widget: false,
    how_work: false,
    choose: false,
  },
  guideShow: true,

  // Actions
  setDomain: (domain) => set({ domain }),
  
  setLang: (lang) => set({ lang }),
  
  updateAuthToken: (authToken) => set({ authToken }),
  
  updateUserToken: (token) => set({ token }),
  
  setGuideShow: (show) => set({ guideShow: show }),

  closeGuide: () => set({ guideShow: false }),
  
  setUserGuide: (userGuide) => set({ userGuide }),
  
  toggleUserGuideStep: (name, check) => set((state) => ({
    userGuide: {
      ...state.userGuide,
      [name]: check
    }
  })),
}));

export default useUserStore;
export const userStore = useUserStore;
export const getUserState = () => useUserStore.getState();