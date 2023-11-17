## System Design Documentation: Simple Match API

The **Simple Match API** is designed to facilitate matching and querying of people. It allows clients to add, remove, and query individuals in a memory-based store.

### API Endpoints

1. **POST /addPersonAndMatch**: Adds a single person to the memory store and performs matching.
2. **DELETE /removePerson/{id}**: Removes a person with the specified ID.
3. **GET /queryPeople**: Queries the stored people.

### Time Complexity Analysis

Let's analyze the time complexity of key operations in the API:

1. **Adding a Person and Matching (POST /addPersonAndMatch)**:
   - Adding a person involves inserting data into the memory store. Assuming the store uses a hash map or similar data structure, the average time complexity for insertion is O(1).
   - Matching involves comparing attributes of the newly added person with existing people. If we have 'n' people in the store, the matching operation would take O(n) time in the worst case (e.g., linear search).

2. **Removing a Person (DELETE /removePerson/{id})**:
   - Removing a person by ID requires locating the person in the store. Again, assuming a hash map, the average time complexity for deletion is O(1).

3. **Querying People (GET /queryPeople)**:
   - Querying all people in the store involves iterating through the entire collection. The time complexity for this operation is O(n), where 'n' is the number of people.
