import { http } from '../api/http';
import { PageQuery, PageResult } from '../types/api';
import { User } from '../types/user';

export const userService = {
  list(query?: PageQuery) {
    return http.get<PageResult<User>>('/users', query);
  },
};
