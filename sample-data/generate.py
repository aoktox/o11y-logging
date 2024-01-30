import os, time, datetime, random

def generate_sample_data():
    timestamp = datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d %H:%M:%S")
    service_names = ["checkout", "payment", "shipping", "authentication", "inventory", "customer_support"]
    status_codes = [200, 201, 204, 400, 401, 403, 404, 500, 503]
    response_time = random.randint(10, 2000)
    user_id = "user-" + str(random.randint(1000, 9999))
    transaction_id = "tx" + str(random.randint(1000, 9999))
    additional_info = "Sample log entry for testing."

    log_entry = f"{timestamp} {random.choice(service_names)} {random.choice(status_codes)} {response_time}ms {user_id} {transaction_id} {additional_info}"
    return log_entry


def main():
    script_directory = os.path.dirname(os.path.realpath(__file__))
    log_file_path = os.path.join(script_directory, "sample.log")
    
    with open(log_file_path, "a") as log_file:
        counter=0
        for _ in range(1000):  # Generate 10 sample log entries
            log_entry = generate_sample_data()
            log_file.write(log_entry + "\n")
            counter+=1
            if counter == 10:
                counter=0
                time.sleep(1)



if __name__ == "__main__":
    main()