export function getPermissionOpenKeys(list: PermissionListItem[], ids: string[]): string[] {
  list.forEach((item) => {
    if (item.id) {
      ids.push(item.id.toString());
    }
    if (item.children) {
      getPermissionOpenKeys(item.children, ids);
    }
  });
  return ids;
}
