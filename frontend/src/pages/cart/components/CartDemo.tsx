import React from "react";
import { BlockStack, Box, Button, Card, InlineStack, Text, Thumbnail } from "@shopify/polaris";
import CustomSwitch from "./CustomSwitch";

interface Icon {
  id: number;
  src: string;
  selected: boolean;
}

interface CartDemoProps {
  iconVisibility: string;
  selectedIcon: Icon | undefined;
  addonTitle: string;
  enabledDescription: string;
  disabledDescription: string;
  footerText: string;
  footerUrl: string;
  selectButton: string;
  checkboxInput: boolean;
  optInColor: string;
  optOutColor: string;
}

const CartDemo: React.FC<CartDemoProps> = ({
  iconVisibility,
  selectedIcon,
  addonTitle,
  enabledDescription,
  disabledDescription,
  footerText,
  footerUrl,
  selectButton,
  checkboxInput,
  optInColor,
  optOutColor,
}) => {
  const [checked, setChecked] = React.useState(false);
  const renderSelectionControl = () => {
    if (selectButton === "0") {
      return (
        <CustomSwitch
          onChange={setChecked}
          checked={checked}
          onColor={optInColor}
          offColor={optOutColor}
          checkedIcon={false}
        />
      );
    } else {
      return (
        <label className="custom-checkbox">
          <input type="checkbox" className="absolute" checked={checked} onChange={e => setChecked(e.target.checked)} />
          <span
            className="checkmark"
            style={{
              backgroundColor: checked ? optInColor : optOutColor,
              borderColor: checked ? optInColor : optOutColor,
              transition: "all 0.2s ease-in-out",
              boxShadow: checked ? `0 0 5px ${optInColor}40` : `0 0 5px ${optOutColor}40`,
              transform: checked ? "scale(1.05)" : "scale(1)",
              width: "20px",
              height: "20px",
              borderRadius: "4px",
              display: "inline-block",
              position: "relative",
              cursor: "pointer",
            }}
          >
            {checked && (
              <span
                style={{
                  position: "absolute",
                  top: "45%", // 调整到略微靠上的位置
                  left: "50%",
                  transform: "translate(-50%, -50%) rotate(-45deg)", // 合并旋转到 transform 中
                  width: "12px", // 调整勾选标记的宽度
                  height: "7px", // 调整勾选标记的高度
                  borderLeft: "2px solid white",
                  borderBottom: "2px solid white",
                  display: "block", // 确保显示为块级元素
                  margin: "0 auto", // 水平居中
                }}
              />
            )}
          </span>
        </label>
      );
    }
  };

  return (
    <Card padding="0">
      <Box background="bg-fill-tertiary" padding="400" borderStartStartRadius="100" borderStartEndRadius="100">
        <InlineStack gap="400" align="space-between">
          <Text as="h6" variant="bodyMd" fontWeight="semibold">
            Cart Page Demo
          </Text>
          <Button variant="tertiary">View in store</Button>
        </InlineStack>
      </Box>
      <Box padding="400">
        <BlockStack gap="300">
          {/* Mock Products */}
          {[1, 2].map((item, idx) => (
            <Box key={idx} padding="300" background="bg-fill-secondary">
              <InlineStack gap="300" align="start">
                <Thumbnail
                  source="https://img.icons8.com/plasticine/100/cat-footprint.png"
                  alt="Cat Slippers"
                  size="medium"
                />
                <BlockStack>
                  <Text as="p" variant="bodyMd" fontWeight="medium">
                    Cute Cat Slippers
                  </Text>
                  <Text as="span">$10.00</Text>
                </BlockStack>
              </InlineStack>
            </Box>
          ))}

          {/* Protection Option */}
          <Card padding="300">
            <InlineStack wrap={false} gap="300" blockAlign="start">
              {iconVisibility === "1" && selectedIcon ? (
                <img src={selectedIcon.src} alt="Protection" style={{ flexShrink: 0, width: "66px", height: "66px" }} />
              ) : (
                <div />
              )}
              <BlockStack gap="150">
                <InlineStack blockAlign="start" wrap={false} align="space-between">
                  <BlockStack align="start">
                    <Box>
                      <Text as="span" variant="bodyMd" fontWeight="semibold">
                        {addonTitle || "Shipping Protection"}
                      </Text>
                      <Text tone="subdued" as="span">
                        (2.00 USD)
                      </Text>
                    </Box>
                    <Text as="p" tone="subdued" variant="bodySm">
                      {(selectButton === "0" && checked) || (selectButton === "1" && checkboxInput)
                        ? enabledDescription
                        : disabledDescription}
                    </Text>
                    {/* Footer Link */}
                    {footerText && (
                      <Box>
                        <Text variant="bodySm" as="span">
                          <a
                            title={footerUrl || ""}
                            style={{
                              color: "#0070f3",
                              textDecoration: "none",
                              cursor: "pointer",
                            }}
                            onClick={e => e.preventDefault()}
                          >
                            {footerText}
                          </a>
                        </Text>
                      </Box>
                    )}
                  </BlockStack>
                  {renderSelectionControl()}
                </InlineStack>
              </BlockStack>
            </InlineStack>
          </Card>

          {/* Checkout Button */}
          <Button fullWidth size="large" variant="primary">
            Checkout {checked || checkboxInput ? "22.00 USD" : "20.00 USD"}
          </Button>
        </BlockStack>
      </Box>
    </Card>
  );
};

export default CartDemo;
