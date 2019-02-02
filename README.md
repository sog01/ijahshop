# Ijah Inventory
Simple API for inventory online shop using go language. This repository codebase implement 3 tier architecture and dependency injection. 

## Prequisite
1. golang
2. sqlite
3. dep (optionaly)

## How to run this app

To run this application locally, simply type **make run**. Application will serve at **http://localhost:8080**. It will also generate binary files on root directory. 

### Optional
Vendor files generated using dep package manager, if you want to re-generate this file, simple run **dep ensure -v --vendor-only**.

## Demonstration

This app already run on my VPS. You can access directly at http://ijah.abdullah-dev.tech.

## Documentation
API Documentation can be accessed at https://web.postman.co/collections/6239183-5e39b4b8-6fba-4a3f-9df3-cd7d4a8d9cca?workspace=95bf9ad3-7234-4920-8278-30a6ae7432f8#26380cf8-1459-4303-ba00-c92784d6b26b

## Data Model
![alt text](https://abdullah-dev.tech/images/ijah_model.jpg "Inventory database model")

### Product
Product is model that represent Catatan Jumlah Barang. These field respectively represent each column (in excel) that shown below :

1. **name** represent Nama Item
2. **sku** represent SKU
3. **stock** represent Jumlah Sekarang

### Purchase
Purchase is model that represent **Catatan Barang Masuk**. These field respectively represent each column (in excel) that shown below :

1. **quantity_order** represent Jumlah Pemesanan
2. **quantity_accepted** represent Jumlah Diterima
3. **description** represent Catatan
4. **invoice_number** represent Nomer Kuitansi
5. **cost** represent Harga Beli
6. **date** represent Waktu

### Orders
Orders is model that represent **Catatan Barang Keluar**. These field respectively represent each column (in excel) that shown below :

1. **quantity** represent Jumlah Keluar
2. **description** represent Catatan
3. **date** represent Waktu
4. **price** represent Harga Jual

For remaining column that doesn't mention at here, has been calculated by field that shown on picture above. It means that column doesn't have original value. That's why these column not be created as a schema. 



 
