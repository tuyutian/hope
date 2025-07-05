import {
    TextField,
    IndexTable,
    LegacyCard,
    IndexFilters,
    useSetIndexFiltersMode,
    useIndexResourceState,
    Text,
    Badge,
    Page,
    Pagination,
    SkeletonBodyText
} from '@shopify/polaris';
import { useState, useEffect, useCallback, useMemo } from 'react';
import { rqGetOrderList } from "@/api";
import debounce from 'lodash.debounce';
import {formatTimestampToUSDate} from "@/utils/tools.ts";

export default function Order() {
    const itemStrings = ['All', 'Paid', 'Refund'];
    const tabs = itemStrings.map((item, index) => ({
        content: item,
        index,
        onAction: () => {},
        id: `${item}-${index}`,
        isLocked: index === 0,
    }));

    const [selected, setSelected] = useState(0);
    const { mode, setMode } = useSetIndexFiltersMode();

    const [orders, setOrders] = useState([]);
    const [totalCount, setTotalCount] = useState(0);
    const [currentPage, setCurrentPage] = useState(1);
    const itemsPerPage = 20;

    const [queryValue, setQueryValue] = useState('');
    const [debouncedQuery, setDebouncedQuery] = useState('');
    const [isLoading, setIsLoading] = useState(true);
    const [isTabLoading, setIsTabLoading] = useState(false);

    const handleFiltersQueryChange = useCallback((value) => {
        const trimmedValue = value.slice(0, 20);
        setQueryValue(trimmedValue);
        debouncedQueryUpdate(trimmedValue);
    }, []);

    const debouncedQueryUpdate = useMemo(
        () => debounce((val) => setDebouncedQuery(val), 500),
        []
    );

    useEffect(() => {
        return () => {
            debouncedQueryUpdate.cancel();
        };
    }, [debouncedQueryUpdate]);

    const handleQueryValueRemove = useCallback(() => {
        setQueryValue('');
        setDebouncedQuery('');
    }, []);

    const handleFiltersClearAll = useCallback(() => {
        handleQueryValueRemove();
        setSelected(0);
    }, [handleQueryValueRemove]);

    const getBadgeStatus = (status) => {
        switch (status) {
            case 'PAID':
                return 'success';
            case 'PARTIALLY_PAID':
                return 'attention';
            case 'PARTIALLY_REFUNDED':
            case 'REFUNDED':
                return 'warning';
            case 'UNPAID':
            default:
                return 'critical';
        }
    };

    const getOrderData = async () => {
        try {
            const res = await rqGetOrderList({
                page: currentPage,
                page_size: itemsPerPage,
                type: itemStrings[selected],
                query: debouncedQuery,
            });

            if (res?.code === 200) {
                const list = res.data.list.map(item => ({
                    id: item.id,
                    order_name: item.order_name,
                    order: <Text as="span" variant="bodyMd" fontWeight="semibold">{item.order_name}</Text>,
                    date: item.order_completion_at > 0 ? formatTimestampToUSDate(item.order_completion_at) : "-",
                    item: item.sku_num,
                    paymentStatus: <Badge tone={getBadgeStatus(item.financial_status)}>
                        {item.financial_status.replace(/_/g, ' ')}
                    </Badge>,
                    total: `${item.total_price_amount} ${item.currency}`,
                    protectionFee: `${item.insurance_amount} ${item.currency}`,
                }));

                setOrders(list);
                setTotalCount(res.data.count);
            }
        } finally {
            setIsLoading(false);
            setIsTabLoading(false);
        }
    };

    useEffect(() => {
        setIsLoading(true);
        getOrderData();
    }, []);

    useEffect(() => {
        if (isTabLoading) {
            getOrderData();
        }
    }, [isTabLoading]);

    useEffect(() => {
        if (currentPage !== 1) {
            setIsLoading(true);
            getOrderData();
        }
    }, [currentPage]);

    const onHandleCancel = () => {
        setQueryValue('');
        setDebouncedQuery('');
        getOrderData();
    };

    const primaryAction = {
        type: 'save',
        onAction: async () => {
            setIsLoading(true);
            setCurrentPage(1);
            await getOrderData();
        },
        disabled: false,
        loading: false,
    };

    const resourceName = {
        singular: 'order',
        plural: 'orders',
    };

    const { selectedResources, allResourcesSelected, handleSelectionChange } = useIndexResourceState(orders);

    const rowMarkup = isLoading ? (
        Array(itemsPerPage).fill(0).map((_, index) => (
            <IndexTable.Row key={`skeleton-${index}`}>
                {Array(6).fill(0).map((_, cellIndex) => (
                    <IndexTable.Cell key={`skeleton-cell-${cellIndex}`}>
                        <SkeletonBodyText lines={1} />
                    </IndexTable.Cell>
                ))}
            </IndexTable.Row>
        ))
    ) : (
        orders.map(
            ({ id, order_name, date, item, paymentStatus, total, protectionFee }, index) => (
                <IndexTable.Row
                    id={id}
                    key={id}
                    selected={selectedResources.includes(id)}
                    position={index}
                >
                    <IndexTable.Cell>
                        <Text as="span" alignment="center">{order_name}</Text>
                    </IndexTable.Cell>
                    <IndexTable.Cell>
                        <Text as="span" alignment="center">{date}</Text>
                    </IndexTable.Cell>
                    <IndexTable.Cell>
                        <Text as="span" alignment="center" numeric>{item}</Text>
                    </IndexTable.Cell>
                    <IndexTable.Cell>
                        <Text as="span" alignment="center" numeric>{paymentStatus}</Text>
                    </IndexTable.Cell>
                    <IndexTable.Cell alignment="center">
                        <Text as="span" alignment="center" numeric>{total}</Text>
                    </IndexTable.Cell>
                    <IndexTable.Cell>
                        <Text as="span" alignment="center" numeric>{protectionFee}</Text>
                    </IndexTable.Cell>
                </IndexTable.Row>
            )
        )
    );

    const totalPages = Math.ceil(totalCount / itemsPerPage);

    if (isLoading) {
        return (
            <Page title="Protection Orders">
                <LegacyCard>
                    <IndexTable
                        resourceName={resourceName}
                        itemCount={itemsPerPage}
                        headings={[
                            { title: 'Order', alignment: 'center' },
                            { title: 'Date', alignment: 'center' },
                            { title: 'Items', alignment: 'center' },
                            { title: 'Payment status', alignment: 'center' },
                            { title: 'Total', alignment: 'center' },
                            { title: 'Protection Fee', alignment: 'center' },
                        ]}
                    >
                        {rowMarkup}
                    </IndexTable>
                </LegacyCard>
            </Page>
        );
    }

    return (
        <Page title="Protection Orders">
            <LegacyCard>
                <IndexFilters
                    primaryAction={primaryAction}
                    onClearAll={handleFiltersClearAll}
                    mode={mode}
                    setMode={setMode}
                    queryValue={queryValue}
                    queryPlaceholder="Search orders"
                    onQueryChange={handleFiltersQueryChange}
                    onQueryClear={handleQueryValueRemove}
                    tabs={tabs}
                    selected={selected}
                    onSelect={(index) => {
                        setIsTabLoading(true);
                        setCurrentPage(1);
                        setSelected(index);
                    }}
                    cancelAction={{
                        onAction: onHandleCancel,
                        disabled: false,
                        loading: false,
                    }}
                    hideFilters
                    hideQueryField={false}
                    disabled={isTabLoading}
                    canCreateNewView={false}
                />
                <IndexTable
                    resourceName={resourceName}
                    itemCount={orders.length}
                    selectedItemsCount={allResourcesSelected ? 'All' : selectedResources.length}
                    onSelectionChange={handleSelectionChange}
                    headings={[
                        { title: 'Order', alignment: 'center' },
                        { title: 'Date', alignment: 'center' },
                        { title: 'Items', alignment: 'center' },
                        { title: 'Payment status', alignment: 'center' },
                        { title: 'Total', alignment: 'center' },
                        { title: 'Protection Fee', alignment: 'center' },
                    ]}
                    loading={false} // Explicitly set to false to hide loading banner
                    hasMoreItems
                >
                    {rowMarkup}
                </IndexTable>
                <div style={{ padding: '1rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <Text variant="bodySm" as="p">
                        Total: {totalCount} orders | Page {currentPage} of {totalPages}
                    </Text>
                    <Pagination
                        hasPrevious={currentPage > 1}
                        onPrevious={() => setCurrentPage(currentPage - 1)}
                        hasNext={currentPage < totalPages}
                        onNext={() => setCurrentPage(currentPage + 1)}
                        disabled={isTabLoading}
                    />
                </div>
            </LegacyCard>
        </Page>
    );
}