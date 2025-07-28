import { useTranslation } from "react-i18next";
import { ContextualSaveBar, Modal } from "@shopify/polaris";
import { SaveBar } from "@shopify/app-bridge-react";
import { useShopifyBridge } from "@/hooks/useShopifyBridge.ts";
import { useCallback, useState, useTransition } from "react";

type Props = {
  dirty: boolean;
  onSave: () => Promise<void>;
  onDiscard: () => void;
};

const PageSaveBar = ({ dirty, onSave, onDiscard }: Props) => {
  const [isPending, startTransition] = useTransition();
  const [isHanding, startDiscard] = useTransition();

  const appBridge = useShopifyBridge();
  const { t } = useTranslation();
  const [discardActive, setDiscardActive] = useState<boolean>(false);
  const handleDiscardActive = useCallback(() => {
    setDiscardActive(!discardActive);
  }, [discardActive]);

  const handleSave = async () => {
    console.log(123);
    await onSave();
    if (appBridge) {
      await appBridge.saveBar.hide("global-save-bar");
    }
  };
  const handleDiscard = async () => {
    await onDiscard();
    handleDiscardActive();
    if (appBridge) {
      await appBridge.saveBar.hide("global-save-bar");
    }
  };

  return (
    <>
      {appBridge ? (
        <SaveBar id="global-save-bar" open={dirty}>
          {/* eslint-disable-next-line react/no-unknown-property */}
          <button variant="primary" onClick={() => startTransition(handleSave)}>
            {t("001723", "Save")}
          </button>
          <button onClick={() => startDiscard(handleDiscardActive)}>{t("001724", "Discard")}</button>
        </SaveBar>
      ) : dirty ? (
        <ContextualSaveBar
          alignContentFlush
          message="Unsaved changes"
          saveAction={{
            onAction: handleSave,
            loading: isPending,
            content: t("001723", "Save"),
          }}
          discardAction={{
            onAction: handleDiscardActive,
            loading: isHanding,
            content: t("001724", "Discard"),
          }}
        />
      ) : (
        ""
      )}
      <Modal
        open={discardActive}
        onClose={handleDiscardActive}
        title={t("001725", "Discard all unsaved changes")}
        primaryAction={{
          destructive: true,
          content: t("001726", "Discard changes"),
          onAction: handleDiscard,
        }}
        secondaryActions={[
          {
            content: t("001727", "Continue editing"),
            onAction: handleDiscardActive,
          },
        ]}
      >
        <Modal.Section>
          <p>{t("001728", "If you discard changes, you’ll delete any edits you made since you last saved.")}</p>
        </Modal.Section>
      </Modal>
    </>
  );
};

export default PageSaveBar;
