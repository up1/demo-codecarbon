import os
import time
from datetime import datetime
from codecarbon import track_emissions

@track_emissions(
    api_key=os.getenv("CODECARBON_API_KEY"),
    experiment_name="demo01",
)

# check member in large list performance
def check_member_in_list(input_list, item):
    """
    Checks if an item is a member of a list and measures the time taken for the operation.

    Args:
        input_list: A list to search through
        item: The element to search for in the list

    Returns:
        tuple: A tuple containing:
            - bool: True if item is in large_list, False otherwise
            - float: Time elapsed during the search operation in seconds
    """
    start_time = time.time()
    is_member = item in input_list
    end_time = time.time()
    search_time = end_time - start_time
    return is_member, search_time

if __name__ == "__main__":
    # Create a large list
    large_list = set(range(10_000_000))  # Set of integers from 0 to 9999999

    # Check for an item in the list
    ITEM_TO_CHECK = 9_999_999
    result, elapsed_time = check_member_in_list(large_list, ITEM_TO_CHECK)

    print(f"Item {ITEM_TO_CHECK} found: {result}")
    print(f"Time taken to check membership: {elapsed_time:.6f} seconds")