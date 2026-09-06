def run_data_agent():
    print("DATA AGENT")
    print("Reading Databricks data...")
    
    customers_path = "/Volumes/dbacademy/default/raw_data/rides/build_data_pipeline_demo/customers.csv"
    sales_path = "/Volumes/dbacademy/default/raw_data/rides/build_data_pipeline_demo/sales.csv"

    print(f"Customers: {customers_path}")
    print(f"Sales: {sales_path}")

    print("Data Agent completed.")