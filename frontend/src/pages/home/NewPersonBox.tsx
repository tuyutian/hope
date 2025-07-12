import {BlockStack, Box, Button, Collapsible, InlineStack, Text,} from "@shopify/polaris";
import React, {useEffect, useMemo, useState} from "react";
import intl from "@/lib/i18n";
import GuideImage from "@/assets/images/dashborad/guide.png";
import CheckCycle from "@/pages/home/CheckCycle.tsx";
import {useAuth} from "@/stores/context.ts";


const NewPersonBox = () => {
  const {user} = useAuth();
  const userGuideSteps = useMemo(() => {
    const userGuide = user.userGuide;
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
      key: "setting_protension",
      title: "Creat Protection Plans",
      dsc: `Please follow the 👉 help docs to enable app embed and complete setup.Rest assured that this step will not affect your store cart, this won't display the widget on your store cart until you've published it from our app.`,
      isAction: true,
      stepMask: GuideImage,
      isOk: userGuide.setting_protension,
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
  }, [user.userGuide]);
  const [getStartClo, setGetStartClo] = useState(true);//引导是否关闭
  const [currentStart, setCurrentStart] = useState(0);//当前的引导步骤

  useEffect(() => {
    Object.values(user.userGuide).map(function (value,index){
      if (!value) {
        setCurrentStart(index)
      }
    })
  }, []);

  const handleItemToggle = (ind: number) => {
    setCurrentStart(ind);
  };

  return <div>
    <Box>
      <s-section>
        <BlockStack gap="200">
          <s-grid gridTemplateColumns="1fr auto">
                <span
                  style={{
                    fontSize: "16px", fontWeight: 600, display: "inline-block", minHeight: "32px",
                  }}
                >
                  {intl.get("Get started with GoodCare") as string}
                </span>
          </s-grid>
          <Box>
            <InlineStack blockAlign="center">
              <Text as="p" variant="bodyMd">
                {intl.get("Start running effortlessly with this personalized guide.",) as string}
              </Text>
            </InlineStack>
          </Box>
          <Collapsible
            open={getStartClo}
            id="basic-collapsible"
            transition={{
              duration: "500ms", timingFunction: "ease-in-out",
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
                    <CheckCycle name={item.key} check={item.isOk} />
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
        </BlockStack>
      </s-section>
    </Box>
  </div>;
};

export default NewPersonBox;
