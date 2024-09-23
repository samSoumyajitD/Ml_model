import pymongo
import json
from bson.json_util import loads

# MongoDB connection details
mongo_uri = "mongodb://localhost:27017"
database_name = "free5gc"
collection_name = "dataCollectionNwdar"

# Connect to MongoDB
client = pymongo.MongoClient(mongo_uri)
db = client[database_name]
collection = db[collection_name]

# Fetch the document(s) from the collection
documents = collection.find()

# List to hold the extracted reportList data
report_data = []

for document in documents:
    # Parse the 'response' field which is a JSON string
    response_json = loads(document['response'])
    
    # Extract the 'reportList' from the response
    report_list = response_json.get("reportList", [])
    
    # Append the extracted reportList to the report_data list
    report_data.extend(report_list)

# Save the extracted report data to a JSON file
with open('../data_files/report_data.json', 'w') as file:
    json.dump(report_data, file, indent=4)

# Print the extracted report data to confirm
print(json.dumps(report_data, indent=4))

# Close the MongoDB connection
client.close()
