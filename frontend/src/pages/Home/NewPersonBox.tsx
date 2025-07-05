import {
    BlockStack, Box, Button, Card, Collapsible, InlineGrid, InlineStack, RadioButton, Text,
} from "@shopify/polaris";
import React, {useEffect, useState} from "react";
import intl from "@/lib/i18n";

const NewPersonBox = ({setNewPerson, data, toastFun, setPageData}) => {
    const list1 = [{
        title: "Enable app embed",
        dsc: `Please follow the 👉 help docs to enable app embed and complete setup.Rest assured that this step will not affect your store cart, this won't display the widget on your store cart until you've published it from our app.`,
        isAction: true,
        stepMask: "/image/in_card01.png",
        isOk: false,
        button: "Enable app embed",
        button2: '',
        vd: 0,
        task: 0,
    }, {
        title: "Creat Protection Plans",
        dsc: `Please follow the 👉 help docs to enable app embed and complete setup.Rest assured that this step will not affect your store cart, this won't display the widget on your store cart until you've published it from our app.`,
        isAction: true,
        stepMask: "/image/in_card01.png",
        isOk: false,
        button: "Enable app embed",
        button2: '',
        vd: 0,
        task: 0,
    }, {
        title: "Setup protection widget",
        dsc: `Please follow the 👉 help docs to enable app embed and complete setup.Rest assured that this step will not affect your store cart, this won't display the widget on your store cart until you've published it from our app.`,
        isAction: true,
        stepMask: "/image/in_card01.png",
        isOk: false,
        button: "Enable app embed",
        button2: '',
        vd: 0,
        task: 0,
    }, {
        title: "How does the protection work",
        dsc: `Please follow the 👉 help docs to enable app embed and complete setup.Rest assured that this step will not affect your store cart, this won't display the widget on your store cart until you've published it from our app.`,
        isAction: true,
        stepMask: "/image/in_card01.png",
        isOk: false,
        button: "Enable app embed",
        button2: '',
        vd: 0,
        task: 0,
    },

    ];


    const [newerStartData, setNewerStartData] = useState(list1);// 新人引导数据
    const [getStartClo, setGetStartClo] = useState(true);//引导是否关闭
    const [currentStart, setCurrentStart] = useState(0);//当前的引导步骤
    useEffect(() => {
        let isMount = false;
        return () => {
            isMount = true;
        };
         
    }, [data]);
    const handleItemToggle = (ind) => {
        setCurrentStart(ind);
    };


    return (<div>
            <s-box className="new_person_box">
                <Card>
                    <BlockStack gap="200">
                        <InlineGrid columns="1fr auto">
                <span
                    style={{
                        fontSize: "16px", fontWeight: 600, display: "inline-block", minHeight: "32px",
                    }}
                >
                  {intl.get("Get started with GoodCare")}
                </span>
                        </InlineGrid>
                        <s-box>
                            <InlineStack blockAlign="center">
                                <Text as="p" variant="bodyMd">
                                    {intl.get("Start running effortlessly with this personalized guide.",)}
                                </Text>
                            </InlineStack>
                        </s-box>
                        <Collapsible
                            open={getStartClo}
                            id="basic-collapsible"
                            transition={{
                                duration: "500ms", timingFunction: "ease-in-out",
                            }}
                            expandOnPrint
                        >
                            <s-box className="get_start_list">
                                {newerStartData.map((item, ind) => (<s-box
                                        key={ind}
                                        className={currentStart === ind ? "home_start_s home_start" : "home_start"}
                                    >
                                        <s-box
                                            className="home_start_p"
                                            style={{
                                                padding: "0 16px", display: "flex", alignItems: "center",

                                            }}
                                            onClick={() => handleItemToggle(ind)}
                                        >

                                            <RadioButton checked label="" />

                                            <Text as="p" variant="bodyMd">
                                                {item.title}
                                            </Text>
                                        </s-box>
                                        <Collapsible
                                            open={Number(currentStart) === Number(ind)}
                                            id="basic-collapsible"
                                            transition={{
                                                duration: "200ms", timingFunction: "ease-in-out",
                                            }}
                                            expandOnPrint
                                        >
                                            <s-box className="home_start_text">
                                                <s-box
                                                    className="home_start_text_title"
                                                >
                            <span
                                dangerouslySetInnerHTML={{
                                    __html: item.dsc,
                                }}
                             />
                                                    {item.button !== "" && (<p className="home_start_text_btn">
                                                            <Button
                                                                variant="primary"
                                                                onClick={() => console.log(999)}
                                                            >
                                                                {item.button}
                                                            </Button>
                                                        </p>)}
                                                </s-box>
                                                <s-box
                                                    className=""
                                                    onClick={() => console.log(8888)}
                                                >
                                                    <img src={item.stepMask} />
                                                </s-box>
                                            </s-box>
                                        </Collapsible>
                                    </s-box>))}
                            </s-box>
                        </Collapsible>
                    </BlockStack>
                </Card>
            </s-box>

        </div>);
};

export default NewPersonBox;
