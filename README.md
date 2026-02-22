# drawbridge

Drawbridge is a lightweight, high-performance port knocking daemon written in Go. It utilizes `libpcap` to monitor network traffic for a predefined sequence of TCP SYN packets. Upon detecting a valid sequence, it dynamically modifies `iptables` rules to grant temporary access to a protected port (e.g., SSH) and automatically revokes access after a configured timeout.

By operating at the packet capture level with Berkeley Packet Filters (BPF), Drawbridge remains stealthy and minimizes overhead, dropping irrelevant traffic before it reaches the application logic.

## Features

- **Stealth Operation:** Listens passively via `libpcap`. The daemon itself does not open any listening sockets.
- **BPF Filtering:** Highly efficient packet filtering offloaded to the kernel.
- **Automatic Revocation:** Access is granted temporarily and revoked automatically via background timers.
- **Systemd Integration:** Includes deployment scripts for easy setup as a background service.

## Prerequisites

- Linux OS with `iptables` installed.
- Root privileges (required for `libpcap` and `iptables` management).
- Go 1.25+ (for building from source).
- `libpcap-dev` (required for CGO compilation of the `gopacket` dependency).

On Debian/Ubuntu:
```bash
sudo apt update
sudo apt install libpcap-dev iptables
```

## Installation

A deployment script is provided to compile the binary, set up the configuration directory, and register the `systemd` service.

```bash
git clone [https://github.com/sadsnake231/drawbridge.git](https://github.com/sadsnake231/drawbridge.git)
cd drawbridge
chmod +x install.sh
sudo ./install.sh
```

## Configuration

The default configuration file is installed at `/etc/drawbridge/config.yaml`. Modify it to suit your network environment.

Sequence recommendations:

- Avoid sequentional ports
- Prefer ports in the `1024-65535` range



```YAML
interface: "eth0"               # Network interface to monitor
sequence:                       # Secret port sequence
  - 1111
  - 2222
  - 3333
knock-timeout: 15s              # Max time allowed to complete the sequence
safe-port: 22                   # The target port to expose (e.g., SSH)
close-timeout: 15m              # Duration to keep the port open
log-file: /var/log/drawbridge.log # Path to the daemon log file
snaplen: 1024                   # Packet capture snapshot length
promisc: false                  # Promiscuous mode flag
bpf-filter: "tcp[tcpflags] & (tcp-syn) != 0" # BPF expression to filter traffic
```

The default setup operates under the following constraints:

- Captures only TCP SYN packets (tcp[tcpflags] & (tcp-syn) != 0). UDP, ICMP, and established TCP connections are ignored at the kernel level.

- Knock Window (knock-timeout): You have 15 seconds to complete the entire port sequence. If the time expires or an incorrect port is hit, the sequence resets.

- Access Duration (close-timeout): Upon a successful knock, the protected port (e.g., port 22) is opened for exactly 15 minutes. After this period, the firewall rule is automatically revoked.

- Capture Efficiency: Reads only the first 1024 bytes (snaplen) of each packet and operates with promiscuous mode disabled (promisc: false) to minimize CPU load.

After modifying the configuration, restart the service to apply changes:
```bash
sudo systemctl restart drawbridge
```

## Usage

```
sudo systemctl start drawbridge
sudo systemctl stop drawbridge
```

## Logs

Drawbridge writes operational logs to the specified log file.

```bash
sudo tail -f /var/log/drawbridge.log # or whatever path you specified in config
```

## Client authentication

To gain access to the protected port, send TCP SYN packets to the configured sequence of ports within the knock-timeout window. You can use tools like `netcat` or `nmap`.

`Netcat`:

```bash
for port in 1111 2222 3333; do nc -z -w 1 <server_ip> $port; done
```

## Uninstallation
To completely remove the binary and the systemd service, run the provided uninstallation script:
```bash
chmod +x uninstall.sh
sudo ./uninstall.sh
```
