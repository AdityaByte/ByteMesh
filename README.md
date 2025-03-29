# ByteMesh - A Distributed File System (DFS) in Go  

**ByteMesh** is a **Distributed File System (DFS)** built using **Golang**, designed to provide **scalable, fault-tolerant, and efficient file storage across multiple machines**. Unlike traditional DFS implementations such as **Hadoop HDFS, Google File System (GFS), or Ceph**, ByteMesh offers a **lightweight yet powerful distributed storage system** with a **TCP-based client-server architecture**.  

## 🚀 Features  

✅ **True Distributed Storage** – The server can be deployed on multiple machines, handling distributed file storage.  
✅ **TCP-Based Client-Server Communication** – Uses raw TCP sockets for high-performance data transfer.  
✅ **Metadata Management** – A dedicated metadata service tracks file locations, replication, and integrity.  
✅ **File Chunking & Replication** – Files are split into chunks and distributed across nodes for redundancy.  
✅ **Scalability** – Designed to scale by adding more storage nodes.     

## 📌 Installation  

### **Prerequisites**  
- Go (Golang) **1.18+**  
- MongoDB
- Git  

### **Clone the Repository**  
```sh
git clone https://github.com/AdityaByte/ByteMesh.git
cd ByteMesh
