# app-metric-simulator
mimic a app to simulate usage of metric, freeze, high cpu ,memory 

## Usage

To run the simulator for high memory and CPU usage with constant frequency for 10 minutes (600 seconds):

```
./dex-simulator-mac simulate --high=memory,cpy --frequency=constant --time=600
```

- `--high`: Resources to stress (memory, cpu, io, or combination, comma separated)
- `--frequency`: Load frequency (`constant` or `random`)
- `--time`: Duration in seconds
- `--crash-time`: Time in seconds after which to simulate a crash (0 means no crash)
- `--max-cpu`: Maximum CPU usage percent (default 60)
- `--max-memory`: Maximum memory usage in GB (default 10GB)

### Resource Types
- `memory`: Simulates high memory usage
- `cpu`: Simulates high CPU usage
- `io`: Simulates intensive disk I/O operations

### Examples

Example commands:

Simulate high I/O load:
```
./dex-simulator-mac simulate --high=io --frequency=constant --time=30
```

Simulate high memory, CPU, and I/O together:
```
./dex-simulator-mac simulate --high=memory,cpu,io --frequency=constant --time=30
```

Simulate a crash after 10 seconds:
```
./dex-simulator-mac simulate --high=memory --frequency=constant --time=30 --crash-time=10
```
 