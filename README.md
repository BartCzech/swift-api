# Swift API - SWIFT Code Management System

A RESTful API built with **Go** and **PostgreSQL**, designed to manage **SWIFT codes** for banks worldwide. The API provides functionality to retrieve, create, and delete SWIFT codes efficiently.

## 🚀 Features
- Retrieve details of a **single SWIFT code** (including headquarters and branches).
- Get all SWIFT codes **for a specific country**.
- Create a **new SWIFT code entry**.
- Delete an **existing SWIFT code**.
- Fully **containerized** using **Docker** and **Docker Compose**.

## 🛠️ Setup Instructions

### **1️⃣ Install Docker & Docker Compose**
#### Linux:
```sh
sudo apt update
sudo apt install docker.io docker-compose -y
sudo systemctl enable --now docker
```
#### Windows/Mac:
- Download and install **Docker Desktop** from [here](https://www.docker.com/products/docker-desktop).
- Ensure Docker is **running** before proceeding.

### **2️⃣ Clone the Repository**
```sh
git clone https://github.com/YourUsername/swift-api.git
cd swift-api
```

### **3️⃣ Include the `.env` File**
Add the `.env` file to the swift-api folder.

### **4️⃣ Build and Start the Containers**
Launch the terminal in the swift-api folder and enter the following command (works on both Windows and Linux):
```sh
docker-compose up --build
```
This will build the **Go application**, start **PostgreSQL** and set up the **environment variables**.

### **5️⃣ Test the API**
After the containers start, check if the API is running:

#### **Health check route:**
```sh
curl http://localhost:8080/ping
```
Expected response:
```json
{"message": "pong"}
```

## 🛑 Stopping the Service
To **stop** the containers **without deleting data**, press **CTRL+C** or run:
```sh
docker-compose down
```

To **stop and remove all data** (including the database):
```sh
docker-compose down -v
```

## 📌 API Endpoints

| Method | Endpoint | Description |
|--------|---------|-------------|
| **GET** | `/ping` | Health check |
| **GET** | `/v1/swift-codes` | Get all SWIFT codes |
| **GET** | `/v1/swift-codes/:swift-code` | Get details of a specific SWIFT code |
| **GET** | `/v1/swift-codes/country/:countryISO2code` | Get SWIFT codes for a country |
| **POST** | `/v1/swift-codes` | Create a new SWIFT code |
| **DELETE** | `/v1/swift-codes/:swift-code` | Delete a SWIFT code |