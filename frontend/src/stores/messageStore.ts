import { create } from "zustand";
import { Message } from "~/types/notify";

type State = {
  message: Message[];
  messageNum: number;
};
type Action = {
  toastMessage: (payload: string | Message, duration?: number, error?: boolean) => void;
  toastComplete: (payload: number) => void;
};

export const useMessageStore = create<State & Action>(set => {
  return {
    message: [],
    messageNum: 0,
    toastMessage: (payload: string | Message, duration = 2000, error = false) => {
      console.log(payload);
      if (typeof payload === "string") {
        console.log(shopify);
        if (typeof shopify === "object") {
          shopify.toast.show(payload, {
            duration: duration,
            isError: error,
          });
          return;
        }
        set(state => ({
          message: [
            ...state.message,
            {
              content: payload,
              error,
              duration,
              id: Math.floor(Math.random() * 100 + 1),
            },
          ],
          messageNum: state.message.length + 1,
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
        message: [
          ...state.message,
          {
            content: payload.content,
            error: payload.error,
            duration: payload.duration ? payload.duration : 2000,
            id: Math.floor(Math.random() * 100 + 1),
          },
        ],
        messageNum: state.message.length + 1,
      }));
    },
    toastComplete: (payload: number) => {
      set(state => {
        const index = state.message.findIndex(e => e.id === payload);
        const newMessage = [...state.message];
        newMessage.splice(index, 1);
        return {
          message: newMessage,
          messageNum: newMessage.length,
        };
      });
    },
  };
});
export default useMessageStore;

// 导出 store 实例
export const messageStore = useMessageStore;

// 导出获取状态的函数（非 hook）
export const getMessageState = () => useMessageStore.getState();
