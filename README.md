# HyFabric-Herb | 中草药溯源系统 (Hyperledger Fabric + Go Gin + REST + JS)

> 基于 **Hyperledger Fabric** 的中草药全流程溯源解决方案：从种植、加工、质检到零售查询。  
> 后端 **Go (Gin)** 提供 **RESTful API**，前端 **JavaScript** 调用，链码默认使用 **Go**。

[![Go Version](https://img.shields.io/badge/Go-%3E=1.21-blue.svg)]()
[![Fabric](https://img.shields.io/badge/Hyperledger-Fabric-2.x-6f42c1.svg)]()
[![License](https://img.shields.io/badge/license-MIT-green.svg)]()

---

## 目录
- [概述](#概述)
- [系统架构](#系统架构)
- [功能特性](#功能特性)
- [目录结构](#目录结构)
- [先决条件](#先决条件)
- [快速开始](#快速开始)
  - [1. 启动本地 Fabric 网络](#1-启动本地-fabric-网络)
  - [2. 部署链码](#2-部署链码)
  - [3. 启动后端 (Go Gin)](#3-启动后端-go-gin)
  - [4. 启动前端 (JavaScript)](#4-启动前端-javascript)
- [环境变量](#环境变量)
- [REST API 概览](#rest-api-概览)
- [链码接口 (示例)](#链码接口-示例)
- [数据模型](#数据模型)
- [开发与测试](#开发与测试)
- [生产部署要点](#生产部署要点)
- [Roadmap](#roadmap)
- [贡献](#贡献)
- [许可证](#许可证)

---

## 概述
中草药供应链复杂、参与主体多、信息孤岛严重。本项目利用 **Fabric** 的多组织/通道、背书策略、私有数据与不可篡改特性，实现：
- 批次级别追踪（采集 → 加工 → 质检 → 运输 → 零售）
- 证书与检验报告上链存证（可选私有数据）
- 标准化 **RESTful API** 供前端/第三方系统集成

---

## 系统架构

```mermaid
flowchart LR
  subgraph Client[Web / Mobile (JavaScript)]
    UI[Traceability UI] --> FE[API Client]
  end

  FE -->|HTTPS/JSON| GIN[Go Gin REST API]
  GIN --> SDK[Fabric SDK (Go)]
  SDK --> PEER1[Peer Org1]
  SDK --> PEER2[Peer Org2]
  PEER1 <---> ORDERER[Orderer RAFT]
  PEER2 <---> ORDERER
  PEER1 <---> CDB1[CouchDB]
  PEER2 <---> CDB2[CouchDB]
  PEER1 --> CC[(Chaincode: herbcc)]
  PEER2 --> CC
