import json
import pandas as pd
from sklearn.ensemble import RandomForestClassifier
from sklearn.model_selection import train_test_split
from sklearn.metrics import classification_report
from sklearn.preprocessing import LabelEncoder
import joblib
import os

# Load the training data
with open('../data/training_data.json', 'r') as file:
    training_data = json.load(file)

# Convert to DataFrame
df = pd.DataFrame(training_data)

# Flatten any list columns
for column in df.select_dtypes(include=['object']).columns:
    if df[column].apply(lambda x: isinstance(x, list)).any():
        # Convert lists to strings
        df[column] = df[column].apply(lambda x: ','.join(map(str, x)) if isinstance(x, list) else x)

# Separate features and target
X = df.drop('label', axis=1)  # Features
y = df['label']  # Target labels

# Convert categorical columns to numeric
label_encoders = {}
for column in X.select_dtypes(include=['object']).columns:
    le = LabelEncoder()
    X[column] = le.fit_transform(X[column])
    label_encoders[column] = le

# Save label encoders
label_encoders_path = '../models/label_encoders.pkl'
joblib.dump(label_encoders, label_encoders_path)

# Save the list of expected columns
expected_columns = X.columns.tolist()
expected_columns_path = '../models/expected_columns.pkl'
joblib.dump(expected_columns, expected_columns_path)

# Split into train and test sets
X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

# Initialize the model
model = RandomForestClassifier(n_estimators=100, random_state=42)

# Train the model
model.fit(X_train, y_train)

# Evaluate the model
y_pred = model.predict(X_test)
print(classification_report(y_test, y_pred))

# Save the trained model
model_path = '../models/trained_model.pkl'
joblib.dump(model, model_path)
