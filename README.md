# App Metric Simulator
A utility to simulate resource usage metrics including high CPU, memory, and I/O load for testing and monitoring purposes.

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

## Installation

### macOS
Simply download and extract the appropriate binary for your system (Intel or Apple Silicon):

```
./dex-simulator-mac # For Intel Macs
./dex-simulator-mac-arm64 # For Apple Silicon
```

You can also use the application bundle `DEXSimulator.app`.

### Windows

#### Handling Windows Security Warnings

When downloading or running `dex-simulator.exe` on Windows, you may encounter security warnings or antivirus detections. This is normal for unsigned executables and doesn't indicate any actual malware. Here's how to handle this:

**Method 1: Use the ZIP Package (Recommended)**

1. Download the `dex-simulator-windows.zip` file
2. Right-click the ZIP file and select "Properties"
3. At the bottom of the Properties window, check "Unblock" if present
4. Click "OK"
5. Extract the ZIP file
6. Run the extracted executable

**Method 2: Add an Exception to Windows Defender**

1. Open Windows Security (search for it in the Start menu)
2. Go to "Virus & threat protection"
3. Click "Manage settings" under "Virus & threat protection settings"
4. Scroll down to "Exclusions" and click "Add or remove exclusions"
5. Click "Add an exclusion" and select "File"
6. Browse to and select `dex-simulator.exe`

**Method 3: Run Despite the Warning**

When the SmartScreen warning appears:
1. Click "More info"
2. Click "Run anyway"

### Usage on Windows

Open a command prompt in the folder containing the executable and run:

```
dex-simulator.exe simulate --high=memory,cpu --frequency=constant --time=60
```
 