import json

# Load the data from responses.json
with open('../data/responses.json', 'r') as file:
    data = json.load(file)

# Initialize an empty list to store the extracted location information
location_info = []

# Iterate through each subscriber's data
for subscriber in data:
    # Extract the necessary fields
    extracted_info = {
        "subscriptionId": subscriber.get("subscriptionId"),
        "supi": subscriber.get("supi"),
        "timestamp": subscriber.get("timeStamp"),
        "location": subscriber.get("location"),
        "timezone": subscriber.get("timezone")
    }
    # Add the extracted information to the list
    location_info.append(extracted_info)

# Store the extracted location information in a new JSON file
with open('../data_files/location_info.json', 'w') as output_file:
    json.dump(location_info, output_file, indent=4)

print("Location information has been extracted and stored in location_info.json")
