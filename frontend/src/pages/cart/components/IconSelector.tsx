import React, { useCallback, useTransition } from "react";
import { DropZone, Icon, InlineStack, Spinner, Text, Thumbnail } from "@shopify/polaris";

interface Icon {
  id: number;
  src: string;
  selected: boolean;
}

interface IconSelectorProps {
  icons: Icon[];
  onIconClick: (id: number) => void;
  onIconUpload: (file: File) => void;
}

const IconSelector: React.FC<IconSelectorProps> = ({ icons, onIconClick, onIconUpload }) => {
  const [uploading, startTransition] = useTransition();
  const handleDropZoneDrop = useCallback((_dropFiles: File[], acceptedFiles: File[], _rejectedFiles: File[]) => {
    startTransition(() => {
      onIconUpload(acceptedFiles[0]);
    });
  }, []);

  return (
    <>
      <Text as="h2" variant="bodyMd">
        Widget Icon
      </Text>
      <InlineStack wrap gap="200" blockAlign="center">
        {icons.map((icon, index) => {
          const isSelected = icon.selected;
          return (
            <div
              key={index}
              onClick={() => onIconClick(icon.id)}
              style={{
                position: "relative",
                padding: "8px",
                borderRadius: "12px",
                border: isSelected ? "2px solid #5c6ac4" : "1px solid #dfe3e8",
                backgroundColor: isSelected ? "#f0f1f3" : "white",
                cursor: "pointer",
              }}
              onMouseEnter={e => {
                if (!isSelected) e.currentTarget.style.borderColor = "#b4bcc4";
              }}
              onMouseLeave={e => {
                if (!isSelected) e.currentTarget.style.borderColor = "#dfe3e8";
              }}
            >
              <Thumbnail source={icon.src} alt="Icon" size="medium" />
              {isSelected && (
                <div
                  style={{
                    position: "absolute",
                    top: "6px",
                    right: "6px",
                    backgroundColor: "#5c6ac4",
                    color: "white",
                    borderRadius: "50%",
                    width: "18px",
                    height: "18px",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    fontSize: "12px",
                    fontWeight: "bold",
                  }}
                >
                  ✓
                </div>
              )}
            </div>
          );
        })}
        <div style={{ width: 75, height: 75 }}>
          <DropZone type="image" allowMultiple={false} onDrop={handleDropZoneDrop}>
            {uploading ? (
              <div className="flex items-center justify-center cursor-default">
                <Spinner accessibilityLabel="Small spinner example" size="small" />
              </div>
            ) : (
              <DropZone.FileUpload />
            )}
          </DropZone>
        </div>
      </InlineStack>
    </>
  );
};

export default IconSelector;
