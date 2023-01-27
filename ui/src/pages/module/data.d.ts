export type TableListItem = {
  id?: number;
  name: string;
  tag: string;
  type: string;
  version: string;
  url: string;
  hash: string;
  middleware_name: string;
  middleware_version: string;
  c: string;
  username: string;
  message: string;
  parameters: ParameterItem[];
  create_time?: string;
  update_time?: string;
  status: number;
};

export type ParameterItem = {
  key: string;
  value: string;
  info: string;
};

export type ModuleTypeItem = {
  label: string;
  value: number;
};

export type TableListItemResponse = {
  list: TableListItem[];
  page_num: number;
  pageSize: number;
  total: number;
};

export type TableListPagination = {
  total: number;
  pageSize: number;
  current: number;
};

export type TableListParams = {
  configName?: string;
  status?: number;
  pageSize?: number;
  currentPage?: number;
  filter?: Record<string, any[]>;
  sorter?: Record<string, any>;
};
