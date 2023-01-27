import React, { useEffect, useState } from 'react';
import { TreeSelect } from 'antd';
import { permissionIndex } from '@/pages/rbac/permission/service';

interface PermissionInputProps {
  value?: number;
  onChange?: (value?: number) => void;
}
const PermissionTreeSelect: React.FC<PermissionInputProps> = ({ value, onChange }) => {
  const [permissionId, setPermissionId] = useState<number>(0);
  const [data, dataHandler] = useState<PermissionTreeItem[]>([]);
  const triggerChange = (changedValue: number) => {
    onChange?.(changedValue);
  };
  const fetchTree = async () => {
    const response = await permissionIndex({});
    const list = response.data || [];
    const first = {
      title: '顶级权限',
      value: 0,
      children: [],
    };
    const treeData = list.map((item) => {
      const childrenData = item.children || [];
      const children = childrenData.map((child) => {
        return {
          title: child.name,
          value: child.id,
        };
      });
      return {
        title: item.name,
        value: item.id,
        children,
      };
    });
    treeData.unshift(first);
    dataHandler(treeData);
  };
  useEffect(() => {
    fetchTree();
  }, []);
  return (
    <TreeSelect
      value={value || permissionId}
      placeholder={'请选择所属权限'}
      treeDefaultExpandAll={true}
      treeData={data}
      onChange={(selectedId) => {
        setPermissionId(selectedId);
        triggerChange(selectedId);
      }}
    />
  );
};
export default PermissionTreeSelect;
