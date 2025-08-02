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
  onCheckboxChange: () => void;
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
  onCheckboxChange,
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
          <input type="checkbox" checked={checkboxInput} onChange={onCheckboxChange} />
          <span
            className="checkmark"
            style={{
              backgroundColor: checkboxInput ? optInColor : optOutColor,
              borderColor: checkboxInput ? optInColor : optOutColor,
            }}
          />
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
          <Button fullWidth size="large">
            Checkout {checked || checkboxInput ? "22.00 USD" : "20.00 USD"}
          </Button>
        </BlockStack>
      </Box>
    </Card>
  );
};

export default CartDemo;
