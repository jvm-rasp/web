export type TableListItem = {
  id?: number;
  host_name: string;
  remote_ip: string;
  local_ip: string;
  attack_type: string;
  check_type: string;
  attack_time: string;
  handle_status: number;
  is_blocked: boolean;
  stack_trace: string;
  http_method: string;
  request_protocol: string;
  request_parameters: string;
  attack_parameters: string;
  header: string;
  body: string;
  request_uri: string;
  level: number;
  create_time?: string;
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
  host_name?: string;
  pageSize?: number;
  current?: number;
  filter?: Record<string, any[]>;
  sorter?: Record<string, any>;
};
