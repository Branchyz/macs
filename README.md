# MACS
A Windows-only Go CLI Utility for searching a IP address associated by a (portion of a) MAC address
The tool performs the following steps:
1. Finds the Local Gateway IP: Utilizes the route command to identify the gateway IP address.
2. Scans the Local Network: Executes an nmap scan on the local subnet.
3. Searches Scan Results: Searches all IP addresses corresponding to the given (portion of a) MAC address from the nmap output.

## Prerequisites
- Go: Ensure Go is installed on your system to build the application. Download it from golang.org.
- nmap: This tool relies on nmap for network scanning. Install it from nmap.org.

## Installation

1. Clone this repository:

   ```bash
   git clone https://github.com/branchyz/macs.git
   ```
2. Build the application:

   ```bash
   go build -o macs.exe
   ```
3. (Optional) Add the binary to your PATH:

    To make the macs command accessible from any directory, move it to a directory that is in your system's PATH, or add the current directory to your PATH. For example:
   
   ```bash
   set PATH=%PATH%;C:\path\to\your\bin
   ```

## Usage
Run the application with the MAC address as the command-line argument:

```bash
macf <MAC_ADDRESS>
```

Replace <MAC_ADDRESS> with the (portion of a) MAC address you want to resolve (e.g., 00:1A:2B:3C:4D:5E or C:4D:5).

The application will output the associated IP address or an error if it can't be found.

## Error Handling
1. **No MAC Address Provided**  
   **Error:** `errNoMac`  
   Description: Ensure you provide (a portion of) a MAC address as a command-line argument. Without it, the program cannot proceed.

2. **Invalid MAC Address**  
   **Error:** `errInvalidMac`  
   Description: The provided MAC address format is incorrect. Verify that the address matches the expected format (e.g., `00:1A:2B:3C:4D:5E` or `00:5E`).

3. **No Gateway Found**  
   **Error:** `errNoGateway`  
   Description: The default gateway could not be determined. This may indicate a network configuration issue.

4. **Unable to Parse Gateway**  
   **Error:** `errCantParseGateway`  
   Description: The program encountered an issue while parsing the gateway information. Verify your network configuration and try again.

