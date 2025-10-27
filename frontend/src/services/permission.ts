import { http } from '../api/http';
import { PageQuery, PageResult } from '../types/api';

export type Permission = {
  id: string;
  name: string;
  desc?: string;
};

export const permissionService = {
  list(query?: PageQuery) {
    return http.get<PageResult<Permission>>('/permissions', query as any);
  },
};
