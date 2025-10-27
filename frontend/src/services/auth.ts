import { http } from '../api/http';
import { CurrentUser, LoginRequest, LoginResult } from '../types/auth';

export const authService = {
  login(payload: LoginRequest) {
    return http.post<LoginResult>('/auth/login', payload);
  },
  me() {
    return http.get<CurrentUser>('/auth/me');
  },
  logout() {
    return http.post<{ ok: boolean }>('/auth/logout');
  },
};
