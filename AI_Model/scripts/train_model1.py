import json
import pandas as pd
from sklearn.ensemble import RandomForestClassifier
from sklearn.model_selection import train_test_split, GridSearchCV
from sklearn.metrics import classification_report, confusion_matrix
from sklearn.preprocessing import LabelEncoder
import joblib
import os

def load_data(file_path):
    with open(file_path, 'r') as file:
        return json.load(file)

def preprocess_data(df):
    # Flatten any list columns
    for column in df.select_dtypes(include=['object']).columns:
        if df[column].apply(lambda x: isinstance(x, list)).any():
            df[column] = df[column].apply(lambda x: ','.join(map(str, x)) if isinstance(x, list) else x)
    return df

def handle_missing_values(df):
    return df.fillna(method='ffill')  # Forward fill as an example

def encode_labels(X):
    label_encoders = {}
    for column in X.select_dtypes(include=['object']).columns:
        le = LabelEncoder()
        X[column] = le.fit_transform(X[column])
        label_encoders[column] = le
    return X, label_encoders

def save_objects(label_encoders, expected_columns, model, encoders_path, expected_columns_path, model_path):
    joblib.dump(label_encoders, encoders_path)
    joblib.dump(expected_columns, expected_columns_path)
    joblib.dump(model, model_path)

def main():
    # Load the training data
    df = pd.DataFrame(load_data('../data/training_datamodified.json'))
    
    df = preprocess_data(df)
    df = handle_missing_values(df)

    X = df.drop('label', axis=1)  # Features
    y = df['label']  # Target labels

    X, label_encoders = encode_labels(X)

    expected_columns = X.columns.tolist()
    
    X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

    model = RandomForestClassifier(n_estimators=100, random_state=42)

    # Optionally perform Grid Search for hyperparameter tuning
    # param_grid = {'n_estimators': [50, 100, 200], 'max_depth': [None, 10, 20]}
    # grid_search = GridSearchCV(model, param_grid, cv=5)
    # grid_search.fit(X_train, y_train)
    # model = grid_search.best_estimator_

    model.fit(X_train, y_train)
    
    y_pred = model.predict(X_test)
    print(classification_report(y_test, y_pred))
    print(confusion_matrix(y_test, y_pred))

    save_objects(label_encoders, expected_columns, model, 
                 '../models/label_encoders1.pkl', 
                 '../models/expected_columns1.pkl', 
                 '../models/trained_model1.pkl')

if __name__ == "__main__":
    main()
