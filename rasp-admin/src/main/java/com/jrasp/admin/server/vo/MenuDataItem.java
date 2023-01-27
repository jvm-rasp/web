package com.jrasp.admin.server.vo;

import lombok.Data;

import java.util.List;

@Data
public class MenuDataItem {
    private Integer id;
    private String icon;
    private String name;
    private String key;
    private String path;
    private Integer parentId;
    private List<MenuDataItem> children;
    private Boolean hideChildrenInMenu;
    private Boolean hideInMenu;
    private List<String> parentKeys;
    private Integer priority;
}
