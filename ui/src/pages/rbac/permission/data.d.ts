type PermissionListItem = {
  id?: number;
  parentId?: number;
  name?: string;
  icon: string;
  url?: string;
  priority?: number;
  hide_in_menu?: number;
  hide_children_in_menu?: number;
  parent_id?: number;
  children?: PermissionListItem[];
};
type PermissionTreeItem = {
  value?: number;
  title?: string;
  label?: string;
  children?: PermissionTreeItem[];
};
