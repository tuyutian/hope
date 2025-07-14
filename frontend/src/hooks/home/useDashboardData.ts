import { useState, useEffect } from 'react';
import { rqGetDashboard } from '@/api';

interface OrderStatistic {
  refund: number;
  orders: number;
  total: number;
  sales: number;
}

interface OrderStatisticsTableItem {
  date: string;
  sales: number;
  refund: number;
}

interface DashboardData {
  orderStaticsTable: OrderStatisticsTableItem[];
  orderStatics: OrderStatistic;
}

const initialData: DashboardData = {
  orderStaticsTable: [
    {
      date: "0",
      sales: 0,
      refund: 0,
    },
  ],
  orderStatics: {
    refund: 0,
    orders: 0,
    total: 0,
    sales: 0,
  }
};

export const useDashboardData = () => {
  const [data, setData] = useState<DashboardData>(initialData);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchData = async (period: string) => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await rqGetDashboard(period);
      
      if (response.code === 0) {
        setData({
          orderStaticsTable: response.data.order_statistics_table || [],
          orderStatics: response.data.order_statistics || initialData.orderStatics,
        });
      } else {
        setError('Failed to fetch dashboard data');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
      console.error('Error fetching dashboard data:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData('30');
  }, []);

  return {
    data,
    loading,
    error,
    fetchData,
    refetch: () => fetchData('30'),
  };
};
