package com.jrasp.admin.common.vo;

import com.baomidou.mybatisplus.core.metadata.IPage;
import lombok.Data;

import java.util.List;

@Data
public class PageResult<T> {
    private Long pageNum;
    private Long pageSize;
    private Long total;
    private List<T> list;

    public static <T> PageResult<T> page(List<T> list,IPage<?> page){
        PageResult<T> result = new PageResult<T>();
        result.setTotal(page.getTotal());
        result.setPageNum(page.getCurrent());
        result.setPageSize(page.getSize());
        result.setList(list);
        return result;
    }
}
