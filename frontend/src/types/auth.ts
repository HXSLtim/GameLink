export type LoginRequest = {
  username: string;
  password: string;
};

export type LoginResult = {
  token: string;
  expires_in?: number;
};

export type CurrentUser = {
  id: string;
  username: string;
  role: string;
};
