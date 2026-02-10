# Test Documentation

## Framework Choice
**Pytest** with **FastAPI TestClient**
- **Pytest**: Industry standard for Python testing - simple syntax, powerful fixtures, rich assertion introspection
- **TestClient**: FastAPI's built-in testing tool - lightweight, no HTTP overhead, preserves FastAPI context

## Test Structure

### **Unit Tests (11 tests total)**
```
test_main.py
├── Root Endpoint Tests
│   ├── Basic structure validation
│   ├── Service info validation
│   ├── System info validation
│   ├── Runtime info validation
│   ├── Request info validation
│   ├── Endpoints list validation
│   └── Complete JSON structure validation
│
├── Health Endpoint Tests 
│   ├── Basic health check
│   └── Uptime monotonic increase
│
├── Error Handling Tests
│   ├── 404 for non-existent endpoints
│   └── 405 for invalid methods
│
└── Parameterized Tests 
    └── Multiple 404 path variations
```

### **Test Categories**
1. **Structure Tests**: Validate JSON schema and required fields
2. **Data Integrity Tests**: Verify data types and logical constraints
3. **Error Case Tests**: Test error responses and HTTP status codes
4. **Behavior Tests**: Verify time-based behavior (uptime increase)

## How to Run Tests

### **Execution Commands**
```bash
# Run all tests
pytest

# Run specific file
pytest tests/test_endpoints.py 

# Run specific test
pytest tests/test_endpoints.py::test_various_nonexistent_paths
```

### **Output**
```
(.venv) arthur@Artur-MacBook-Pro app_python % pytest                                                        
============================================================ test session starts ============================================================
platform darwin -- Python 3.12.4, pytest-9.0.2, pluggy-1.6.0
rootdir: /Users/arthur/PycharmProjects/DevOps-Core-Course/app_python
plugins: anyio-4.12.1
collected 17 items                                                                                                                          

tests/test_endpoints.py .................                                                                                             [100%]

============================================================ 17 passed in 0.27s =============================================================
```
