# app-metric-simulator
mimic a app to simulate usage of metric, freeze, high cpu ,memory 

## Usage

To run the simulator for high memory and CPU usage with constant frequency for 10 minutes (600 seconds):

```
./dex-simulator-mac simulate --high=memory,cpy --frequency=constant --time=600
```

- `--high`: Resources to stress (memory, cpu, or both, comma separated)
- `--frequency`: Load frequency (`constant` or `random`)
- `--time`: Duration in seconds
