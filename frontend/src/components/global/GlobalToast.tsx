import {Message} from "@/types/notify.ts";
import {useMessageStore} from "@/stores/messageStore.ts";
import {Toast} from "@shopify/polaris";

export const GlobalToast = function () {
  const message = useMessageStore(state => state.message);
  const toastComplete = useMessageStore(state => state.toastComplete);

  const completeToast = (toast: Message) => {
    if (toast.id) {
      toastComplete(toast.id);
    }
    return null;
  };

  if (message.length === 0) return null;

  return (
    <>
      {message.map((toast: Message) =>
        toast.content !== "" ? (
          <Toast
            key={toast.id}
            content={toast.content}
            error={toast.error}
            onDismiss={() => completeToast(toast)}
            duration={toast.duration}
          />
        ) : (
          completeToast(toast)
        )
      )}
    </>
  );
}

export default GlobalToast;