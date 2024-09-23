import subprocess
import sys

def run_script(script_name):
    """Run a Python script and print its output."""
    try:
        result = subprocess.run(['python', script_name], capture_output=True, text=True, check=True)
        print(f"Output of {script_name}:\n{result.stdout}")
        if result.stderr:
            print(f"Error in {script_name}:\n{result.stderr}")
    except subprocess.CalledProcessError as e:
        print(f"An error occurred while running {script_name}: {e}")
        sys.exit(1)

def main():
    scripts = [
        'subscriberDatafetch.py',
        'fetchLocation.py',
        'real_time_update.py',
        'detect_fraud.py'
    ]

    for script in scripts:
        run_script(script)

if __name__ == '__main__':
    main()
