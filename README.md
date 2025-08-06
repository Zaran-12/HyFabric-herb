# HyFabric-Herb | 中草药溯源系统 (Hyperledger Fabric + Go Gin + REST + JS)

> 基于 **Hyperledger Fabric** 的中草药全流程溯源解决方案：从种植、加工、质检到零售查询。  
> 后端 **Go (Gin)** 提供 **RESTful API**，前端 **JavaScript** 调用，链码默认使用 **Go**。

[![Go Version](https://img.shields.io/badge/Go-%3E=1.21-blue.svg)]()
[![Fabric](https://img.shields.io/badge/Hyperledger-Fabric-2.x-6f42c1.svg)]()
[![License](https://img.shields.io/badge/license-MIT-green.svg)]()

---



## 概述
中草药供应链复杂、参与主体多、信息孤岛严重。本项目利用 **Fabric** 的多组织/通道、背书策略、私有数据与不可篡改特性，实现：
- 批次级别追踪（采集 → 加工 → 质检 → 运输 → 零售）
- 证书与检验报告上链存证（可选私有数据）
- 标准化 **RESTful API** 供前端/第三方系统集成


