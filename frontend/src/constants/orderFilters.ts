
// 时间范围选项
export const TIME_RANGES = {
  TODAY: 'today',
  LAST_7_DAYS: 'last_7_days',
  LAST_14_DAYS: 'last_14_days',
  LAST_30_DAYS: 'last_30_days',
  LAST_60_DAYS: 'last_60_days',
  LAST_90_DAYS: 'last_90_days',
  LAST_180_DAYS: 'last_180_days',
  LAST_360_DAYS: 'last_360_days',
  CUSTOM: 'custom',
} as const;


// 支付状态
export const PAYMENT_STATUS = {
  PAID: 'PAID',
  REFUNDED: 'REFUNDED',
  PARTIAL_REFUNDED: 'PARTIAL_REFUNDED',
} as const;


// 履约状态
export const FULFILLMENT_STATUS = {
  FULFILLED: 'FULFILLED',
  UNFULFILLED: 'UNFULFILLED',
  PARTIAL_FULFILLED: 'PARTIAL_FULFILLED',
} as const;


// 排序选项
export const SORT_OPTIONS = {
  ORDER_NUMBER: 'order_number',
  PAYMENT_DATE: 'payment_date',
  PROTECTION_FEE: 'protection_fee',
} as const;


export const SORT_DIRECTION = {
  ASC: 'asc',
  DESC: 'desc',
} as const;


// 完整的排序选项（包含方向）
export const FULL_SORT_OPTIONS = {
  ORDER_NUMBER_ASC: `${SORT_OPTIONS.ORDER_NUMBER} ${SORT_DIRECTION.ASC}`,
  ORDER_NUMBER_DESC: `${SORT_OPTIONS.ORDER_NUMBER} ${SORT_DIRECTION.DESC}`,
  PAYMENT_DATE_ASC: `${SORT_OPTIONS.PAYMENT_DATE} ${SORT_DIRECTION.ASC}`,
  PAYMENT_DATE_DESC: `${SORT_OPTIONS.PAYMENT_DATE} ${SORT_DIRECTION.DESC}`,
  PROTECTION_FEE_ASC: `${SORT_OPTIONS.PROTECTION_FEE} ${SORT_DIRECTION.ASC}`,
  PROTECTION_FEE_DESC: `${SORT_OPTIONS.PROTECTION_FEE} ${SORT_DIRECTION.DESC}`,
} as const;