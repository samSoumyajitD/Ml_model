import json
import pandas as pd
from sklearn.ensemble import RandomForestClassifier
import joblib
import os

# File paths
model_file_path = '../models/trained_model.pkl'
location_data_path = '../data/location_info.json'
fraud_cases_path = '../data/detected_fraud.json'
label_encoders_path = '../models/label_encoders.pkl'
expected_columns_path = '../models/expected_columns.pkl'

# Load the model
if not os.path.exists(model_file_path):
    raise FileNotFoundError(f"The model file {model_file_path} does not exist.")
model = joblib.load(model_file_path)

# Load location data
if not os.path.exists(location_data_path):
    raise FileNotFoundError(f"The location data file {location_data_path} does not exist.")
with open(location_data_path, 'r') as file:
    location_data = json.load(file)

# Convert JSON to DataFrame
location_df = pd.DataFrame(location_data)

# Flatten the nested JSON structure
def flatten_location_info(df):
    df['location_nrLocation_tai_plmnId_mcc'] = df['location'].apply(lambda x: x['nrLocation']['tai']['plmnId']['mcc'])
    df['location_nrLocation_tai_plmnId_mnc'] = df['location'].apply(lambda x: x['nrLocation']['tai']['plmnId']['mnc'])
    df['location_nrLocation_tai_tac'] = df['location'].apply(lambda x: x['nrLocation']['tai']['tac'])
    df['location_nrLocation_ncgi_plmnId_mcc'] = df['location'].apply(lambda x: x['nrLocation']['ncgi']['plmnId']['mcc'])
    df['location_nrLocation_ncgi_plmnId_mnc'] = df['location'].apply(lambda x: x['nrLocation']['ncgi']['plmnId']['mnc'])
    df['location_nrLocation_ncgi_nrCellId'] = df['location'].apply(lambda x: x['nrLocation']['ncgi']['nrCellId'])
    df['location_ageOfLocationInformation'] = df['location'].apply(lambda x: x['nrLocation']['ageOfLocationInformation'])
    df['location_ueLocationTimestamp'] = df['location'].apply(lambda x: x['nrLocation']['ueLocationTimestamp'])
    df = df.drop(columns=['location'])
    return df

location_df = flatten_location_info(location_df)

# Debug: Check the first few rows of flattened data
print("Flattened Location Data:")
print(location_df.head())

# Load label encoders
if not os.path.exists(label_encoders_path):
    raise FileNotFoundError(f"The label encoders file {label_encoders_path} does not exist.")
label_encoders = joblib.load(label_encoders_path)

# Apply label encoding to the features
for column, le in label_encoders.items():
    if column in location_df.columns:
        location_df[column] = le.transform(location_df[column])

# Debug: Check if label encoding was applied correctly
print("Data after Label Encoding:")
print(location_df.head())

# Load expected columns
if not os.path.exists(expected_columns_path):
    raise FileNotFoundError(f"The expected columns file {expected_columns_path} does not exist.")
expected_columns = joblib.load(expected_columns_path)

# Ensure the DataFrame columns match the model's expectations
for column in expected_columns:
    if column not in location_df.columns:
        location_df[column] = 0  # Add missing columns with default values

# Debug: Check for missing columns
missing_columns = [col for col in expected_columns if col not in location_df.columns]
if missing_columns:
    print(f"Missing columns: {missing_columns}")

location_df = location_df[expected_columns]  # Reorder columns to match the model's expectation

# Debug: Check the final DataFrame columns
print("DataFrame Columns for Prediction:")
print(location_df.head())

# Predict potential fraud
fraud_predictions = model.predict(location_df)

# Debug: Print the first few predictions
print("Fraud Predictions:")
print(fraud_predictions[:10])  # Print the first 10 predictions

# Store detected fraud cases
location_df['fraud'] = fraud_predictions
fraud_cases = location_df[location_df['fraud'] == 1]  # Assuming '1' indicates fraud

# Debug: Check if any fraud cases were detected
if fraud_cases.empty:
    print("No fraud detected.")
else:
    print("Detected Fraud Cases:")
    print(fraud_cases.head())

# Save detected fraud cases to a JSON file
with open(fraud_cases_path, 'w') as file:
    fraud_cases.to_json(file, orient='records', lines=True)
