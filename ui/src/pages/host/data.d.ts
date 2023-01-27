export type TableListItem = {
  id?: number;
  host_name: string;
  ip: string;
  config_id: number;
  agent_mode: string;
  heartbeat_time: string;
  create_time?: string;
  os_type: string;
  total_mem: number;
  cpu_counts: number;
  free_disk: number;
  install_dir: string;
  exe_file_hash: string;
  build_date_time: string;
  build_git_branch: string;
  build_git_commit: string;
  nacos_info: string;
  agent_config_update_time: string;
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

// 发送给后端的数据
export type TableListParams = {
  hostName?: string;
  ip?: string;
  status?: number;
  pageSize?: number;
  current?: number;
  filter?: Record<string, any[]>;
  sorter?: Record<string, any>;
};

// 发送给后端的数据
export type UpdateConfigDetail = {
  configId: number;
  hostNames: string;
};
