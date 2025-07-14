import {BlockStack, Box, Button, Collapsible, InlineStack, Popover, Text,} from "@shopify/polaris";
import React, { useEffect, useMemo, useState} from "react";
import intl from "@/lib/i18n";
import GuideImage from "@/assets/images/dashborad/guide.png";
import CheckCycle from "@/pages/home/CheckCycle.tsx";
import {ChevronDownIcon, ChevronUpIcon, MenuHorizontalIcon, XIcon} from "@shopify/polaris-icons";
import {getUserState} from "@/stores/userStore.ts";
import {UpdateUserSetting} from "@/api";
import {SettingCode} from "@/constants/settingCode.ts";
import {UserGuide} from "@/types/user.ts";


const NewPersonBox = () => {
  const {userGuide,closeGuide} = getUserState();
  const userGuideSteps = useMemo(() => {
    return [{
      key: "enabled",
      title: "Enable app embed",
      dsc: `Please follow the 👉 help docs to enable app embed and complete setup.Rest assured that this step will not affect your store cart, this won't display the widget on your store cart until you've published it from our app.`,
      isAction: true,
      stepMask: GuideImage,
      isOk: userGuide.enabled,
      button: "Enable app embed",
      button2: "",
      vd: 0,
      task: 0,
    }, {
      key: "setup_widget",
      title: "Setup protection widget",
      dsc: `Please follow the 👉 help docs to enable app embed and complete setup.Rest assured that this step will not affect your store cart, this won't display the widget on your store cart until you've published it from our app.`,
      isAction: true,
      stepMask: GuideImage,
      isOk: userGuide.setup_widget,
      button: "Enable app embed",
      button2: "",
      vd: 0,
      task: 0,
    }, {
      key: "how_work",
      title: "How does the protection work",
      dsc: `Please follow the 👉 help docs to enable app embed and complete setup.Rest assured that this step will not affect your store cart, this won't display the widget on your store cart until you've published it from our app.`,
      isAction: true,
      stepMask: GuideImage,
      isOk: userGuide.how_work,
      button: "Enable app embed",
      button2: "",
      vd: 0,
      task: 0,
    },
    ];
  }, [userGuide]);
  const [dismissShow, setDismissShow] = useState(false);//引导是否关闭
  const [colGuide, setColGuide] = useState(true);//引导是否关闭
  const [currentStart, setCurrentStart] = useState(0);//当前的引导步骤

  useEffect(() => {
    Object.values(userGuide).map(function (value, index) {
      if (!value) {
        setCurrentStart(index);
      }
    });
  }, []);

  const handleItemToggle = (ind: number) => {
    setCurrentStart(ind);
  };

  const toggleDismiss = function () {
    setDismissShow(!dismissShow);
  };
  const toggleGuide = function () {
    setColGuide(!colGuide);
  };
  const handleDismiss = async function () {
    await UpdateUserSetting(SettingCode.DashboardGuideHide,"1")
    closeGuide()
  };
  return <div>
    <Box>
      <s-section>
        <BlockStack>
          <InlineStack align="space-between" blockAlign="center">
            <Text as="h2" variant="headingMd">
              {intl.get("Get started with GoodCare") as string}
            </Text>

            <InlineStack gap="200" align="center" blockAlign="center">
              <div className="h-5">
                <Popover active={dismissShow}
                         activator={<Button variant="tertiary" icon={MenuHorizontalIcon} onClick={toggleDismiss} />}
                         onClose={toggleDismiss}>
                  <Popover.Pane>
                    <Button variant="tertiary" icon={XIcon} onClick={handleDismiss}>Dismiss</Button>
                  </Popover.Pane>
                </Popover>
              </div>
              <Button size="slim" variant="tertiary" icon={colGuide ? ChevronUpIcon : ChevronDownIcon}
                      onClick={toggleGuide} />
            </InlineStack>
          </InlineStack>
          <Collapsible id="dashboard-guide" open={colGuide} transition={{
            duration: "300ms", timingFunction: "ease-in-out"
          }}>
            <Box>
              <InlineStack blockAlign="center">
                <Text as="p" variant="bodyMd">
                  {intl.get("Start running effortlessly with this personalized guide.",) as string}
                </Text>
              </InlineStack>
            </Box>
            <Collapsible
              open
              id="basic-collapsible"
              transition={{
                duration: "300ms", timingFunction: "ease-in-out",
              }}
              expandOnPrint
            >
              <Box>
                {userGuideSteps.map((item, ind) => (<Box
                  key={item.key}
                  background={currentStart === ind ? "bg-fill-active" : "bg-fill"}
                  paddingInline="200" paddingBlock="400"
                  borderRadius="200"
                >
                  <div className={currentStart !== ind ? "cursor-pointer" : ""} onClick={() => {
                    currentStart !== ind && handleItemToggle(ind);
                  }}>
                    <InlineStack gap="400" wrap={false} blockAlign="start">
                      <CheckCycle name={item.key as keyof UserGuide} check={item.isOk} />
                      <Box>
                        <div className="Polaris-CalloutCard">
                          <div className="Polaris-CalloutCard__Content">
                            <div className={currentStart !== ind ? "mt-0.5" : "Polaris-CalloutCard__Title"}><h2
                              className="Polaris-Text--root Polaris-Text--headingSm">{item.title}</h2>
                            </div>
                            <Collapsible
                              open={Number(currentStart) === Number(ind)}
                              id="basic-collapsible"
                              transition={{
                                duration: "200ms", timingFunction: "ease-in-out",
                              }}
                              expandOnPrint
                            >
                              <BlockStack>
                              <span className="Polaris-Text--root Polaris-Text--bodyMd">
                              <p
                                dangerouslySetInnerHTML={{__html: item.dsc}} />
                          </span>
                                {item.button !== "" && (
                                  <div className="Polaris-CalloutCard__Buttons">
                                    <Button
                                      variant="primary"
                                      onClick={() => console.log(999)}
                                    >
                                      {item.button}
                                    </Button>
                                  </div>)}
                              </BlockStack>

                            </Collapsible>
                          </div>
                          {currentStart === ind && <img alt=""
                                                        src={item.stepMask}
                                                        className="Polaris-CalloutCard__Image" />}
                        </div>
                      </Box>
                    </InlineStack>
                  </div>
                </Box>))}
              </Box>
            </Collapsible>
          </Collapsible>
        </BlockStack>
      </s-section>
    </Box>
  </div>;
};

export default NewPersonBox;
