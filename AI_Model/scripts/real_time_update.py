import json
import pandas as pd
from sklearn.ensemble import RandomForestClassifier
import joblib
import os

# File paths
feedback_file_path = '../data/feedback_data.json'
model_file_path = '../models/trained_model.pkl'
label_encoders_path = '../models/label_encoders.pkl'

# Check if feedback file exists
if not os.path.exists(feedback_file_path):
    raise FileNotFoundError(f"The file {feedback_file_path} does not exist.")

# Load feedback data
try:
    with open(feedback_file_path, 'r') as file:
        feedback_data = json.load(file)
except json.JSONDecodeError as e:
    raise ValueError(f"Error decoding JSON: {e}")

# Convert feedback data to DataFrame
feedback_df = pd.DataFrame(feedback_data)

# Check if model file exists
if not os.path.exists(model_file_path):
    raise FileNotFoundError(f"The model file {model_file_path} does not exist.")

# Load the existing model
model = joblib.load(model_file_path)

# Check if label encoders file exists
if not os.path.exists(label_encoders_path):
    raise FileNotFoundError(f"The label encoders file {label_encoders_path} does not exist.")

# Load label encoders
try:
    label_encoders = joblib.load(label_encoders_path)
except EOFError:
    raise ValueError(f"File {label_encoders_path} is corrupted or empty.")

# Ensure 'label' column exists in feedback data
if 'label' not in feedback_df.columns:
    raise KeyError("The 'label' column is missing in the feedback data.")

# Separate features and target
X_feedback = feedback_df.drop('label', axis=1)
y_feedback = feedback_df['label']

# Process each column with label encoders
for column, le in label_encoders.items():
    if column in X_feedback.columns:
        # Convert lists to strings
        X_feedback[column] = X_feedback[column].apply(lambda x: ','.join(map(str, x)) if isinstance(x, list) else x)
        # Transform using label encoder
        X_feedback[column] = le.transform(X_feedback[column])

# Ensure the columns match the expected columns from training
expected_columns_path = '../models/expected_columns.pkl'
if not os.path.exists(expected_columns_path):
    raise FileNotFoundError(f"The expected columns file {expected_columns_path} does not exist.")
expected_columns = joblib.load(expected_columns_path)

for column in expected_columns:
    if column not in X_feedback.columns:
        X_feedback[column] = 0  # Add missing columns with default values

X_feedback = X_feedback[expected_columns]  # Reorder columns to match the model's expectation

# Update the model with new feedback
model.fit(X_feedback, y_feedback)

# Save the updated model
joblib.dump(model, model_file_path)
