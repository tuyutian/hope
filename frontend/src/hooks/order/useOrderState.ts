import { useState, useCallback } from 'react';

interface UseOrderStateReturn {
  isTabLoading: boolean;
  setIsTabLoading: (loading: boolean) => void;
  startTabLoading: () => void;
  stopTabLoading: () => void;
}

export const useOrderState = (): UseOrderStateReturn => {
  const [isTabLoading, setIsTabLoading] = useState(false);

  const startTabLoading = useCallback(() => {
    setIsTabLoading(true);
  }, []);

  const stopTabLoading = useCallback(() => {
    // 添加延迟以提供视觉反馈
    setTimeout(() => {
      setIsTabLoading(false);
    }, 300);
  }, []);

  return {
    isTabLoading,
    setIsTabLoading,
    startTabLoading,
    stopTabLoading,
  };
};