import json
import pandas as pd
from sklearn.ensemble import RandomForestClassifier
import joblib
import os

# File paths
model_file_path = '../models/trained_model1.pkl'
location_data_path = '../data/location_modified.json'
detected_fraud_path = '../data/detected_fraud.json'  # Update the path for detected fraud
label_encoders_path = '../models/label_encoders1.pkl'
expected_columns_path = '../models/expected_columns1.pkl'

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

# Debug: Print the DataFrame structure
print("Loaded Location Data:")
print(location_df.head())  # Print the first few rows
print("Columns in DataFrame:", location_df.columns.tolist())  # Print the column names

# Load label encoders
if not os.path.exists(label_encoders_path):
    raise FileNotFoundError(f"The label encoders file {label_encoders_path} does not exist.")
label_encoders = joblib.load(label_encoders_path)

# Apply label encoding to the features
for column, le in label_encoders.items():
    if column in location_df.columns:
        try:
            location_df[column] = le.transform(location_df[column])
        except ValueError as e:
            print(f"ValueError for column '{column}': {e}")
            # Handle unseen labels by replacing them with a default value
            location_df[column] = location_df[column].apply(
                lambda x: le.transform([x])[0] if x in le.classes_ else -1  # Replace unseen with -1
            )

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

# Map numeric predictions to string labels
fraud_labels = ['non-fraudulent' if pred == 0 else 'fraudulent' for pred in fraud_predictions]

# Debug: Print the first few predictions
print("Fraud Predictions:")
print(fraud_labels[:10])  # Print the first 10 predictions

# Store detected fraud cases
location_df['fraud'] = fraud_labels
fraud_cases = location_df[location_df['fraud'] == 'fraudulent']  # Select rows labeled as 'fraudulent'

# Debug: Check if any fraud cases were detected
if fraud_cases.empty:
    print("No fraud detected.")
else:
    print("Detected Fraud Cases:")
    print(fraud_cases.head())

# Save detected fraud cases to the new JSON file
with open(detected_fraud_path, 'w') as file:
    fraud_cases.to_json(file, orient='records', lines=True)

print(f"Detected fraud cases saved to {detected_fraud_path}")
