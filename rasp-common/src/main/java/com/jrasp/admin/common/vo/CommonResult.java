package com.jrasp.admin.common.vo;

import io.swagger.annotations.ApiModelProperty;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class CommonResult<T> {
    @ApiModelProperty("响应码")
    private int code;
    @ApiModelProperty("状态信息")
    private String msg;
    @ApiModelProperty("业务数据")
    private T data;
    public static <T> CommonResult<T> success(String msg,T data){
        CommonResult<T> result = new CommonResult<>();
        result.setCode(200);
        result.setMsg(msg);
        result.setData(data);
        return result;
    }
    public static <T> CommonResult<T> success(String msg){
        CommonResult<T> result = new CommonResult<>();
        result.setCode(200);
        result.setMsg(msg);
        return result;
    }
    public static <T> CommonResult<T> success(T data){
        String msg="获取数据成功！";
        CommonResult<T> result = new CommonResult<>();
        result.setCode(200);
        result.setMsg(msg);
        result.setData(data);
        return result;
    }
    public static <T> CommonResult<T> success(){
        CommonResult<T> result = new CommonResult<>();
        result.setCode(200);
        result.setMsg("操作成功！");
        return result;
    }
    public static <T> CommonResult<T> error(String msg){
        CommonResult<T> result = new CommonResult<>();
        result.setCode(0);
        result.setMsg(msg);
        return result;
    }
    public static <T> CommonResult<T> exception(Integer code,String message){
        CommonResult<T> result = new CommonResult<>();
        result.setCode(code);
        result.setMsg(message);
        return result;
    }
}
