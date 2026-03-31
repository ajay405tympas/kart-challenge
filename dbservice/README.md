# Coupon ETL Project

## Setup

1. Install Go
2. Install dependency:
   go mod tidy

3. Update DB connection string in main.go

4. Run:
   go run main.go

## Input Files
Place your coupon files:
- couponbase1.txt
- couponbase2.txt
- couponbase3.txt
- Please unzip  the file and copy the file into dBService folder
- Please note you need to have postgresSQL running in the machine and you need to create a database table 'coupons' and by executing the command run

## Output
Data will be loaded into PostgreSQL table:
coupons(code TEXT PRIMARY KEY, counter INT)
