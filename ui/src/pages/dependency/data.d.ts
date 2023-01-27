export type TableListItem = {
  id?: number;
  host_name: string;
  ip: string;
  config_name: string;
  heartbeat_time: string;
  create_time?: string;
  updated_time?: string;
  status: number;
};

export type TableListItemResponse = {
  list: TableListItem[];
  page_num: number;
  page_size: number;
  total: number;
};

export type TableListPagination = {
  total: number;
  pageSize: number;
  current: number;
};

// 发送给后段的数据
export type TableListParams = {
  hostName?: string;
  ip?: string;
  status?: number;
  pageSize?: number;
  currentPage?: number;
  filter?: Record<string, any[]>;
  sorter?: Record<string, any>;
};
