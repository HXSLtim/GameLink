import { http } from '../api/http';
import { PageQuery, PageResult } from '../types/api';
import { Order } from '../types/order';

export const orderService = {
  list(query?: PageQuery) {
    return http.get<PageResult<Order>>('/orders', query);
  },
};
