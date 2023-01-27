export type TableListItem = {
  id?: number;
  username: string;
  avatar: string;
  mobile: string;
  created_at?: string;
  updated_at?: string;
  role_id: number;
  status: number;
  role?: {
    id: number;
    name: string;
  };
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

export type TableListData = {
  list: TableListItem[];
  pagination: Partial<TableListPagination>;
};

export type TableListParams = {
  username?: string;
  mobile?: string;
  status?: number;
  pageSize?: number;
  currentPage?: number;
  filter?: Record<string, any[]>;
  sorter?: Record<string, any>;
};
