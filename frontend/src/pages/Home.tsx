import {BlockStack, Box, Collapsible, Icon, InlineStack, Text,} from "@shopify/polaris";

import NewPersonBox from "@/pages/home/NewPersonBox";
import FulfilledOrders from "@/pages/home/components/FulfilledOrders.tsx";
import "@/pages/home/index.css";
import {ChevronUpIcon} from "@shopify/polaris-icons";
import {userStore} from "@/stores/userStore.ts";

const Home = () => {
  const guideShow = userStore(state => state.guideShow);
  return (
    <s-page>
      <s-heading><Text fontWeight="bold" as="h1" variant="headingLg">👋 Hi, Welcome to goodcare protection
      </Text></s-heading>
      <s-stack gap="base">
        <s-box>
          <FulfilledOrders />
        </s-box>

        {guideShow && <NewPersonBox />}

        <s-section padding="none" heading="">
          <div className="p-4">
            <Text as="h3" variant="headingSm">
              Frequently asked questions
            </Text>
          </div>
          <div className="border-b-1 border-gray-200" />

          <div className="p-4 flex flex-col gap-4">
            <BlockStack gap="200">
              <InlineStack blockAlign="center" align="space-between">
                <Text as="h3" variant="headingXs">
                  Do you really offer 24/7 live chat?
                </Text>
                <div>
                  <Icon tone="base" source={ChevronUpIcon} />
                </div>
              </InlineStack>

              <Collapsible open id="test">

                <Text as="p" variant="bodyMd" tone="subdued">Offering World Class Support is something we strive for
                  here at Loloyal. We have real people to help 24*7 on live chat ready to help. No bots like other
                  apps.</Text>
              </Collapsible>
            </BlockStack>

            <BlockStack gap="200">
              <InlineStack blockAlign="center" align="space-between">
                <Text as="h3" variant="headingXs">
                  Do you offer help to move from other loyalty apps like S**e or move an existing program to Loloyal??
                </Text>
                <div>
                  <Icon tone="base" source={ChevronUpIcon} />
                </div>
              </InlineStack>
              <Collapsible open id="test1">

                <Text as="p" variant="bodyMd" tone="subdued">Absolutely, just get in touch on live chat and we can help
                  you instantly. Or you can export your customers&#39; existing point balances and easily import them
                  into
                  Loloyal with a simple CSV file.</Text>
              </Collapsible>
            </BlockStack>

            <BlockStack gap="200">
              <InlineStack blockAlign="center" align="space-between">
                <Text as="h3" variant="headingXs">
                  What are monthly order limits?
                </Text>
                <div>
                  <Icon tone="base" source={ChevronUpIcon} />
                </div>
              </InlineStack>
              <Collapsible open={false} id="test2">

                <Text as="p" variant="bodyMd" tone="subdued">Offering World Class Support is something we strive for
                  here at Loloyal. We have real people to help 24*7 on live chat ready to help. No bots like other
                  apps.</Text>
              </Collapsible>
            </BlockStack>

            <BlockStack gap="200">
              <InlineStack blockAlign="center" align="space-between">
                <Text as="h3" variant="headingXs">
                  Why are you so much cheaper than the competition?
                </Text>
                <div>
                  <Icon tone="base" source={ChevronUpIcon} />
                </div>
              </InlineStack>
              <Collapsible open={false} id="test3">

                <Text as="p" variant="bodyMd" tone="subdued">Offering World Class Support is something we strive for
                  here at Loloyal. We have real people to help 24*7 on live chat ready to help. No bots like other
                  apps.</Text>
              </Collapsible>
            </BlockStack>

            <BlockStack gap="200">
              <InlineStack blockAlign="center" align="space-between">
                <Text as="h3" variant="headingXs">
                  Can I change my plan later?
                </Text>
                <div>
                  <Icon tone="base" source={ChevronUpIcon} />
                </div>
              </InlineStack>
              <Collapsible open={false} id="test4">
                <Text as="p" variant="bodyMd" tone="subdued">Offering World Class Support is something we strive for
                  here at Loloyal. We have real people to help 24*7 on live chat ready to help. No bots like other
                  apps.</Text>
              </Collapsible>
            </BlockStack>
          </div>
        </s-section>
        <s-section heading="Need any help?">
          <InlineStack blockAlign="center" align="space-between">
            <Box>
              Email
            </Box>
            <Box>
              Live chat
            </Box>
            <Box>
              Help center
            </Box>
          </InlineStack>
        </s-section>
      </s-stack>
    </s-page>
  );
};

export default Home;
