import {
    Layout,
    LegacyCard,
    SkeletonBodyText, SkeletonDisplayText,
    SkeletonPage,
    TextContainer
} from '@shopify/polaris';
import React from 'react';

const SkeletonScreen = ()  => {
    return (
        <SkeletonPage primaryAction>
        <Layout>
            <Layout.Section>
                <LegacyCard sectioned>
                    <SkeletonBodyText />
                </LegacyCard>
                <LegacyCard sectioned>
                    <TextContainer>
                        <SkeletonDisplayText size="small" />
                        <SkeletonBodyText />
                    </TextContainer>
                </LegacyCard>
                <LegacyCard sectioned>
                    <TextContainer>
                        <SkeletonDisplayText size="small" />
                        <SkeletonBodyText />
                    </TextContainer>
                </LegacyCard>
            </Layout.Section>
            <Layout.Section variant="oneThird">
                <LegacyCard>
                    <LegacyCard.Section>
                        <TextContainer>
                            <SkeletonDisplayText size="small" />
                            <SkeletonBodyText lines={2} />
                        </TextContainer>
                    </LegacyCard.Section>
                    <LegacyCard.Section>
                        <SkeletonBodyText lines={1} />
                    </LegacyCard.Section>
                </LegacyCard>
                <LegacyCard subdued>
                    <LegacyCard.Section>
                        <TextContainer>
                            <SkeletonDisplayText size="small" />
                            <SkeletonBodyText lines={2} />
                        </TextContainer>
                    </LegacyCard.Section>
                    <LegacyCard.Section>
                        <SkeletonBodyText lines={2} />
                    </LegacyCard.Section>
                </LegacyCard>
            </Layout.Section>
        </Layout>
    </SkeletonPage>
    );
}

export default SkeletonScreen;
