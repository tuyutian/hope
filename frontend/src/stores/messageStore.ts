import {create} from "zustand";
import {Message} from "~/types";

type State = {
  message: Message[];
  rate: number;
  commentModalShow: boolean;
  scroll: string;
  newNotice: boolean;
  messageNum: number;

}
type Action = {
  setScroll: (scroll: string) => void;
  setCommentModal: (show: boolean) => void;
  setRate: (rate: number) => void;
  toastMessage: (payload: string | Message, duration?: number, error?: boolean) => void;
  toastComplete: (payload: number) => void;
  setNewNotice: (newNotice: boolean) => void;
  clearToast: () => void;
}

export const useMessageStore = create<State&Action>((set) => ({
  message: [],
  rate: 0,
  commentModalShow: false,
  scroll: "",
  newNotice: false,
  messageNum: 0,

  setScroll: (scroll: string) => set({scroll}),

  setCommentModal: (show: boolean) => set({commentModalShow: show}),

  setRate: (rate: number) => set({rate}),

  toastMessage: (payload: string | Message, duration = 2000, error = false) => {
    if (typeof payload === "string") {
      if (typeof shopify === "object") {
        shopify.toast.show(payload, {
          duration: duration,
          isError: error,
        });
        return;
      }
      set(state => ({
        message: [...state.message, {
          content: payload,
          error,
          duration,
          id: Math.floor(Math.random() * 100 + 1),
        }],
        messageNum: state.message.length + 1
      }));
      return;
    }
    if (typeof shopify === "object") {
      shopify.toast.show(payload.content, {
        duration: payload.duration ? payload.duration : duration,
        isError: payload.error,
      });
      return;
    }
    set(state => ({
      message: [...state.message, {
        content: payload.content,
        error: payload.error,
        duration: payload.duration ? payload.duration : 2000,
        id: Math.floor(Math.random() * 100 + 1),
      }],
      messageNum: state.message.length + 1
    }));
  },

  toastComplete: (payload: number) => {
    set(state => {
      const index = state.message.findIndex(e => e.id === payload);
      const newMessage = [...state.message];
      newMessage.splice(index, 1);
      return {
        message: newMessage,
        messageNum: newMessage.length
      };
    });
  },

  setNewNotice: (newNotice: boolean) => set({newNotice}),

  clearToast: () => set({message: [], messageNum: 0})
}));