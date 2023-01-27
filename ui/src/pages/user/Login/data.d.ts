type LoginParams = {
  username?: string;
  password?: string;
};
type LoginResult = {
  code?: number;
  msg?: string;
  data?: {
    token?: string;
    url?: string;
  };
};
