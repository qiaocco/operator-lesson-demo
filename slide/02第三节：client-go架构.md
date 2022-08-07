---
marp: true
theme: gaia
paginate: true
footer: '@白丁云原生'
backgroundColor: white
style: |
    code {
        background: black;
    }
---
<!--
_class: lead
-->
# client-go架构

---
## 代码结构

![width:26cm height:14cm](./images/client-go.png)

- kubernetes: 包含了所有访问k8s api的clientset
- informers: 包含所有内置资源的informer，便于操作k8s的资源
- listers: 包含所有对象lister，用于读取缓存中k8s资源对象的信息
- discovery: 用于发现 api server支持的api
- dynamic: 通过dynamic client，可以操作任意k8s api对象，包括我们自定义的资源
- plugin: 插件
- transport: 创建连接，认证的逻辑，clientset使用
- tools: 工具

---
## 架构
![width:26cm height:14cm](./images/design.png)

---

## 控制器逻辑
- 观察：通过监控Kubernetes资源对象变化的事件来获取当前对象状态，我们只需要注入EventHandler让client-go将变化的事件对象信息放入WorkQueue中。
- 分析：确定当前状态和期望状态的不同，由Worker完成。
- 执行：执行能够驱动对象当前状态变化的操作，由Worker完成。
- 更新：更新对象的当前状态，由Worker完成。



