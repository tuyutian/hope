import BillingTableSkeleton from "@/pages/billing/components/BillingTableSkeleton.tsx";
import { BlockStack, Box, Card, EmptyState, IndexTable, InlineStack, Page, Pagination, Text } from "@shopify/polaris";
import { getMessageState } from "@/stores/messageStore.ts";
import { useState, useEffect } from "react";
import { useQuery } from "@tanstack/react-query";
import { useNavigate, useSearchParams } from "react-router";
import { ApiResponse, billingService } from "@/api";
import type { NonEmptyArray } from "@shopify/polaris/build/ts/src/types";
import { IndexTableHeading } from "@shopify/polaris/build/ts/src/components/IndexTable/IndexTable";
import { FilterParams } from "@/types/billing.ts";
import { formatUnixTimestampRange } from "@/utils/dateUtils";

type TableData = {
  list: {
    id: number;
    order_name: string;
    charged_at: string;
    total_price_amount: number;
    protectify_amount: number;
    commission_amount: number;
  }[];
  total: number;
};
export default function BillingDetail() {
  const toastMessage = getMessageState().toastMessage;
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const [filters, setFilters] = useState<FilterParams>({
    sort: "desc",
    page: 1,
    size: 10,
    minTime: 0,
    maxTime: 0,
  });

  // 从URL参数获取时间范围参数
  useEffect(() => {
    const urlMinTime = searchParams.get('minTime');
    const urlMaxTime = searchParams.get('maxTime');
    const minTime = urlMinTime ? Number(urlMinTime) : 0;
    const maxTime = urlMaxTime ? Number(urlMaxTime) : 0;

    if(minTime===0&&maxTime===0){
      navigate("/billing");
      return;
    }
    
    // 只有当URL参数存在时才更新filters，否则保持默认空值
    if (urlMinTime || urlMaxTime) {
      setFilters(prev => ({
        ...prev,
        minTime: minTime,
        maxTime: maxTime,
      }));
    }
  }, [searchParams]);

  const { data, error, isLoading, isFetching } = useQuery<ApiResponse, Error, TableData>({
    queryKey: ["billing-detail-table", filters],
    queryFn: async () => {
      try {
        return await billingService.getDetailData(filters);
      } catch (error) {
        console.log(error);
        return {
          code: 0,
          data: {
            list: [],
            total: 0,
          },
          message: "暂无数据",
        };
      }
    },
    select: res => (res?.code === 0 ? res.data : { list: [], total: 0 }),
    staleTime: 5 * 60 * 1000,
    gcTime: 10 * 60 * 1000,
  });

  if (error || !data) {
    if (error) {
      toastMessage(error.message, 5000, true);
    }
    return (
      <EmptyState
        heading="Request data failed"
        image="https://cdn.shopify.com/s/files/1/0262/4071/2726/files/emptystate-files.png"
        fullWidth
      >
        <p>Please contact us to fix this error!</p>
      </EmptyState>
    );
  }

  if (isLoading) {
    return <BillingTableSkeleton />;
  }

  // 处理分页逻辑
  const handlePreviousPage = () => {
    if (filters.page > 1) {
      setFilters(prev => ({ ...prev, page: prev.page - 1 }));
    }
  };

  const handleNextPage = () => {
    const totalPages = Math.ceil((data?.total || 0) / filters.size);
    if (filters.page < totalPages) {
      setFilters(prev => ({ ...prev, page: prev.page + 1 }));
    }
  };

  // 计算分页状态
  const totalPages = Math.ceil((data?.total || 0) / filters.size);
  const hasPrevious = filters.page > 1;
  const hasNext = filters.page < totalPages;
  const headers: NonEmptyArray<IndexTableHeading> = [
    { id: "order", title: "Order" },
    { id: "payment_date", title: "Payment date" },
    { id: "sales", title: "Sales" },
    { id: "protection_fee", title: "Protection fee" },
    { id: "protection_billing", title: "Protection billing" },
  ];
  const resourceName = {
    singular: "record",
    plural: "records",
  };

  const rowMarkup = data.list.map(
    ({ id, order_name, charged_at, total_price_amount, protectify_amount, commission_amount }, index) => (
      <IndexTable.Row id={String(id)} key={id} position={index}>
        <IndexTable.Cell>
          <Text variant="bodyMd" fontWeight="bold" as="span">
            {order_name}
          </Text>
        </IndexTable.Cell>
        <IndexTable.Cell>{charged_at}</IndexTable.Cell>
        <IndexTable.Cell>{total_price_amount}</IndexTable.Cell>
        <IndexTable.Cell>
          <Text as="span" alignment="end" numeric>
            {protectify_amount}
          </Text>
        </IndexTable.Cell>
        <IndexTable.Cell>
          <Text as="span" alignment="end" numeric>
            {commission_amount}
          </Text>
        </IndexTable.Cell>
      </IndexTable.Row>
    )
  );

  return (
    <Page
      backAction={{
        url: "/billing",
      }}
      title="Protection billing"
      subtitle={`Billing cycle: ${formatUnixTimestampRange(filters.minTime, filters.maxTime)}`}
    >
      {" "}
      <Card padding="0">
        <Box padding="400">
          <Text as="h2" variant="headingSm">
            Past Protection Billing
          </Text>
        </Box>
        <Box paddingBlockEnd="400">
          <BlockStack gap="300">
            <IndexTable
              resourceName={resourceName}
              itemCount={data.list.length}
              loading={isFetching}
              headings={headers}
              selectable={false}
              emptyState={
                <EmptyState
                  heading="You currently have no new policy"
                  image="https://cdn.shopify.com/s/files/1/0262/4071/2726/files/emptystate-files.png"
                  fullWidth
                >
                  <p>Set up your protectify plugin, embed it in your store and start increasing your revenue.</p>
                </EmptyState>
              }
            >
              {rowMarkup}
            </IndexTable>

            {/* 只有在有数据时才显示分页 */}
            {(data?.total || 0) > 0 && (
              <InlineStack align="center">
                <Pagination
                  hasPrevious={hasPrevious}
                  hasNext={hasNext}
                  onPrevious={handlePreviousPage}
                  onNext={handleNextPage}
                  label={`${filters.page} / ${totalPages} `}
                />
              </InlineStack>
            )}
          </BlockStack>
        </Box>
      </Card>
    </Page>
  );
}
