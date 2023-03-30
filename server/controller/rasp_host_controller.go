package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"server/common"
	"server/repository"
	"server/response"
	"server/socket"
	"server/vo"
)

type IRaspHostController interface {
	GetRaspHosts(c *gin.Context)
	BatchDeleteHostByIds(c *gin.Context)
	PushConfig(c *gin.Context)
}

const DEFAULT_CONFIG = "{\"agentMode\":\"dynamic\",\"version\":\"1.1.1\",\"configId\":1,\"moduleAutoUpdate\":false,\"agentConfigs\":{\"check_disable\":false,\"redirect_url\":\"https://www.jrasp.com/block.html\",\"block_status_code\":302,\"json_block_content\":\"{\\\"error\\\":true, \\\"reason\\\": \\\"Request blocked by JRASP (https://www.jrasp.com)\\\"}\",\"xml_block_content\":\"<?xml version=\\\"1.0\\\"?><doc><error>true</error><reason>Request blocked by JRASP</reason></doc>\",\"html_block_content\":\"</script><script>location.href=\\\"https://www.jrasp.com/block.html\\\"</script>\"},\"moduleConfigs\":[{\"moduleName\":\"deserialization-algorithm\",\"parameters\":{\"json_black_list_action\":0,\"json_white_class_list\":[],\"json_black_class_list\":[\"org.apache.commons.collections.Transformer\",\"java.lang.Thread\",\"java.net.Socket\",\"java.net.URL\",\"java.net.InetAddress\",\"java.lang.Class\",\"oracle.jdbc.rowset.OracleJDBCRowSet\",\"oracle.jdbc.connector.OracleManagedConnectionFactory\",\"java.lang.UNIXProcess\",\"java.lang.AutoCloseable\",\"java.lang.Runnable\",\"java.util.EventListener\",\"java.io.PrintWriter\",\"java.io.FileInputStream\",\"java.io.FileOutputStream\",\"java.util.PriorityQueue\"],\"json_black_package_list\":[\"org.apache.commons.collections.functors\",\"org.apache.commons.collections4.functors\",\"org.apache.commons.collections4.comparators\",\"org.python.core\",\"org.apache.tomcat\",\"org.apache.xalan\",\"javax.xml\",\"org.springframework\",\"org.apache.commons.beanutils\",\"org.codehaus.groovy.runtime\",\"javax.net\",\"com.mchange\",\"org.apache.wicket.util\",\"java.util.jar\",\"org.mozilla.javascript\",\"java.rmi\",\"java.util.prefs\",\"com.sun\",\"java.util.logging\",\"org.apache.bcel\",\"org.apache.commons.fileupload\",\"org.hibernate\",\"org.jboss\",\"org.apache.myfaces.context.servlet\",\"org.apache.ibatis.datasource\",\"org.apache.log4j\",\"org.apache.logging\",\"org.apache.commons.dbcp\",\"com.ibatis.sqlmap.engine.datasource\",\"javassist\",\"oracle.net\",\"com.alibaba.fastjson.annotation\",\"com.zaxxer.hikari\",\"ch.qos.logback\",\"com.mysql.cj.jdbc.admin\",\"org.apache.ibatis.parsing\",\"org.apache.ibatis.executor\",\"com.caucho\"],\"ois_black_list_action\":0,\"ois_black_class_list\":[\"org.codehaus.groovy.runtime.ConvertedClosure\",\"org.codehaus.groovy.runtime.ConversionHandler\",\"org.codehaus.groovy.runtime.MethodClosure\",\"org.springframework.transaction.support.AbstractPlatformTransactionManager\",\"java.rmi.server.UnicastRemoteObject\",\"java.rmi.server.RemoteObjectInvocationHandler\",\"com.bea.core.repackaged.springframework.transaction.support.AbstractPlatformTransactionManager\",\"java.rmi.server.RemoteObject\",\"com.tangosol.coherence.rest.util.extractor.MvelExtractor\",\"java.lang.Runtime\",\"oracle.eclipselink.coherence.integrated.internal.cache.LockVersionExtractor\",\"org.eclipse.persistence.internal.descriptors.MethodAttributeAccessor\",\"org.eclipse.persistence.internal.descriptors.InstanceVariableAttributeAccessor\",\"org.apache.commons.fileupload.disk.DiskFileItem\",\"oracle.jdbc.pool.OraclePooledConnection\",\"com.tangosol.util.extractor.ReflectionExtractor\",\"com.tangosol.internal.util.SimpleBinaryEntry\",\"com.tangosol.coherence.component.util.daemon.queueProcessor.service.grid.partitionedService.PartitionedCache$Storage$BinaryEntry\",\"com.sun.rowset.JdbcRowSetImpl\",\"org.eclipse.persistence.internal.indirection.ProxyIndirectionHandler\",\"bsh.XThis\",\"bsh.Interpreter\",\"com.mchange.v2.c3p0.PoolBackedDataSource\",\"com.mchange.v2.c3p0.impl.PoolBackedDataSourceBase\",\"org.apache.commons.beanutils.BeanComparator\",\"java.util.PriorityQueue\",\"java.lang.reflect.Proxy\",\"clojure.lang.PersistentArrayMap\",\"org.apache.commons.io.output.DeferredFileOutputStream\",\"org.apache.commons.io.output.ThresholdingOutputStream\",\"org.apache.wicket.util.upload.DiskFileItem\",\"org.apache.wicket.util.io.DeferredFileOutputStream\",\"org.apache.wicket.util.io.ThresholdingOutputStream\",\"com.sun.org.apache.bcel.internal.util.ClassLoader\",\"com.sun.syndication.feed.impl.ObjectBean\",\"org.springframework.beans.factory.ObjectFactory\",\"org.springframework.aop.framework.AdvisedSupport\",\"org.springframework.aop.target.SingletonTargetSource\",\"com.vaadin.data.util.NestedMethodProperty\",\"com.vaadin.data.util.PropertysetItem\",\"javax.management.BadAttributeValueExpException\",\"org.apache.myfaces.context.servlet.FacesContextImpl\",\"org.apache.myfaces.context.servlet.FacesContextImplBase\"],\"ois_black_package_list\":[\"org.apache.commons.collections.functors\",\"org.apache.commons.collections4.functors\",\"com.sun.org.apache.xalan.internal.xsltc.trax\",\"org.apache.xalan.xsltc.trax\",\"javassist\",\"java.rmi.activation\",\"sun.rmi.server\",\"com.bea.core.repackaged.springframework.aop.aspectj\",\"com.bea.core.repackaged.springframework.beans.factory.support\",\"org.python.core\",\"com.bea.core.repackaged.aspectj.weaver.tools.cache\",\"com.bea.core.repackaged.aspectj.weaver.tools\",\"com.bea.core.repackaged.aspectj.weaver.reflect\",\"com.bea.core.repackaged.aspectj.weaver\",\"com.oracle.wls.shaded.org.apache.xalan.xsltc.trax\",\"oracle.eclipselink.coherence.integrated.internal.querying\",\"oracle.eclipselink.coherence.integrated.internal.cache\",\"javax.swing.plaf.synth\",\"javax.swing.plaf.metal\",\"com.tangosol.internal.util.invoke\",\"com.tangosol.internal.util.invoke.lambda\",\"com.tangosol.util.extractor\",\"com.tangosol.coherence.rest.util.extractor\",\"com.tangosol.coherence.rest.util\",\"com.tangosol.coherence.component.application.console\",\"org.mozilla.javascript\",\"org.apache.myfaces.el\",\"org.apache.myfaces.view.facelets.el\"],\"xml_black_list_action\":0,\"xml_black_class_list\":[\"java.io.PrintWriter\",\"java.io.FileInputStream\",\"java.io.FileOutputStream\",\"java.util.PriorityQueue\",\"javax.sql.rowset.BaseRowSet\",\"javax.activation.DataSource\",\"java.nio.channels.Channel\",\"java.io.InputStream\",\"java.lang.ProcessBuilder\",\"java.lang.Runtime\",\"javafx.collections.ObservableList\",\"java.beans.EventHandler\",\"sun.swing.SwingLazyValue\",\"java.io.File\"],\"xml_black_package_list\":[\"sun.reflect\",\"sun.tracing\",\"com.sun.corba\",\"javax.crypto\",\"jdk.nashorn.internal\",\"sun.awt.datatransfer\",\"com.sun.tools\",\"javax.imageio\",\"com.sun.rowset\"],\"xml_black_key_list\":[\".jndi.\",\".rmi.\",\".bcel.\",\".xsltc.trax.TemplatesImpl\",\".ws.client.sei.\",\"$URLData\",\"$LazyIterator\",\"$GetterSetterReflection\",\"$PrivilegedGetter\",\"$ProxyLazyValue\",\"$ServiceNameIterator\"]}},{\"moduleName\":\"deserialization-hook\",\"parameters\":{\"disable\":false}},{\"moduleName\":\"expression-algorithm\",\"parameters\":{\"ognl_min_length\":30,\"ognl_max_limit_length\":200,\"ognl_black_list_action\":0,\"ognl_max_limit_length_action\":0,\"ognl_black_list\":[\"java.lang.Runtime\",\"java.lang.Class\",\"java.lang.ClassLoader\",\"java.lang.System\",\"java.lang.ProcessBuilder\",\"java.lang.Object\",\"java.lang.Shutdown\",\"ognl.OgnlContext\",\"ognl.TypeConverter\",\"ognl.MemberAccess\",\"_memberAccess\",\"ognl.ClassResolver\",\"java.io.File\",\"javax.script.ScriptEngineManager\",\"excludedClasses\",\"excludedPackageNamePatterns\",\"excludedPackageNames\",\"com.opensymphony.xwork2.ActionContext\"],\"spel_min_length\":30,\"spel_max_limit_length\":400,\"spel_black_list_action\":0,\"spel_max_limit_length_action\":0,\"spel_black_list\":[\"java.lang.Runtime\",\"java.lang.ProcessBuilder\",\"javax.script.ScriptEngineManager\",\"java.lang.System\",\"org.springframework.cglib.core.ReflectUtils\",\"java.io.File\",\"javax.management.remote.rmi.RMIConnector\"]}},{\"moduleName\":\"expression-hook\",\"parameters\":{\"disable\":false}},{\"moduleName\":\"file-algorithm\",\"parameters\":{\"file_delete_action\":0,\"file_list_action\":0,\"file_read_action\":0,\"file_upload_action\":0,\"danger_dir_list\":[\"/\",\"/home\",\"/etc\",\"/usr\",\"/usr/local\",\"/var/log\",\"/proc\",\"/sys\",\"/root\",\"C:\\\\\",           \"D:\\\\\",           \"E:\\\\\"         ],         \"file_upload_black_list\": [           \".jsp\",           \".asp\",           \".phar\",           \".phtml\",           \".sh\",           \".py\",           \".pl\",           \".rb\"         ],         \"travel_str\": [           \"../\",           \"..\\\\\"         ]       }     },     {       \"moduleName\": \"file-hook\",       \"parameters\": {         \"disable\": false       }     },     {       \"moduleName\": \"http-algorithm\",       \"parameters\": {         \"disable\": false       }     },     {       \"moduleName\": \"http-hook\",       \"parameters\": {         \"disable\": false       }     },     {       \"moduleName\": \"jndi-hook\",       \"parameters\": {         \"disable\": false,         \"jndi_black_list_action\": 0,         \"danger_protocol\": [           \"ldap://\",           \"rmi://\"         ]       }     },     {       \"moduleName\": \"mysql-hook\",       \"parameters\": {         \"disable\": false       }     },     {       \"moduleName\": \"rce-algorithm\",       \"parameters\": {         \"rce_action\": 0,         \"rce_white_list\": [         ],         \"rce_danger_stack_list\": [           \"com.thoughtworks.xstream.XStream.unmarshal\",           \"java.beans.XMLDecoder.readObject\",           \"java.io.ObjectInputStream.readObject\",           \"org.apache.dubbo.common.serialize.hessian2.Hessian2ObjectInput.readObject\",           \"com.alibaba.fastjson.JSON.parse\",           \"com.fasterxml.jackson.databind.ObjectMapper.readValue\",           \"payload.execCommand\",           \"net.rebeyond.behinder\",           \"org.springframework.expression.spel.support.ReflectiveMethodExecutor.execute\",           \"freemarker.template.utility.Execute.exec\",           \"freemarker.core.Expression.eval\",           \"bsh.Reflect.invokeMethod\",           \"org.jboss.el.util.ReflectionUtil.invokeMethod\",           \"org.codehaus.groovy.runtime.ProcessGroovyMethods.execute\",           \"org.codehaus.groovy.runtime.callsite.AbstractCallSite.call\",           \"ScriptFunction.invoke\",           \"com.caucho.hessian.io.HessianInput.readObject\",           \"org.apache.velocity.runtime.parser.node.ASTMethod.execute\",           \"org.apache.commons.jexl3.internal.Interpreter.call\",           \"javax.script.AbstractScriptEngine.eval\",           \"javax.el.ELProcessor.getValue\",           \"ognl.OgnlRuntime.invokeMethod\",           \"javax.naming.InitialContext.lookup\",           \"org.mvel2.MVEL.executeExpression\",           \"org.mvel.MVEL.executeExpression\",           \"ysoserial.Pwner\",           \"org.yaml.snakeyaml.Yaml.load\",           \"org.mozilla.javascript.Context.evaluateString\",           \"command.Exec.equals\",           \"java.lang.ref.Finalizer.runFinalizer\",           \"java.sql.DriverManager.getConnection\"         ],         \"rce_block_list\": [           \"curl\",           \"wget\",           \"echo\",           \"touch\",           \"gawk\",           \"telnet\",           \"xterm\",           \"perl\",           \"python\",           \"python3\",           \"ruby\",           \"lua\",           \"whoami\",           \"php\",           \"pwd\",           \"ifconfig\",           \"alias\",           \"export\"         ]       }     },     {       \"moduleName\": \"rce-hook\",       \"parameters\": {         \"disable\": false       }     },     {       \"moduleName\": \"sql-algorithm\",       \"parameters\": {         \"mysql_block_action\": 0,         \"mysql_min_limit_length\": 16,         \"mysql_max_limit_length\": 65535,         \"mysql_white_list\": []       }     },     {       \"moduleName\": \"ssrf-algorithm\",       \"parameters\": {         \"disable\": false       }     },     {       \"moduleName\": \"ssrf-hook\",       \"parameters\": {         \"disable\": false       }     },     {       \"moduleName\": \"xxe-hook\",       \"parameters\": {         \"disable\": false       }     }   ] }"

type RaspHostController struct {
	RaspHostRepository   repository.IRaspHostRepository
	RaspConfigRepository repository.IRaspConfigRepository
}

func NewRaspHostController() IRaspHostController {
	raspHostRepository := repository.NewRaspHostRepository()
	raspConfigRepository := repository.NewRaspConfigRepository()
	raspHostController := RaspHostController{
		RaspHostRepository:   raspHostRepository,
		RaspConfigRepository: raspConfigRepository,
	}
	return raspHostController
}

func (h RaspHostController) GetRaspHosts(c *gin.Context) {
	var req vo.RaspHostListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}
	// 获取
	raspHosts, total, err := h.RaspHostRepository.GetRaspHosts(&req)
	if err != nil {
		response.Fail(c, nil, "获取实例列表失败")
		return
	}
	response.Success(c, gin.H{
		"data": raspHosts, "total": total,
	}, "获取实例列表失败")
}

// 批量删除接口
func (h RaspHostController) BatchDeleteHostByIds(c *gin.Context) {
	var req vo.DeleteRaspHostRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}
	// 删除接口
	err := h.RaspHostRepository.DeleteRaspHost(req.Ids)
	if err != nil {
		response.Fail(c, nil, "删除实例失败: "+err.Error())
		return
	}
	response.Success(c, nil, "删除实例成功")
}

// 发布配置
func (h RaspHostController) PushConfig(c *gin.Context) {
	var req vo.PushConfigRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}

	// TODO 新建发布记录表
	configId := req.ConfigId
	config, err := h.RaspConfigRepository.GetRaspConfigById(configId)
	if err != nil {
		response.Fail(c, nil, "配置不存在")
		return
	}
	content, err := json.Marshal(config)
	if err != nil {
		response.Fail(c, nil, "获取配置文本失败:"+err.Error())
		return
	}
	// TODO 改成实际配置
	content = []byte(DEFAULT_CONFIG)
	for _, hostName := range req.HostNames {
		// 先判断连接是否存在
		m := socket.WebsocketManager.Group[hostName]
		if m != nil {
			client := m[hostName]
			if client != nil {
				socket.WebsocketManager.Send(hostName, hostName, content)
				response.Success(c, nil, "配置下发成功")
				return
			}
		}
	}
	response.Fail(c, nil, "配置下发失败")
}
