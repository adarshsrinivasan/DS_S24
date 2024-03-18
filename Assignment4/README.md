# CSCI 5673 - Distributed Systems
### Assignment 1 - Online E-commerce marketplace using TCP/IP sockets

The design comprises seven logical components:
1. Client-Buyer
2. Client-Seller
3. Server-Buyer
4. Server-Seller
5. Customer database - consists of Buyer and Seller user details
6. Product database - consists of Product details - using **Mongo** db.
7. Financial Transactions database - using **Postgres** db

Assumptions:
1. For Customer database, we are using **Postgres** db since the details are structured and can be stored in the form of tables.
2. For Product database, we are using **MongoDB** since the item details such as the keywords are lists which are unstructured. In this case, using Mongo can be efficient, for searching and inserting the data.

Current State of the system:
- [x] Server listens and handles Buyer requests
- [x] Server listens and handles Seller requests
- [x] Client - Buyer connects to the Server - Buyer and exchanges information
- [x] Client - Seller connects to the Server - Seller and exchanges information
- [x] Customer database up and running
- [x] Product database up and running

