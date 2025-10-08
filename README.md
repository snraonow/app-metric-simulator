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

When downloading or running `dex-simulator.exe` on Windows, you may encounter security warnings or antivirus detections such as "Trojan:Win32/Sabsik". This is a false positive that occurs with many Go applications and doesn't indicate actual malware. Here's how to handle this:

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

**Method 4: Try Alternative Package (For Sabsik Detection)**

If you're specifically getting the "Trojan:Win32/Sabsik" warning:

1. Download the `app-simulator.zip` file instead
2. This contains the same program but with a generic filename that may avoid detection
3. Follow the same unblocking steps as Method 1

**Method 5: Disable Real-time Protection Temporarily**

As a last resort, you can temporarily disable real-time protection:

1. Open Windows Security
2. Click "Virus & threat protection"
3. Under "Virus & threat protection settings" click "Manage settings"
4. Toggle "Real-time protection" to Off
5. Run the application
6. **Important:** Turn real-time protection back on when finished

### Usage on Windows

Open a command prompt in the folder containing the executable and run:

```
dex-simulator.exe simulate --high=memory,cpu --frequency=constant --time=60
```

### Why Does This Happen?

The "Trojan:Win32/Sabsik" detection is a false positive that commonly occurs with Go applications. This happens because:

1. Go executables contain certain patterns similar to known malware
2. The simulator performs resource-intensive operations that security software flags as suspicious
3. The application is unsigned (doesn't have a digital signature)

Rest assured that this application does not contain any malware or malicious code. The source code is fully available for inspection.
 